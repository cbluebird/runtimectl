package util

import (
	"math/rand"
	"time"
)

func RandomDateString() string {
	// Define the time range
	rand.Seed(time.Now().UnixNano())
	start := time.Date(2024, 7, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2024, 7, 31, 23, 59, 59, 0, time.UTC)

	// Generate a random time within the range
	randomTime := start.Add(time.Duration(rand.Int63n(end.Unix()-start.Unix())) * time.Second)

	// Format the time to the desired string format
	return randomTime.Format("2006-01-02-1504")
}
