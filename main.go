package main

import (
	"abandoned-cart-listener/cronjobs"
	"abandoned-cart-listener/listener"
	"abandoned-cart-listener/redisclient"

	"os"
	"os/signal"
	"syscall"
)

func main() {
	defer redisclient.Client.Close()

	// Start cron jobs
	cronjobs.StartCronJobs()

	// Start listener
	listener.Start()

	// Set up channel to listen for signals
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// Block until a signal is received
	<-stop

	// Stop cron jobs gracefully
	cronjobs.StopCronJobs()
}
