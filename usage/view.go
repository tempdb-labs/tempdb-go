// package main

// import (
// 	"log"

// 	tempdb "github.com/tempdb-labs/tempdb-go/lib"
// )

// func main() {
// 	client, err := tempdb.NewClient(tempdb.Config{
// 		Addr: "0.0.0.0:8081",
// 		URL:  "tempdb://admin:admin@workspace:8020/ecommerce",
// 	})
// 	if err != nil {
// 		log.Fatalf("Failed to get client: %v", err)
// 	}
// 	defer client.Close()

// 	user := map[string]interface{}{
// 		"name":    "Thembinkosi",
// 		"surname": "Mkhonta",
// 		"preferences": map[string]interface{}{
// 			"mode":          "dark",
// 			"notifications": "no",
// 		},
// 	}

// 	result, err := client.Store("user_01", user)
// 	if err != nil {
// 		return
// 	}

// 	log.Println("user stored successfully %v", result)
// }

package main