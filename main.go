package main

import (
	"./cache"
	"fmt"
	"time"
)

func main() {
	shortCheckIntervalTime := 1 * time.Second
	shortHoldInMemoryTime := 1 * time.Millisecond
	fastCache := cache.InitNewCache(shortCheckIntervalTime, shortHoldInMemoryTime, true)

	fastCache.Add("some", 42)
	fastCache.Add("another", 100)
	value, err := fastCache.Get("some")
	if err != nil {
		fmt.Println("Error on get from memory cahce", err)
	}
	fmt.Println("Value from memory cache", value)

	// через 1 сек будет выполнена выгрузка на диск
	time.Sleep(2 * time.Second)

	another, err := fastCache.Get("another")
	if err != nil {
		fmt.Println("Error on get from disk", err)
	}
	fmt.Println("Value from disk", another)
}
