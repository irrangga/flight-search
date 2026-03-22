package main

import (
	"fmt"
	"log"

	"flight-search/internal/aggregator"
	"flight-search/internal/transport"
)

func main() {
	agg := aggregator.NewAggregator()
	handler := transport.NewHandler(agg)

	r := transport.SetupRouter(handler)

	port := ":8080"
	fmt.Printf("🚀 Flight Search API running on http://localhost%s\n", port)
	if err := r.Run(port); err != nil {
		log.Fatal(err)
	}
}
