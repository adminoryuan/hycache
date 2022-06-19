package main

import (
	"fmt"
	"testing"
	"time"
)

func TestC(t *testing.T) {
	ch := NewCache()
	ch.Put("w", "123", 5)
	ch.Put("a", "123", 10)
	ch.Put("c", "123", 11)
	time.Sleep(10 * time.Second)

	a := ch.Get("w")

	fmt.Println(a)
	fmt.Println(ch.Get("a"))
	fmt.Println(ch.Get("c"))

}
