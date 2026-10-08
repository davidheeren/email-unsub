package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"

	"google.golang.org/api/gmail/v1"
)


func runTestData() {
	pattern := filepath.Join("testdata", "unsub", "*.json")
	paths, err := filepath.Glob(pattern)
	if err != nil {
		log.Fatal(err)
	}

	const numWorkers = 50
	var wg sync.WaitGroup
	jobs := make(chan string)

	for w := 0; w < numWorkers; w++ {
		wg.Add(1)
		go runTestDataWorker(jobs, &wg)
	}

	for _, p := range paths {
		jobs <- p
	}

	close(jobs)
	wg.Wait()
	fmt.Println("testing complete")
}

func runTestDataWorker(jobs <-chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	for path := range jobs {
		file, err := os.Open(path)
		if err != nil {
			fmt.Printf("there was an issue opening %v . continuing : %v", path, err)
			continue
		}

		defer file.Close()
		decoder := json.NewDecoder(file)
		var msg gmail.Message
		err = decoder.Decode(&msg)
		if err != nil {
			fmt.Printf("there was an issue decoding %v . continuing : %v", path, err)
			continue
		}

		unsubUrl, err := scanUnsubUrls(&msg)
		if err != nil {
			fmt.Printf("there was an eror scanning %v . test failed : %v", path, err)
			continue
		}
		if unsubUrl == "" {
			fmt.Printf("no unsub url when scanning %v . test failed\n", path)
			continue
		}
	}
}
