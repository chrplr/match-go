package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"time"
)

type Item struct {
	Text     string
	SourceID int
}

func main() {
	rand.Seed(time.Now().UnixNano())

	maxConsecutive := flag.Int("n", 2, "Max number of consecutive items from the same file")
	outputFile := flag.String("o", "mixed_list.txt", "Output file path")
	flag.Parse()

	inputFiles := flag.Args()
	if len(inputFiles) < 1 {
		log.Fatal("Usage: mix -n <max_consecutive> -o <output_file> <input_file1> <input_file2> ...")
	}

	buckets := make([][]*Item, len(inputFiles))
	totalItems := 0
	for i, path := range inputFiles {
		file, err := os.Open(path)
		if err != nil {
			log.Fatalf("Failed to open input file %s: %v", path, err)
		}
		defer file.Close()

		reader := bufio.NewReader(file)
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				if err == io.EOF {
					break
				}
				log.Fatalf("Failed to read from %s: %v", path, err)
			}
			buckets[i] = append(buckets[i], &Item{Text: line, SourceID: i})
			totalItems++
		}
	}

	out, err := os.Create(*outputFile)
	if err != nil {
		log.Fatalf("Failed to create output file %s: %v", *outputFile, err)
	}
	defer out.Close()

	var result []*Item
	recentSources := make([]int, 0, *maxConsecutive)

	for len(result) < totalItems {
		var allowedItems []*Item
		allowedSources := make(map[int]bool)

		// Determine allowed sources
		for i := 0; i < len(buckets); i++ {
			if len(buckets[i]) == 0 {
				continue
			}
			consecutiveCount := 0
			for _, sourceID := range recentSources {
				if sourceID == i {
					consecutiveCount++
				}
			}
			if consecutiveCount < *maxConsecutive {
				allowedSources[i] = true
			}
		}
		
		// If no sources are allowed, relax the constraint for one turn
		if len(allowedSources) == 0 {
			for i := 0; i < len(buckets); i++ {
				if len(buckets[i]) > 0 {
					allowedSources[i] = true
				}
			}
		}

		// Pool items from allowed sources
		for sourceID := range allowedSources {
			allowedItems = append(allowedItems, buckets[sourceID]...)
		}

		if len(allowedItems) == 0 {
			break // All buckets are empty
		}

		// Pick a random item from the pool
		pickedIdx := rand.Intn(len(allowedItems))
		pickedItem := allowedItems[pickedIdx]

		// Add to result and update state
		result = append(result, pickedItem)
		recentSources = append(recentSources, pickedItem.SourceID)
		if len(recentSources) > *maxConsecutive {
			recentSources = recentSources[1:]
		}

		// Remove from bucket
		sourceBucket := buckets[pickedItem.SourceID]
		for i, item := range sourceBucket {
			if item == pickedItem {
				buckets[pickedItem.SourceID] = append(sourceBucket[:i], sourceBucket[i+1:]...)
				break
			}
		}
	}

	// Write to file
	for _, item := range result {
		_, err := out.WriteString(item.Text)
		if err != nil {
			log.Fatalf("Failed to write to output file: %v", err)
		}
	}

	fmt.Printf("Successfully created '%s' with %d items.\n", *outputFile, len(result))
}
