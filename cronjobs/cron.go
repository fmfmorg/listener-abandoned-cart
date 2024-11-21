package cronjobs

import (
	"abandoned-cart-listener/config"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/robfig/cron/v3"
)

var c *cron.Cron

func StartCronJobs() {
	c = cron.New(cron.WithLocation(time.FixedZone("GMT", 0)))

	// Schedule NotifyUpdateData to run every day at 12 AM London time
	_, err := c.AddFunc("0 0 * * *", notifyUpdateData)
	if err != nil {
		log.Fatalf("Could not schedule NotifyUpdateData cron job: %v", err)
	}

	c.Start()
	log.Println("Cron jobs started")
}

func StopCronJobs() {
	ctx := c.Stop()
	select {
	case <-ctx.Done():
		log.Println("Cron jobs stopped")
	case <-time.After(5 * time.Second):
		log.Println("Timed out waiting for cron jobs to stop")
	}
}

func notifyUpdateData() {
	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/webshop/daily-update", config.WebshopApiUrl), nil)
	if err != nil {
		fmt.Println("a: ", err)
		return
	}

	req.Header.Add("X-Request-Source", "SERVER")
	_, err = client.Do(req)
	if err != nil {
		fmt.Println("b: ", err)
		return
	}
}
