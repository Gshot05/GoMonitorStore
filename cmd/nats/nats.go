// cmd/nats/nats.go
package main

import (
	"fmt"
	"log"

	"github.com/nats-io/nats.go"
)

func main() {
	// Подключение к серверу
	natsURL := "localhost"
	log.Println("Подключение...")
	nc, err := nats.Connect(natsURL)
	if err != nil {
		log.Fatal(err)
	} else {
		log.Print("Подключение успешно!")
	}
	defer nc.Close()

	// Подписка на тему
	_, err = nc.Subscribe("log", func(m *nats.Msg) {
		fmt.Printf("Полученное сообщение: %s\n", string(m.Data))
	})
	if err != nil {
		log.Fatalf("Error subscribing to topic: %v", err)
	}

	// Deadlock
	select {}
}
