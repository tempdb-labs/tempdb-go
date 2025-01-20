package main

import (
	"log"

	tempdb "github.com/tempdb-labs/tempdb-go/lib"
)

func main() {

	client, err := tempdb.NewClient(tempdb.Config{
		Addr: "0.0.0.0:8081",
		URL:  "tempdb://admin:admin@workspace:8020/ecommerce-docs",
	})
	if err != nil {
		log.Fatalf("Failed to get client: %v", err)
	}
	defer client.Close()

	// Insert a document
	doc := map[string]interface{}{
		"name":  "John Doe",
		"age":   30,
		"email": "john@example.com",
	}
	docID, err := client.InsertDoc(doc)
	

	log.Println("document id ", docID)

	// Get the document
	retrievedDoc, err := client.GetDoc(docID)

	log.Println("retrieved doc: ", retrievedDoc)

	// Update the document
	update := map[string]interface{}{
		"age": 31,
	}
	updatedDoc, err := client.UpdateDoc(docID, update)
	log.Println("updatedDoc doc: ", updatedDoc)

	// Query documents
	filter := map[string]interface{}{
		"age": 31,
	}
	matchingDocs, err := client.QueryDocs(filter)
	log.Println("matchingDocs doc: ", matchingDocs)

}
