package main

import (
	"golang.org/x/time/rate"
	"sync"
	"time"
)

// IPRateLimiter .
type IPRateLimiter struct {
	ips map[string]*LimiterEntry
	mu  *sync.RWMutex
	r   rate.Limit
	b   int
}

type LimiterEntry struct {
	limiter *rate.Limiter
	lastUse time.Time
}

// NewIPRateLimiter .
func NewIPRateLimiter(r rate.Limit, b int) *IPRateLimiter {
	i := &IPRateLimiter{
		ips: make(map[string]*LimiterEntry),
		mu:  &sync.RWMutex{},
		r:   r,
		b:   b,
	}

	go func() {
		for range time.Tick(time.Second * 30) {
			i.mu.Lock()
			for s := range i.ips {
				l := i.ips[s]
				var diff = time.Now().Sub(l.lastUse).Seconds()
				if diff > 10 {
					delete(i.ips, s)
				}
			}
			i.mu.Unlock()
		}
	}()

	return i
}

// AddIP creates a new rate limiter and adds it to the ips map,
// using the IP address as the key
func (i *IPRateLimiter) AddIP(ip string) *rate.Limiter {
	i.mu.Lock()
	defer i.mu.Unlock()

	limiter := rate.NewLimiter(i.r, i.b)

	le := &LimiterEntry{
		limiter: limiter,
		lastUse: time.Now(),
	}

	i.ips[ip] = le

	return limiter
}

// GetLimiter returns the rate limiter for the provided IP address if it exists.
// Otherwise calls AddIP to add IP address to the map
func (i *IPRateLimiter) GetLimiter(ip string) *rate.Limiter {
	i.mu.Lock()
	limiter, exists := i.ips[ip]

	if !exists {
		i.mu.Unlock()
		return i.AddIP(ip)
	}

	limiter.lastUse = time.Now()

	i.mu.Unlock()

	return limiter.limiter
}
