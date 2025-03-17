// package main

// import (
// 	"fmt"
// 	"log"
// 	"time"

// 	tempdb "github.com/tempdb-labs/tempdb-go/lib"
// )

// func main() {
// 	// Initialize the client configuration
// 	client, err := tempdb.NewClient(tempdb.Config{
// 		Addr: "0.0.0.0:8081",
// 		// Addr: "tempdb1.tempdb.xyz:8081",
// 		URL: "tempdb://admin:5M!d$7pG68;5@workspace:890cbc66ca3b/ecommerce",
// 	})
// 	if err != nil {
// 		log.Print(err)
// 	}

// 	defer client.Close()

// 	// === Pub/Sub Examples ===

// 	// Subscribe to a channel and handle incoming messages
// 	// err = client.Subscribe("mychannel", func(message string) {
// 	// 	fmt.Printf("Received message on 'mychannel': %s\n", message)
// 	// })
// 	// if err != nil {
// 	// 	log.Printf("Subscribe error: %v", err)
// 	// }

// 	// // Publish a message to the channel
// 	// count, err := client.Publish("mychannel", "Hello, subscribers!")
// 	// if err != nil {
// 	// 	log.Printf("Publish error: %v", err)
// 	// }
// 	// fmt.Printf("Published to %d subscribers on 'mychannel'\n", count)

// 	// Get the number of subscribers for the channel
// 	numSub, err := client.PubSubNumSub("mychannel")
// 	if err != nil {
// 		log.Printf("PubSubNumSub error: %v", err)
// 	}
// 	fmt.Printf("Number of subscribers to 'mychannel': %d\n", numSub)

// 	// List all active Pub/Sub channels
// 	channels, err := client.PubSubChannels()
// 	if err != nil {
// 		log.Printf("PubSubChannels error: %v", err)
// 	}
// 	fmt.Println("Active Pub/Sub channels:", channels)

// 	// Unsubscribe from the channel
// 	err = client.Unsubscribe("mychannel")
// 	if err != nil {
// 		log.Printf("Unsubscribe error: %v", err)
// 	}
// 	fmt.Println("Unsubscribed from 'mychannel'")

// 	// Wait briefly to ensure subscription messages are processed
// 	time.Sleep(1 * time.Second)

// 	// === Event Stream Examples ===

// 	// Add an event to a stream
// 	id, err := client.XAdd("mystream", map[string]string{"event": "login"})
// 	if err != nil {
// 		log.Printf("XAdd error: %v", err)
// 	}
// 	fmt.Println("Added event to 'mystream' with ID:", id)

// 	// Read events from the stream (starting from the beginning, up to 10 entries)
// 	entries, err := client.XRead("mystream", "-", 10)
// 	if err != nil {
// 		log.Printf("XRead error: %v", err)
// 	}
// 	for i, entry := range entries {
// 		fmt.Printf("Stream entry %d from 'mystream': %v\n", i+1, entry)
// 	}

// 	// List all event streams in the database
// 	streams, err := client.XList()
// 	if err != nil {
// 		log.Printf("XList error: %v", err)
// 	}
// 	fmt.Println("Active event streams:", streams)

// 	// Delete the stream
// 	err = client.XDel("mystream")
// 	if err != nil {
// 		log.Printf("XDel error: %v", err)
// 	}
// 	fmt.Println("Deleted 'mystream'")

// 	// === Message Queue Examples ===

// 	// Enqueue a message into a queue
// 	err = client.Enqueue("myqueue", map[string]string{"task": "process"})
// 	if err != nil {
// 		log.Printf("Enqueue error: %v", err)
// 	}
// 	fmt.Println("Enqueued message to 'myqueue'")

// 	// Peek at the next message without removing it
// 	peekMsg, err := client.QPeek("myqueue")
// 	if err != nil {
// 		log.Printf("QPeek error: %v", err)
// 	}
// 	fmt.Printf("Peeked at next message in 'myqueue': %v\n", peekMsg)

// 	// Get the length of the queue
// 	length, err := client.QLen("myqueue")
// 	if err != nil {
// 		log.Printf("QLen error: %v", err)
// 	}
// 	fmt.Printf("Length of 'myqueue': %d\n", length)

// 	// Dequeue the message from the queue
// 	msg, err := client.Dequeue("myqueue")
// 	if err != nil {
// 		log.Printf("Dequeue error: %v", err)
// 	}
// 	fmt.Printf("Dequeued message from 'myqueue': %v\n", msg)

// 	// List all message queues in the database
// 	queues, err := client.QList()
// 	if err != nil {
// 		log.Printf("QList error: %v", err)
// 	}
// 	fmt.Println("Active message queues:", queues)
// }

package main