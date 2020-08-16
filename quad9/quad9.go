package quad9

import (
	"context"
	"net"
	"strings"
	"sync"
	"time"
)

const (
	// Quad9Sec provides Secured service
	Quad9Sec string = "9.9.9.9:53"
	// Quad9Unc provides Unsecured service
	Quad9Unc string = "9.9.9.10:53"
	// NotExistsSentence will appear in error message
	NotExistsSentence string = "no such host"
)

// Resolver is a DNS query object
type Resolver interface {
	LookupHost(ctx context.Context, host string) (addrs []string, err error)
}

func creatResolver(server string) *net.Resolver {
	return &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{
				Timeout: time.Millisecond * time.Duration(10000),
			}
			return d.DialContext(ctx, "udp", server)
		},
	}
}

// Querier the main object provides IsBlocked functionality
type Querier struct {
	secQuerier Resolver
	reqQuerier Resolver
}

// CreateQuerier the constructor for Querier
func CreateQuerier() *Querier {
	return &Querier{
		creatResolver(Quad9Sec),
		creatResolver(Quad9Unc),
	}
}

type benchResolver struct {
	avgRespTime time.Duration
}

func (br *benchResolver) LookupHost(ctx context.Context, host string) (addrs []string, err error) {
	time.Sleep(br.avgRespTime)
	return []string{"1.2.3.4"}, nil
}

// CreateBenchQuerier the constructor for benchmarking the throughput during different delay of DNS server
func CreateBenchQuerier(d time.Duration) *Querier {
	return &Querier{
		&benchResolver{d},
		&benchResolver{d},
	}
}

func (q *Querier) getProbingResultsAsync(domain string) (error, error) {
	wg := sync.WaitGroup{}
	wg.Add(2)

	errSecChan := make(chan error, 1)
	go func() {
		defer wg.Done()
		_, errSec := q.secQuerier.LookupHost(context.Background(), domain)
		errSecChan <- errSec
	}()

	errUncChan := make(chan error, 1)
	go func() {
		defer wg.Done()
		_, errUnc := q.reqQuerier.LookupHost(context.Background(), domain)
		errUncChan <- errUnc
	}()

	wg.Wait()
	return <-errUncChan, <-errSecChan
}

// IsBlocked is main functionality of this package
func (q *Querier) IsBlocked(domain string) (bool, error) {

	errUnc, errSec := q.getProbingResultsAsync(domain)

	if errUnc != nil {
		// domain is dead possibly
		return false, errUnc
	}

	if errSec != nil && strings.Contains(errSec.Error(), NotExistsSentence) {
		if errUnc == nil {
			return true, nil
		}
		// unsecured service broken
		return false, errUnc
	} else if errSec != nil {
		// secured service broken
		return false, errSec
	}
	// secured service not nil, means this domain is Okay (probably)
	return false, nil
}
