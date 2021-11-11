package main

import (
	"time"

	"github.com/go-co-op/gocron"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	s := gocron.NewScheduler(time.UTC)

	// Run at 22:00 UTC / 05:00 UTC+7
	s.Cron("0 22 * * *").Do(func() {
		wassenger := new(Wassenger)
		wassenger.Monitor(true)
	})

	// Run every minute 10, 30, 50
	s.Cron("10,30,50 * * * *").Do(func() {
		wassenger := new(Wassenger)
		wassenger.Monitor(false)
	})

	s.StartBlocking()
}
