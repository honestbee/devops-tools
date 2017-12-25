package main

import (
	"fmt"
	"log"

	cloudflare "github.com/cloudflare/cloudflare-go"
)

func clearCache(apiKey string, email string, files []string, domain string, url string) error {
	for i := range files {
		files[i] = "https://" + url + files[i]
	}
	purgeCacheRequest := cloudflare.PurgeCacheRequest{
		Files: files,
	}

	// Construct a new API object
	api, err := cloudflare.New(apiKey, email)
	if err != nil {
		log.Fatal(err)
	}

	// Fetch the zone ID
	id, err := api.ZoneIDByName(domain) // Assuming example.com exists in your Cloudflare account already
	if err != nil {
		log.Fatal(err)
	}

	purgeCacheRespone, err := api.PurgeCache(id, purgeCacheRequest)
	if err != nil {
		log.Fatal(err)
	}

	_, err = fmt.Printf("status: %v", purgeCacheRespone.Response.Success)
	return err
}
