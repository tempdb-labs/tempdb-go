package main

// func main() {
// 	client, err := tempdb.NewClient(tempdb.Config{
// 		Addr: "0.0.0.0:8081",
// 		URL:  "tempdb://admin:admin@workspace:8020/ecommerce",
// 	})
// 	if err != nil {
// 		log.Fatalf("Failed to get client: %v", err)
// 	}
// 	defer client.Close()

// 	result, err := client.Query(
// 		tempdb.NewQueryBuilder().
// 			WhereEqual("age_group", "25-45").
// 			Limit(10).
// 			Build(),
// 	)
// 	if err != nil {
// 		return
// 	}

// 	log.Println(result)
// }
