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
	q              *quad9.Querier
)

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: go run main.go [options] [root]\n")
	flag.PrintDefaults()
}

func init() {
	flag.BoolVar(&benchMode, "b", false, "Using mock DNS server for benchmark")
	flag.IntVar(&dnsServerDelay, "d", 20, "The DNS server delay in milliseconds (ms)")
	flag.Usage = usage

}

func main() {
	flag.Parse()

	if benchMode {
		q = quad9.CreateBenchQuerier(time.Duration(dnsServerDelay) * time.Millisecond)
	} else {
		q = quad9.CreateQuerier()
	}

	server := gin.Default()
	server.GET("/checkBlocklist", func(c *gin.Context) {
		hostname := c.Query("hostname")

		r, _ := q.IsBlocked(hostname)
		c.JSON(http.StatusOK, gin.H{
			"blocked": r,
		})
	})
	server.Run(":12345")
}
