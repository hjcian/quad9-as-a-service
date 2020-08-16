package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"q9aas/quad9"
	"time"

	"github.com/bluele/gcache"
	"github.com/gin-gonic/gin"
)

var (
	benchMode      bool
	dnsServerDelay int
	cacheSize      int
	cacheExpiry    int
	q              *quad9.Querier
)

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: go run main.go [options] [root]\n")
	flag.PrintDefaults()
}

func init() {
	flag.BoolVar(&benchMode, "b", false, "Using mock DNS server for benchmark")
	flag.IntVar(&dnsServerDelay, "d", 20, "The DNS server delay in milliseconds (ms)")
	flag.IntVar(&cacheSize, "s", -1, "Cache size. (<=0 means no cache)")
	flag.IntVar(&cacheExpiry, "e", -1, "Cache with an expiration in minutes (min.). (<=0 means forever)")
	flag.Usage = usage
}

func main() {
	flag.Parse()

	if benchMode {
		q = quad9.CreateBenchQuerier(time.Duration(dnsServerDelay) * time.Millisecond)
	} else {
		q = quad9.CreateQuerier()
	}
	cache := gcache.New(cacheSize).
		LRU().
		Build()

	get := func(hostname string) bool {
		// no cahce
		if cacheSize <= 0 {
			r, _ := q.IsBlocked(hostname)
			return r
		}

		// cache hitted
		v, err := cache.Get(hostname)
		if err == nil {
			return v.(bool)
		}

		// cache missed
		r, _ := q.IsBlocked(hostname)
		if cacheExpiry <= 0 {
			// cache forever
			cache.Set(hostname, r)
		} else {
			// cache with expiry
			cache.SetWithExpire(hostname, r, time.Duration(cacheExpiry)*time.Minute)
		}
		return r
	}

	server := gin.Default()
	server.GET("/checkBlocklist", func(c *gin.Context) {
		hostname := c.Query("hostname")

		r := get(hostname)
		c.JSON(http.StatusOK, gin.H{
			"blocked": r,
		})
	})

	server.Run(":12345")
}
