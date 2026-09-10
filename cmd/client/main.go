package main

import (
	"flag"
	"fmt"
	"log"
	"main/internal/client"
	"os"
	"strings"
	"time"
)

func main() {
	asyncFlag := flag.Bool("async", false, "async")
	fileFlag := flag.String("file", "", "file with video URLs")
	outputDirFlag := flag.String("out", ".", "output directory")
	timeFlag := flag.Bool("t", false, "show total time")
	videoFlag := flag.String("v", "", "video URLs")
	flag.Parse()

	if *videoFlag != "" && *fileFlag != "" {
		log.Fatal("-v and -file cannot be used together")
	}

	var videoURLs []string

	switch {
	case *videoFlag != "":
		videoURLs = strings.Fields(*videoFlag)

	case *fileFlag != "":
		data, err := os.ReadFile(*fileFlag)
		if err != nil {
			log.Fatalf("Error reading file: %v", err)
		}

		videoURLs = strings.Fields(string(data))

	default:
		videoURLs = flag.Args()
	}

	if len(videoURLs) == 0 {
		fmt.Println("You must specify at least one video URL")
		os.Exit(1)
	}

	if err := os.MkdirAll(*outputDirFlag, os.ModePerm); err != nil {
		log.Fatalf("Error make directory: %v", err)
	}

	const serverAddr = "localhost:50051"
	c, err := client.NewClient(serverAddr)
	if err != nil {
		log.Fatalf("Error make client: %v", err)
	}
	defer c.Close()

	start := time.Now()

	if *asyncFlag {
		results := c.AsyncGetThumbnail(videoURLs)
		for _, result := range results {
			if result.Err != nil {
				log.Printf("Error async get thumbnail: %v", result.Err)
				continue
			}

			if err = client.WriteFile(
				result.Response.VideoID,
				result.Response.Data,
				*outputDirFlag,
			); err != nil {
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

			if err = client.WriteFile(
				response.VideoID,
				response.Data,
				*outputDirFlag,
			); err != nil {
				log.Printf("Error write file: %v", err)
				continue
			}
		}
	}

	if *timeFlag {
		fmt.Printf("Total time: %s\n", time.Since(start))
	}
}
