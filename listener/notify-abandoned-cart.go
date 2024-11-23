package listener

import (
	"fmt"
	"listener-connection/config"
	"net/http"
)

func notifyAbandonedCart(cartID string) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/webshop/user-left/%s", config.WebshopApiUrl, cartID), nil)
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
