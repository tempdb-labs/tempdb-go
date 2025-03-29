// / This package contains document database type commands.
// / for examples you need, simple uncomment and then start using them
package usage

import (
	"log"

	tempdb "github.com/tempdb-labs/tempdb-go/lib"
)

func DocumentExamples(client *tempdb.TempDBClient) {
	// Insert a document
	doc := map[string]any{
		"name":  "John Doe",
		"age":   30,
		"email": "john@example.com",
	}
	docID, err := client.InsertDoc(doc)
	if err != nil {
		log.Println("error inserting document: ", err)
	}
	log.Println("document id ", docID)

	// Get the document
	retrievedDoc, _ := client.GetDoc(docID)

	log.Println("retrieved doc: ", retrievedDoc)

	// Update the document
	update := map[string]any{
		"age": 31,
	}
	updatedDoc, err := client.UpdateDoc(docID, update)
	if err != nil {
		log.Println("failed to update the document: ", err)
		return
	}
	
	log.Println("updatedDoc doc: ", updatedDoc)
}
