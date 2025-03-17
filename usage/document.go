// / This package contains document database type commands.
// / for examples you need, simple uncomment and then start using them
package main

// func main() {

// 	client, err := tempdb.NewClient(tempdb.Config{
// 		Addr: "0.0.0.0:8081",
// 		// Addr: "tempdb1.tempdb.xyz:8081",
// 		URL: "tempdb://admin:5M!d$7pG68;5@workspace:890cbc66ca3b/ecommerce",
// 	})

// 	if err != nil {
// 		log.Fatalf("Failed to get client: %v", err)
// 	}
// 	defer client.Close()

// 	// Insert a document
// 	doc := map[string]interface{}{
// 		"name":  "John Doe",
// 		"age":   30,
// 		"email": "john@example.com",
// 	}
// 	docID, err := client.InsertDoc(doc)

// 	log.Println("document id ", docID)

// 	// Get the document
// 	retrievedDoc, err := client.GetDoc(docID)

// 	log.Println("retrieved doc: ", retrievedDoc)

// 	// Update the document
// 	update := map[string]interface{}{
// 		"age": 31,
// 	}
// 	updatedDoc, err := client.UpdateDoc(docID, update)
// 	log.Println("updatedDoc doc: ", updatedDoc)
// }
