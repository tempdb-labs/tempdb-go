// / This package contains key value database type commands.
// / for examples you need, simple uncomment and then start using them
// Note, these commands are to get you started
package usage

import (
	"log"

	tempdb "github.com/tempdb-labs/tempdb-go/lib"
)

func KeyValueCommands(client *tempdb.TempDBClient) {
	// insert a key-value into the database
	res, err := client.Store("user-one", map[string]any{
		"name": "Bruce", "surname": "Wayne",
	})
	if err != nil {
		log.Println("failed to store value: ", err)
		return
	}

	log.Println("result after storing user: ", res)

	// getting the user from the database
	user, err := client.Get("user-one")
	if err != nil {
		log.Println("failed to get user: ", err)
		return
	}
	log.Println("retrieved user details: ", user)

	// deleting the user using the key
	_, err = client.Delete("user-one")
	if err != nil {
		log.Println("failed to delete user: ", err)
		return
	}

	// setting a value with custom ttl
	res2, err := client.SetEx("user-two", 3600, map[string]interface{}{
		"testing": "setting value with ttl",
	})
	if err != nil {
		log.Println("failed to store value: ", err)
		return
	}

	log.Println("value store: ", res2)

}
