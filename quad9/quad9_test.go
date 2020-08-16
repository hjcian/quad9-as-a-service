package quad9

import (
	"context"
	"fmt"
	"testing"
	"time"
)

const (
	GoodOne      string = "google.com"
	BadOne       string = "bad.com"
	NotExistsOne string = "foobar.blah"
)

type mockSecResolver struct{}

func (mr *mockSecResolver) LookupHost(ctx context.Context, host string) (addrs []string, err error) {
	time.Sleep(50 * time.Millisecond)
	switch host {
	case GoodOne:
		return []string{"1.1.1.1"}, nil
	case BadOne:
		// Secured service return 'NXDOMAIN' response if a site is blocked
		return []string{}, fmt.Errorf(NotExistsSentence)
	case NotExistsOne:
		return []string{}, fmt.Errorf(NotExistsSentence)
	}
	return nil, nil
}

type mockUncResolver struct{}

func (mr *mockUncResolver) LookupHost(ctx context.Context, host string) (addrs []string, err error) {
	time.Sleep(50 * time.Millisecond)
	switch host {
	case GoodOne:
		return []string{"1.1.1.1"}, nil
	case BadOne:
		// Unsecured service return address response if a site is exists
		return []string{"1.1.1.1"}, nil
	case NotExistsOne:
		return []string{}, fmt.Errorf(NotExistsSentence)
	}
	return nil, nil
}
func Test_cache_getter(t *testing.T) {
	q := CreateBenchQuerier(time.Duration(100) * time.Millisecond)

	getTime := func(host string, cacheGetter CacheGetter) time.Duration {
		start := time.Now()
		cacheGetter(host)
		diff := time.Since(start)
		return diff
	}

	t.Run("Test basic cache", func(t *testing.T) {
		cacheSize := 2
		cacheExpiry := -1
		cacheGetter := q.NewCacheGetter(cacheSize, cacheExpiry)

		host := "aaa.com"
		diff1 := getTime(host, cacheGetter)
		t.Logf("[First] took %v (%v) \n", diff1, diff1.Milliseconds())

		diff2 := getTime(host, cacheGetter)
		t.Logf("[Second] took %v (%v) \n", diff2, diff2.Milliseconds())

		if diff2.Milliseconds() > 100 {
			t.Errorf("%v should tool below 100ms but consumed %v \n", host, diff2)
		}
	})

	t.Run("Test cache is expirable", func(t *testing.T) {
		cacheSize := 2
		cacheExpiry := 1 // in seconds
		cacheGetter := q.NewCacheGetter(cacheSize, cacheExpiry)

		host := "a.com"

		diff1 := getTime(host, cacheGetter)
		t.Logf("[First] took %v (%v) \n", diff1, diff1.Milliseconds())

		// sleep 1 second for expiration
		time.Sleep(1 * time.Second)

		diff2 := getTime(host, cacheGetter)
		t.Logf("[Second] took %v (%v) \n", diff2, diff2.Milliseconds())

		if diff2.Milliseconds() < 100 {
			t.Errorf("%v should tool longer than 100ms but consumed %v \n", host, diff2)
		}
	})
}

func Test_main(t *testing.T) {
	q := Querier{
		secQuerier: &mockSecResolver{},
		reqQuerier: &mockUncResolver{},
	}

	t.Run("Test Exists", func(t *testing.T) {
		testSuits := []struct {
			name      string
			domain    string
			isBlocked bool
		}{
			{"Good", GoodOne, false},
			{"Bad", BadOne, true},
		}
		for _, test := range testSuits {
			t.Logf("[%v] testing...\n", test.name)
			got, err := q.IsBlocked(test.domain)
			if err != nil {
				t.Errorf("[%v] expect no error, got one: '%v' \n", test.name, err)
			}
			if got != test.isBlocked {
				t.Errorf("[%v] expect '%v', got '%v' \n", test.name, test.isBlocked, got)
			}
		}
	})

	t.Run("Test Not Exists", func(t *testing.T) {
		testSuits := []struct {
			name      string
			domain    string
			isBlocked bool
		}{
			{"Not Exists", NotExistsOne, false},
		}
		for _, test := range testSuits {
			t.Logf("[%v] testing...\n", test.name)
			_, err := q.IsBlocked(test.domain)
			if err == nil {
				t.Errorf("[%v] expect err not nil (%v) \n", test.name, err)
			}
		}
	})
}

func (q *Querier) getProbingResultsSync(domain string) (error, error) {
	_, errSec := q.secQuerier.LookupHost(context.Background(), domain)
	_, errUnc := q.reqQuerier.LookupHost(context.Background(), domain)
	return errSec, errUnc
}

func Benchmark_getProbingResults(b *testing.B) {
	q := Querier{
		secQuerier: &mockSecResolver{},
		reqQuerier: &mockUncResolver{},
	}
	b.Run("getProbingResultsSync", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			q.getProbingResultsSync(GoodOne)
		}
	})
	b.Run("getProbingResultsAsync", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			q.getProbingResultsAsync(GoodOne)
		}

	})
}
