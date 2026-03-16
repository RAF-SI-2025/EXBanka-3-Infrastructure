package service_test

import (
	"errors"
	"testing"
	"time"

	"github.com/RAF-SI-2025/EXBanka-3-Infrastructure/internal/service"
)

// mockFetcher is a test double for RateFetcher.
type mockFetcher struct {
	rates map[string]float64
	err   error
	calls int
}

func (m *mockFetcher) FetchRates() (map[string]float64, error) {
	m.calls++
	return m.rates, m.err
}

func newTestService(rates map[string]float64, ttl time.Duration) (*service.ExchangeRateService, *mockFetcher) {
	mock := &mockFetcher{rates: rates}
	svc := service.NewExchangeRateServiceWithFetcher(mock, ttl)
	return svc, mock
}

func TestGetRate_EURtoUSD(t *testing.T) {
	svc, _ := newTestService(map[string]float64{"USD": 1.08, "RSD": 117.0}, time.Hour)

	rate, err := svc.GetRate("EUR", "USD")
	if err != nil {
		t.Fatalf("GetRate(EUR, USD) error: %v", err)
	}
	if rate != 1.08 {
		t.Errorf("GetRate(EUR, USD) = %v, want 1.08", rate)
	}
}

func TestGetRate_USDtoEUR(t *testing.T) {
	svc, _ := newTestService(map[string]float64{"USD": 2.0}, time.Hour)

	rate, err := svc.GetRate("USD", "EUR")
	if err != nil {
		t.Fatalf("GetRate(USD, EUR) error: %v", err)
	}
	if rate != 0.5 {
		t.Errorf("GetRate(USD, EUR) = %v, want 0.5", rate)
	}
}

func TestGetRate_CrossCurrency(t *testing.T) {
	// USD→GBP = rates[GBP] / rates[USD] = 1.0 / 2.0 = 0.5
	svc, _ := newTestService(map[string]float64{"USD": 2.0, "GBP": 1.0}, time.Hour)

	rate, err := svc.GetRate("USD", "GBP")
	if err != nil {
		t.Fatalf("GetRate(USD, GBP) error: %v", err)
	}
	if rate != 0.5 {
		t.Errorf("GetRate(USD, GBP) = %v, want 0.5", rate)
	}
}

func TestGetRate_SameCurrency(t *testing.T) {
	svc, _ := newTestService(map[string]float64{"USD": 1.08}, time.Hour)

	rate, err := svc.GetRate("USD", "USD")
	if err != nil {
		t.Fatalf("GetRate(USD, USD) error: %v", err)
	}
	if rate != 1.0 {
		t.Errorf("GetRate(USD, USD) = %v, want 1.0", rate)
	}
}

func TestGetRate_UnknownCurrency_ReturnsError(t *testing.T) {
	svc, _ := newTestService(map[string]float64{"USD": 1.08}, time.Hour)

	_, err := svc.GetRate("EUR", "XYZ")
	if err == nil {
		t.Error("GetRate(EUR, XYZ) expected error for unknown currency, got nil")
	}
}

func TestGetRate_UsesCache(t *testing.T) {
	svc, mock := newTestService(map[string]float64{"USD": 1.08}, time.Hour)

	// Two calls within TTL — fetcher should only be called once
	svc.GetRate("EUR", "USD") //nolint:errcheck
	svc.GetRate("EUR", "USD") //nolint:errcheck

	if mock.calls != 1 {
		t.Errorf("FetchRates called %d times, want 1 (second call should use cache)", mock.calls)
	}
}

func TestGetRate_CacheExpiry_TriggersRefetch(t *testing.T) {
	svc, mock := newTestService(map[string]float64{"USD": 1.08}, time.Millisecond)

	svc.GetRate("EUR", "USD") //nolint:errcheck
	time.Sleep(2 * time.Millisecond)
	svc.GetRate("EUR", "USD") //nolint:errcheck

	if mock.calls < 2 {
		t.Errorf("FetchRates called %d times after cache expiry, want >= 2", mock.calls)
	}
}

func TestGetRate_FallbackOnError(t *testing.T) {
	mock := &mockFetcher{rates: map[string]float64{"USD": 1.08}}
	svc := service.NewExchangeRateServiceWithFetcher(mock, time.Millisecond)

	// Prime the cache
	svc.GetRate("EUR", "USD") //nolint:errcheck

	// Make fetcher fail
	mock.err = errors.New("API down")
	mock.rates = nil
	time.Sleep(2 * time.Millisecond)

	// Should still return stale cached rate (fallback)
	rate, err := svc.GetRate("EUR", "USD")
	if err != nil {
		t.Fatalf("GetRate with stale cache fallback returned error: %v", err)
	}
	if rate != 1.08 {
		t.Errorf("GetRate fallback = %v, want 1.08", rate)
	}
}

func TestGetRate_NoCache_FetchError_ReturnsError(t *testing.T) {
	mock := &mockFetcher{err: errors.New("API unavailable")}
	svc := service.NewExchangeRateServiceWithFetcher(mock, time.Hour)

	_, err := svc.GetRate("EUR", "USD")
	if err == nil {
		t.Error("GetRate with no cache and fetch error expected error, got nil")
	}
}

func TestGetAllRates_ReturnsNonEmpty(t *testing.T) {
	svc, _ := newTestService(map[string]float64{"USD": 1.08, "RSD": 117.0}, time.Hour)

	rates := svc.GetAllRates()
	if len(rates) == 0 {
		t.Error("GetAllRates() returned empty slice, want non-empty")
	}
}

func TestGetAllRates_ContainsExpectedPairs(t *testing.T) {
	svc, _ := newTestService(map[string]float64{"USD": 1.08}, time.Hour)

	rates := svc.GetAllRates()
	found := false
	for _, r := range rates {
		if r.From == "EUR" && r.To == "USD" {
			found = true
			break
		}
	}
	if !found {
		t.Error("GetAllRates() missing EUR→USD pair")
	}
}
