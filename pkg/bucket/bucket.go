package bucket

import (
	"sync"
	"time"
)

type Bucket interface {
	Check(key string) bool
	Cleanup()
	ResetKey(key string)
}

type bucket struct {
	mu        sync.Mutex
	maxTokens int
	maxRetry  int
	tokens    map[string]*Token
	interval  float64
}

type Token struct {
	expire     time.Time
	countRetry int
}

func New(maxTokens int, maxRetry int, interval float64) Bucket {
	return &bucket{
		mu:        sync.Mutex{},
		maxTokens: maxTokens,
		maxRetry:  maxRetry,
		tokens:    make(map[string]*Token, maxTokens),
		interval:  interval,
	}
}

func (b *bucket) Check(key string) bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	now := time.Now()
	var t *Token
	if token, ok := b.tokens[key]; !ok {
		// Достигнут лимит корзины.
		if len(b.tokens) >= b.maxTokens {
			return false
		}

		t = &Token{
			expire:     now.Add(time.Duration(b.interval) * time.Second),
			countRetry: 0,
		}
		b.tokens[key] = t
	} else {
		t = token
	}

	// Если последняя попытка была давно, то обнуляем счетчик по ключу.
	if now.After(t.expire) {
		t.countRetry = 1
	} else {
		t.countRetry++
	}
	// При каждой попытке указываем новое время истечения действия счетчика.
	t.expire = now.Add(time.Duration(b.interval) * time.Second)
	return t.countRetry <= b.maxRetry
}

func (b *bucket) Cleanup() {
	b.mu.Lock()
	defer b.mu.Unlock()
	now := time.Now()
	for key, token := range b.tokens {
		if now.After(token.expire) {
			delete(b.tokens, key)
		}
	}
}

func (b *bucket) ResetKey(key string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	delete(b.tokens, key)
}
