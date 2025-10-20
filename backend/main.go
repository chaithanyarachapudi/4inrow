package main

import (
    "log"
    "net/http"
    "time"
)

func main() {
    // Initialize DB (optional warning if failed)
    if err := InitDB(); err != nil {
        log.Println("Warning: DB init failed:", err)
    }

    // Start Kafka consumer for analytics in a goroutine
    StartKafkaConsumer(
        []string{"localhost:29092"}, // brokers
        "analytics-consumer",        // consumer group
        "game-events",               // topic
    )

    // Initialize GameManager
    gm := NewGameManager()

    // WebSocket endpoint for players
    http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
        ServeWS(gm, w, r)
    })

    srv := &http.Server{
        Addr:              ":8080",
        ReadHeaderTimeout: 5 * time.Second,
    }

    log.Println("Server starting on :8080")
    if err := srv.ListenAndServe(); err != nil {
        log.Fatal("Server failed:", err)
    }
}
