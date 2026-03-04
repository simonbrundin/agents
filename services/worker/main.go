package main

import (
	"log"
	"time"
)

func main() {
	log.Println("worker starting (stub)")
	// In a real implementation the worker would connect to Redis, pull a job and execute
	for {
		log.Println("worker idle...")
		time.Sleep(30 * time.Second)
	}
}
