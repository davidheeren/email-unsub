package main

// source: https://gobyexample.com/worker-pools

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"google.golang.org/api/gmail/v1"
)

func downloadTestData(srv *gmail.Service, emailCount int) error {
	resp, err := srv.Users.Messages.List(emailUser).MaxResults(int64(emailCount)).Do()
	if err != nil {
		return fmt.Errorf("unable to retrieve messages: %v", err)
	}

	const numWorkers = 10
	var wg sync.WaitGroup
	jobs := make(chan string)

	for w := 0; w < numWorkers; w++ {
		wg.Add(1)
		go downloadEmailWorker(srv, jobs, &wg)
	}

	for _, m := range resp.Messages {
		jobs <- m.Id
	}

	close(jobs)
	wg.Wait()
	fmt.Println("Downloads complete")
	return nil
}

func downloadEmailWorker(srv *gmail.Service, jobs <-chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	for emailId := range jobs {
		msgFull, err := srv.Users.Messages.Get(emailUser, emailId).Format("full").Do()
		if err != nil {
			fmt.Printf("Could not retrieve full message. continuing... : %v", err)
			continue
		}

		dataFull, err := json.MarshalIndent(msgFull, "", "    ")
		if err != nil {
			fmt.Printf("Could not marshal full message. continuing... : %v", err)
			continue
		}

		msgRaw, err := srv.Users.Messages.Get(emailUser, emailId).Format("raw").Do()
		if err != nil {
			fmt.Printf("Could not retrieve raw message. continuing... : %v", err)
			continue
		}

		dataRaw, err := base64.URLEncoding.DecodeString(msgRaw.Raw)
		if err != nil {
			fmt.Printf("Could not decode raw message. continuing... : %v", err)
			continue
		}

		dataDir := "testdata/uncategorized"
		err = os.MkdirAll(dataDir, 0755)
		if err != nil {
			fmt.Printf("Could not create directory. continuing... : %v", err)
			continue
		}

		err = os.WriteFile(filepath.Join(dataDir, emailId+".json"), dataFull, 0644)
		if err != nil {
			fmt.Printf("Could not write to json file. continuing... : %v", err)
			continue
		}

		err = os.WriteFile(filepath.Join(dataDir, emailId+".eml"), dataRaw, 0644)
		if err != nil {
			fmt.Printf("Could not write to eml file. continuing... : %v", err)
			continue
		}

		fmt.Printf("downloaded %v\n", msgFull.Id)
	}
}
