package service

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"sync"
	"time"
)

// ExchangeRate holds a single currency pair and its rate.
type ExchangeRate struct {
	From string
	To   string
	Rate float64
}

// RateFetcher is the interface for fetching EUR-based rates from an external source.
type RateFetcher interface {
	FetchRates() (map[string]float64, error)
}

// ExchangeRateService caches exchange rates with a configurable TTL and falls back
// to stale data when the upstream API is unreachable.
type ExchangeRateService struct {
	fetcher   RateFetcher
	ttl       time.Duration
	mu        sync.RWMutex
	rates     map[string]float64 // EUR-based: "USD" -> 1.08 means 1 EUR = 1.08 USD
	fetchedAt time.Time
	hasCache  bool
}

// NewExchangeRateService creates a service backed by the Frankfurter API with the given TTL.
func NewExchangeRateService(ttl time.Duration) *ExchangeRateService {
	return &ExchangeRateService{
		fetcher: &frankfurterFetcher{client: &http.Client{Timeout: 10 * time.Second}},
		ttl:     ttl,
	}
}

// NewExchangeRateServiceWithFetcher creates a service with an injected RateFetcher (for testing).
func NewExchangeRateServiceWithFetcher(fetcher RateFetcher, ttl time.Duration) *ExchangeRateService {
	return &ExchangeRateService{
		fetcher: fetcher,
		ttl:     ttl,
	}
}

// GetRate returns the exchange rate from one currency to another.
// Rates are derived from EUR-based cache; cross-currency rates computed accordingly.
// Falls back to stale cache if the API is unreachable.
func (s *ExchangeRateService) GetRate(from, to string) (float64, error) {
	if from == to {
		return 1.0, nil
	}
	if err := s.ensureCache(); err != nil && !s.hasCache {
		return 0, fmt.Errorf("no exchange rates available: %w", err)
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.convertRate(from, to)
}

// GetAllRates returns all known currency pair rates from the cache.
func (s *ExchangeRateService) GetAllRates() []ExchangeRate {
	_ = s.ensureCache()
	s.mu.RLock()
	defer s.mu.RUnlock()

	currencies := s.knownCurrencies()
	result := make([]ExchangeRate, 0, len(currencies)*len(currencies))
	for _, from := range currencies {
		for _, to := range currencies {
			if from == to {
				continue
			}
			if rate, err := s.convertRate(from, to); err == nil {
				result = append(result, ExchangeRate{From: from, To: to, Rate: rate})
			}
		}
	}
	return result
}

// ensureCache fetches fresh rates if the cache is stale or empty.
// On fetch failure, logs a warning and keeps the existing cache.
func (s *ExchangeRateService) ensureCache() error {
	s.mu.RLock()
	fresh := s.hasCache && time.Since(s.fetchedAt) < s.ttl
	s.mu.RUnlock()
	if fresh {
		return nil
	}

	rates, err := s.fetcher.FetchRates()
	if err != nil {
		slog.Warn("failed to fetch exchange rates, using cached values", "error", err)
		return err
	}

	s.mu.Lock()
	s.rates = rates
	s.fetchedAt = time.Now()
	s.hasCache = true
	s.mu.Unlock()
	return nil
}

// convertRate derives the from→to rate using the EUR-based cache.
// Must be called with at least a read lock held.
func (s *ExchangeRateService) convertRate(from, to string) (float64, error) {
	if from == "EUR" {
		rate, ok := s.rates[to]
		if !ok {
			return 0, fmt.Errorf("unknown currency %q", to)
		}
		return rate, nil
	}
	if to == "EUR" {
		rate, ok := s.rates[from]
		if !ok {
			return 0, fmt.Errorf("unknown currency %q", from)
		}
		return 1.0 / rate, nil
	}
	fromRate, ok := s.rates[from]
	if !ok {
		return 0, fmt.Errorf("unknown currency %q", from)
	}
	toRate, ok := s.rates[to]
	if !ok {
		return 0, fmt.Errorf("unknown currency %q", to)
	}
	return toRate / fromRate, nil
}

// knownCurrencies returns EUR plus all currencies in the cache.
// Must be called with at least a read lock held.
func (s *ExchangeRateService) knownCurrencies() []string {
	currencies := make([]string, 0, len(s.rates)+1)
	currencies = append(currencies, "EUR")
	for k := range s.rates {
		currencies = append(currencies, k)
	}
	return currencies
}

// frankfurterFetcher fetches EUR-based rates from api.frankfurter.app.
type frankfurterFetcher struct {
	client *http.Client
}

func (f *frankfurterFetcher) FetchRates() (map[string]float64, error) {
	resp, err := f.client.Get("https://api.frankfurter.app/latest?from=EUR")
	if err != nil {
		return nil, fmt.Errorf("frankfurter request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response body: %w", err)
	}

	var result struct {
		Rates map[string]float64 `json:"rates"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("parsing frankfurter response: %w", err)
	}
	return result.Rates, nil
}
