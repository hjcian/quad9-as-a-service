package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"q9aas/quad9"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	benchMode      bool
	dnsServerDelay int
	cacheSize      int
	cacheExpiry    int
)

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: go run main.go [options] [root]\n")
	flag.PrintDefaults()
}

func init() {
	flag.BoolVar(&benchMode, "b", false, "Using mock DNS server for benchmark")
	flag.IntVar(&dnsServerDelay, "d", 20, "The DNS server delay in milliseconds (ms)")
	flag.IntVar(&cacheSize, "s", -1, "Cache size. (<=0 means no cache)")
	flag.IntVar(&cacheExpiry, "e", -1, "Cache with an expiration in seconds (sec.). (<=0 means forever)")
	flag.Usage = usage
}

func main() {
	flag.Parse()

	var q *quad9.Querier
	if benchMode {
		q = quad9.CreateBenchQuerier(time.Duration(dnsServerDelay) * time.Millisecond)
	} else {
		q = quad9.CreateQuerier()
	}
	cacheGetter := q.NewCacheGetter(cacheSize, cacheExpiry)

	server := gin.Default()
	server.GET("/checkBlocklist", func(c *gin.Context) {
		hostname := c.Query("hostname")

		r := cacheGetter(hostname)
		c.JSON(http.StatusOK, gin.H{
			"blocked": r,
		})
	})

	server.Run(":12345")
}
