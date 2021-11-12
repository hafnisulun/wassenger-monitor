package main

import (
	"time"

	"github.com/go-co-op/gocron"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	s := gocron.NewScheduler(time.UTC)

	status := ""
	// Run every 15 minutes
	s.Cron("*/15 * * * *").Do(func() {
		wassenger := new(Wassenger)
		wassenger.Monitor(&status)
	})

	s.StartBlocking()
}
