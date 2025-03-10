package main

import (
	"fmt"
	"log"
	"time"

	"github.com/tempdb-labs/tempdb-go/lib"
)

func main() {
	config := lib.Config{
		Addr: "0.0.0.0:8081",
		URL:  "<your-url>",
	}

	// Run each service in a goroutine
	go RunChatService(config)
	go RunEventLogger(config)
	go RunTaskProcessor(config)

	// Keep main running
	log.Println("Microservices started. Press Ctrl+C to stop.")
	select {}
}

func RunChatService(config lib.Config) {
	client, err := lib.NewClient(config)
	if err != nil {
		log.Fatalf("Chat service failed to connect: %v", err)
	}
	defer client.Close()

	// Subscribe to the chat channel
	err = client.Subscribe("chatroom", func(message string) {
		fmt.Printf("[Chat] Received: %s\n", message)
	})
	if err != nil {
		log.Printf("Chat subscribe error: %v", err)
		return
	}

	// Publish messages periodically
	go func() {
		for i := 1; i <= 5; i++ {
			msg := fmt.Sprintf("User%d says: Hello!", i)
			count, err := client.Publish("chatroom", msg)
			if err != nil {
				log.Printf("Publish error: %v", err)
			} else {
				fmt.Printf("[Chat] Sent '%s' to %d subscribers\n", msg, count)
			}
			time.Sleep(1 * time.Second)
		}
	}()

	// Keep the service running
	select {}
}

func RunEventLogger(config lib.Config) {
	client, err := lib.NewClient(config)
	if err != nil {
		log.Fatalf("Event logger failed to connect: %v", err)
	}
	defer client.Close()

	// Log some events
	events := []map[string]string{
		{"event": "user_login", "user": "alice"},
		{"event": "purchase", "item": "book"},
		{"event": "logout", "user": "alice"},
	}
	for _, event := range events {
		id, err := client.XAdd("user_events", event)
		if err != nil {
			log.Printf("XAdd error: %v", err)
		} else {
			fmt.Printf("[Logger] Added event with ID: %s\n", id)
		}
		time.Sleep(500 * time.Millisecond)
	}

	// Read and display events
	entries, err := client.XRead("user_events", "-", 10)
	if err != nil {
		log.Printf("XRead error: %v", err)
	} else {
		fmt.Println("[Logger] All events:")
		for _, entry := range entries {
			fmt.Printf("  %v\n", entry)
		}
	}

	// Keep running to allow manual testing
	select {}
}

func RunTaskProcessor(config lib.Config) {
	client, err := lib.NewClient(config)
	if err != nil {
		log.Fatalf("Task processor failed to connect: %v", err)
	}
	defer client.Close()

	// Enqueue some tasks
	tasks := []map[string]string{
		{"task": "send_email", "to": "alice@example.com"},
		{"task": "process_payment", "amount": "100"},
		{"task": "generate_report", "type": "daily"},
	}
	for _, task := range tasks {
		err := client.Enqueue("task_queue", task)
		if err != nil {
			log.Printf("Enqueue error: %v", err)
		} else {
			fmt.Printf("[Processor] Enqueued: %v\n", task)
		}
		time.Sleep(500 * time.Millisecond)
	}

	// Process tasks
	go func() {
		for {
			task, err := client.Dequeue("task_queue")
			if err != nil {
				if err.Error() != "EMPTY" {
					log.Printf("Dequeue error: %v", err)
				}
				time.Sleep(1 * time.Second)
				continue
			}
			fmt.Printf("[Processor] Processing task: %v\n", task)
			time.Sleep(2 * time.Second) // Simulate work
		}
	}()

	// Keep running
	select {}
}
