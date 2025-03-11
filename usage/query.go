package main

import (
	"fmt"
	"log"

	tempdb "github.com/tempdb-labs/tempdb-go/lib"
)

func main() {
	client, err := tempdb.NewClient(tempdb.Config{
		Addr: "0.0.0.0:8081",
		URL:  "tempdb://admin:Q{)6X!mG[hTK@workspace:cb4552273c5c/ecommerce",
	})
	if err != nil {
		log.Fatalf("Failed to get client: %v", err)
	}
	defer client.Close()

	// Example 1: Using raw command string
	// Get total sales by payment method
	result1, err := client.Query("GROUPBY /payment_method SUM /net_amount")
	if err != nil {
		log.Fatalf("error: %v", result1)
	}
	fmt.Printf("Sales by payment method and sum: %v\n", result1)

	// pipeline := tempdb.NewAggregation().Count().Build()
	// if err != nil {
	// 	log.Println(err)
	// }
	// log.Println(pipeline)

	// Example 2: usign query builder
	// Get average purchase amount by age group for female customers
	// pipeline := tempdb.NewAggregation().Filter("gender", "eq", "Female").GroupBy("age_group").Average("net_amount")

	// result2, err := client.AggregateWithBuilder(pipeline)
	// if err != nil {
	// 	log.Fatalf("error: %v", result2)
	// }

	// fmt.Printf("average purchase by age group (females): %v\n", result2)

	// // Example 3: Complex analysis
	// // Get count and toal sales by location where discount was used
	// pipeline2 := tempdb.NewAggregation().Filter("discount_availed", "eq", "Yes").GroupBy("location").Count().Sum("new_amount")
	// result3, err := client.AggregateWithBuilder(pipeline2)
	// if err != nil {
	// 	log.Fatalf("error: %v", err)
	// }

	// fmt.Printf("sales analysis by location: %v\n", result3)

	// Example 4: Customer behaviour analysis
	// get the average purchase amount by paymenet method and gender
	// behaviourPipe := tempdb.NewQuery().GroupBy("payment_method").GroupBy("gender").Average("net_amount")
	// result4, err := client.QueryWithBuilder(behaviourPipe)
	// if err != nil {
	// 	log.Fatalf("error: %v", result4)
	// }

	// log.Println(behaviourPipe)

	// // Exampe 5: time based analysis
	// timePipe := tempdb.NewAggregation().Filter("net_amount", "gt", "1000").GroupBy("age_group").Count()
	// result5, err := client.AggregateWithBuilder(timePipe)
	// if err != nil {
	// 	log.Fatalf("error: %f", result5)
	// }
}
