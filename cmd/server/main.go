package main

import (
	"flag"
	"promotions-service/api/v1/router"
)

func main() {
	port := flag.String("port", "8080", "Port for test HTTP server")
	flag.Parse()

	router.Serve(*port)
}
