package main

import (
	"fmt"
	"os"
)

var goal int

func task(c chan int) {
	p := <-c
	if p > goal {
		os.Exit(0)
	}
	fmt.Println(p)
	nc := make(chan int)
	go task(nc)

	for {
		i := <-c
		if i%p != 0 {
			nc <- i
		}
	}
}
