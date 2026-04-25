package main

import (
	"log"

	"github.com/ghassenelkamel/portfolio-v2/internal/config"
	"github.com/ghassenelkamel/portfolio-v2/internal/httpapi"
)

func main() {
	cfg := config.Load()
	srv := httpapi.NewServer(cfg)
	log.Fatal(srv.ListenAndServe())
}
