package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	sanitize "github.com/mrz1836/go-sanitize"
	"github.com/tidwall/gjson"
)

const (
	bingSearchURL = "https://api.bing.microsoft.com/v7.0/images/search?q=%s"
)

func main() {
	apiKey := os.Getenv("BING_API_KEY")

	// Open the CSV file
	f, err := os.Open("data/companies.csv")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	r := csv.NewReader(f)
	for {
		// Read each record from csv
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}

		companyName := strings.TrimSpace(record[0])
		tickerCode := strings.TrimSpace(record[1])

		// Create folder if not exist
		dir := fmt.Sprintf("./data/learning-set/%s", tickerCode)
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			os.MkdirAll(dir, 0755)
		}

		// Search for images
		searchQuery := fmt.Sprintf("%s logo", companyName)
		searchQuery = sanitize.URI(searchQuery)

		client := &http.Client{}

		req, err := http.NewRequest("GET", fmt.Sprintf(bingSearchURL, searchQuery), nil)
		if err != nil {
			log.Fatalf("An Error Occured %v", err)
		}
		req.Header.Add("Ocp-Apim-Subscription-Key", apiKey)
	
		resp, err := client.Do(req)
		if err != nil {
			log.Fatalf("An Error Occured %v", err)
		}
		defer resp.Body.Close()

		// Parse response
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}
		bodyString := string(bodyBytes)
		fmt.Print(bodyString)

		// parse JSON to get a list of images
		images := gjson.Get(bodyString, "value").Array()
		if len(images) == 0 {
			fmt.Printf("No images found for %s\n", companyName)
			continue
		}

		// For each image, download it
		for _, image := range images {
			imageURL := image.Get("contentUrl").String()

			
			// Download and save the image
			err = downloadFile(imageURL, dir)
			if err != nil {
				fmt.Printf("Error downloading image: %s", err)
			}
		}
	}
}

func downloadFile(url string, targetDir string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath.Join(targetDir, filepath.Base(url)))
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}