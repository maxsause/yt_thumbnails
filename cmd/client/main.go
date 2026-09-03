package main

import (
	"flag"
	"fmt"
	"log"
	"main/internal/client"
	"os"
)

func main() {
	async := flag.Bool("async", false, "async")
	outputDir := flag.String("out", ".", "output directory")
	flag.Parse()

	videoURLs := flag.Args()
	if len(videoURLs) == 0 {
		fmt.Println("You must specify at least one video URL")
		os.Exit(1)
	}

	if err := os.MkdirAll(*outputDir, os.ModePerm); err != nil {
		log.Fatalf("Error make directory: %v", err)
	}

	const serverAddr = "localhost:50051"
	c, err := client.NewClient(serverAddr)
	if err != nil {
		log.Fatalf("Error make client: %v", err)
	}
	defer c.Close()

	if *async {
		results := c.AsyncGetThumbnail(videoURLs)
		for _, result := range results {
			if result.Err != nil {
				log.Printf("Error async get thumbnail: %v", result.Err)
				continue
			}
			if err = client.WriteFile(result.Response.VideoID, result.Response.Data, *outputDir); err != nil {
				log.Printf("Error writing file: %v", err)
				continue
			}
		}

	} else {
		for _, videoURL := range videoURLs {
			response, err := c.GetThumbnail(videoURL)
			if err != nil {
				log.Printf("Error get thumbnail: %v", err)
				continue
			}
			if err = client.WriteFile(response.VideoID, response.Data, *outputDir); err != nil {
				log.Printf("Error write file: %v", err)
				continue
			}
		}
	}
}
