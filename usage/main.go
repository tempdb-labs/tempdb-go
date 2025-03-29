// This package is for you to test out some of the commands using examples in /usage/examples folder
package main

import (
	"log"

	tempdb "github.com/tempdb-labs/tempdb-go/lib"
	usage "github.com/tempdb-labs/tempdb-go/usage/examples"
)

func main() {
	client, err := tempdb.NewClient(tempdb.Config{
		// tempdb://<username>:<password>@<database-name>:<type>
		URL:  "<enter your url here. You can also create one at tempdb.xyz>",
		Addr: "0.0.0.0:8081",
	})
	if err != nil {
		log.Println("failed to initialise client: ", err)
	}
	defer client.Close()

	// based on which type you created, run any of the following by commenting out
	// assuming key-value by default

	// 1. Key value commands
	usage.KeyValueCommands(client)

	// 2. Document commands(assuming you will create a database type for document dbs)
	usage.DocumentExamples(client)

	// 3. Messaging
	usage.MessagingCommands(client)

	// 4. Querying commands
	usage.QueryCommand(client)

}
