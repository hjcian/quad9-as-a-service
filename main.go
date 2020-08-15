package main

import (
	"fmt"
	"q9aas/quad9"
)

func main() {
	q := quad9.CreateQuerier()
	fmt.Println("hello q9ass")
	r, _ := q.IsBlocked("google.com")
	fmt.Printf("google.com %v blocked", func() string {
		if r {
			return "is"
		}
		return "is not"
	}())
}
