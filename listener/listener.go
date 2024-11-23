package listener

import (
	"encoding/json"
	"listener-abandoned-cart/config"
	"listener-abandoned-cart/models"
	"listener-abandoned-cart/redisclient"
	"log"
	"time"
)

var (
	heartbeatMap = make(map[string]time.Time)
)

func Start() {
	// Subscribe to Redis channels
	pubsub := redisclient.Client.Subscribe(
		redisclient.Ctx,
		config.ConnectChannel,
		config.DisconnectChannel,
		config.HeartbeatChannel,
		config.PaymentStartChannel,
		config.PaymentEndChannel,
	)
	defer pubsub.Close()

	ch := pubsub.Channel()

	go monitorHeartbeats()

	for msg := range ch {
		switch msg.Channel {
		case config.ConnectChannel:
			handleConnect(msg.Payload)
		case config.DisconnectChannel:
			handleDisconnect(msg.Payload)
		case config.HeartbeatChannel:
			// Update the heartbeat timestamp for the WebSocket server
			serverID := msg.Payload
			heartbeatMap[serverID] = time.Now()
		case config.PaymentStartChannel:
			annountPaymentUpdate(msg.Payload, "start")
		case config.PaymentEndChannel:
			annountPaymentUpdate(msg.Payload, "fail")
		}
	}
}

func annountPaymentUpdate(cartID, status string) {
	message := models.PaymentUpdateWsMessage{
		CartID:  cartID,
		Payload: status, // 3 statuses: start, success, fail
	}
	messages := []models.PaymentUpdateWsMessage{message}
	msg, err := json.Marshal(messages)
	if err != nil {
		return
	}
	redisclient.Client.Publish(redisclient.Ctx, config.RedisPaymentUpdateChannel, msg)
}

func handleConnect(cartID string) {
	// redisclient.Client.Incr(redisclient.Ctx, cartID)
	count, err := redisclient.Client.HIncrBy(redisclient.Ctx, config.HashKeyAllConnections, cartID, 1).Result()
	if err != nil {
		log.Printf("Error decrementing connection count: %v", err)
		return
	}

	if count == 1 {
		notifyUserReturned(cartID)
	}
}

func handleDisconnect(cartID string) {
	count, err := redisclient.Client.HIncrBy(redisclient.Ctx, config.HashKeyAllConnections, cartID, -1).Result()
	if err != nil {
		log.Printf("Error decrementing connection count: %v", err)
		return
	}

	if count <= 0 {
		redisclient.Client.HDel(redisclient.Ctx, config.HashKeyAllConnections, cartID)
		notifyAbandonedCart(cartID)
	}
}

func monitorHeartbeats() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		now := time.Now()
		for serverID, lastHeartbeat := range heartbeatMap {
			if now.Sub(lastHeartbeat) > 30*time.Second {
				log.Printf("WebSocket server %s has not sent a heartbeat in over 30 seconds. Assuming it has shut down.", serverID)
				// Handle the lost connections for this server
				handleServerShutdown(serverID)
				delete(heartbeatMap, serverID)
			}
		}
	}
}

func handleServerShutdown(serverID string) {
	// Get the list of cart IDs associated with the server that shut down
	cartIDs, err := redisclient.Client.HKeys(redisclient.Ctx, serverID).Result()
	if err != nil {
		log.Printf("Error getting cart IDs for server %s: %v", serverID, err)
		return
	}

	// Decrement the connection count for each cart ID and notify if abandoned
	for _, cartID := range cartIDs {
		count, _ := redisclient.Client.HGet(redisclient.Ctx, serverID, cartID).Int64()
		newCount, _ := redisclient.Client.HIncrBy(redisclient.Ctx, config.HashKeyAllConnections, cartID, count*-1).Result()
		if newCount <= 0 {
			redisclient.Client.HDel(redisclient.Ctx, config.HashKeyAllConnections, cartID)
			notifyAbandonedCart(cartID)
		}
	}

	// Clean up the server's entry in Redis
	redisclient.Client.Del(redisclient.Ctx, serverID)
}
