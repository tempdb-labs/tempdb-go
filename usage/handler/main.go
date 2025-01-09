package main

import (
	"log"
	"net/http"

	"github.com/ThembinkosiThemba/zen"
	tempdb "github.com/tempdb-labs/tempdb-go/lib"
)

func main() {
	app := zen.New()

	client, err := tempdb.NewClient(tempdb.Config{
		Addr: "0.0.0.0:8081",
		URL:  "tempdb://admin:admin@workspace:8020/ecommerce",
	})
	if err != nil {
		log.Fatalf("Failed to get client: %v", err)
	}
	defer client.Close()

	app.Apply(zen.Logger())
	Routes(app, client)

	app.Serve(":8082")

}

func Routes(r *zen.Engine, client *tempdb.TempDBClient) {
	r.GET("/", func(ctx *zen.Context) {
		res, err := client.Query("GROUPBY /payment_method SUM /net_amount")
		if err != nil {
			ctx.Error(http.StatusInternalServerError, err.Error())
			return
		}

		ctx.Success(http.StatusOK, res, "OK")
	})

	r.GET("/1", func(ctx *zen.Context) {
		// Example 2: usign query builder
		// Get average purchase amount by age group for female customers
		pipeline := tempdb.NewQuery().Filter("gender", "eq", "Female").GroupBy("age_group").Average("net_amount")

		result2, err := client.QueryWithBuilder(pipeline)
		if err != nil {
			ctx.Error(http.StatusInternalServerError, err.Error())
			return
		}

		ctx.Success(http.StatusOK, result2, "OK")
	})

	r.GET("/2", func(ctx *zen.Context) {
		// // Example 3: Complex analysis
		// // Get count and toal sales by location where discount was used
		pipeline2 := tempdb.NewQuery().Filter("discount_availed", "eq", "Yes").GroupBy("location").Count().Sum("new_amount")
		result3, err := client.QueryWithBuilder(pipeline2)
		if err != nil {
			zen.Fatalf("error: %v", err)
		}
		ctx.Success(http.StatusOK, result3, "OK")
	})

	r.GET("/3", func(ctx *zen.Context) {
		// Example 4: Customer behaviour analysis
		// get the average purchase amount by paymenet method and gender
		behaviourPipe := tempdb.NewQuery().GroupBy("payment_method").GroupBy("gender").Average("net_amount")
		result4, err := client.QueryWithBuilder(behaviourPipe)
		if err != nil {
			log.Fatalf("error: %v", result4)
		}

		ctx.Success(http.StatusOK, result4, "OK")
	})

	r.GET("/4", func(ctx *zen.Context) {
		// Exampe 5: time based analysis
		timePipe := tempdb.NewQuery().Filter("net_amount", "gt", "1000").GroupBy("age_group").Count()
		result5, err := client.QueryWithBuilder(timePipe)
		if err != nil {
			log.Fatalf("error: %f", result5)
		}

		ctx.Success(http.StatusOK, result5, "OK")
	})

	r.GET("/count", func(ctx *zen.Context) {
		// countPipe := tempdb.NewQuery().Count()
		// result, err := client.QueryWithBuilder(countPipe)
		result, err := client.Query("COUNT")
		if err != nil {
			ctx.Error(http.StatusInternalServerError, err.Error())
			return
		}

		ctx.Success(http.StatusOK, result, "OK")
	})

	r.GET("/logs", func(ctx *zen.Context) {
		result, err := client.ViewLogs()
		if err != nil {
			ctx.Error(http.StatusInternalServerError, err.Error())
			return
		}

		ctx.Success(http.StatusOK, result, "OK")
	})
}
