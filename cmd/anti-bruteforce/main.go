package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Al-Sher/otus_anti_bruteforce/internal/app"
	"github.com/Al-Sher/otus_anti_bruteforce/internal/config"
	"github.com/Al-Sher/otus_anti_bruteforce/internal/list"
	"github.com/Al-Sher/otus_anti_bruteforce/internal/logger"
	"github.com/Al-Sher/otus_anti_bruteforce/internal/ratelimit"
	"github.com/Al-Sher/otus_anti_bruteforce/internal/server"
	"github.com/redis/go-redis/v9"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	ctx := context.Background()

	logg, err := logger.New(cfg.LogLevel)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	rURL, err := redis.ParseURL(cfg.RedisURL)
	if err != nil {
		logg.Error(err)
		os.Exit(1)
	}

	rClient := redis.NewClient(rURL)

	whiteList, err := list.New(ctx, cfg.WhiteListRedisKey, *rClient)
	if err != nil {
		logg.Error(err)
		os.Exit(1)
	}

	blackList, err := list.New(ctx, cfg.BlackListRedisKey, *rClient)
	if err != nil {
		logg.Error(err)
		os.Exit(1)
	}

	rt := ratelimit.New(cfg.LoginLimit, cfg.PasswordLimit, cfg.IPLimit, cfg.BucketSize, cfg.BlockInterval)

	a := app.New(rt, whiteList, blackList)
	s := server.New(cfg.Addr, a, logg)
	ctx, cancel := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	go shutdown(ctx, s, logg)
	go clearRateLimit(ctx, rt, cfg.BlockInterval)

	logg.Info("anti-bruteforce is running...")

	if err := s.Start(ctx); err != nil {
		logg.Error(err)
		cancel()
	}

	<-ctx.Done()
}

func shutdown(ctx context.Context, s server.Server, logg logger.Logger) {
	<-ctx.Done()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	if err := s.Stop(ctx); err != nil {
		logg.Error("failed to stop http server: " + err.Error())
	}

	logg.Info("anti-bruteforce is shutdown...")
}

func clearRateLimit(ctx context.Context, rt ratelimit.RateLimit, interval float64) {
	c := time.Tick(time.Duration(interval) * time.Second)
	for {
		select {
		case <-c:
			rt.Cleanup()
		case <-ctx.Done():
			return
		}
	}
}
