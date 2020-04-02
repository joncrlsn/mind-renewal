/*
Copyright © 2020 Jon Carlson <joncrlsn@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package main

//
// Handles ESV text search
//

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/pkg/errors"
)

const (
	baseSearchUrl = "https://api.esv.org/v3/passage/search"
)

// SearchResults holds the results from one search
type SearchResults struct {
	Results []Result `json:"results"`
}

// Result represents a verse that matches the search text
type Result struct {
	Reference string `json:"reference"`
	Content   string `json:"content"`
}

// displaySearchResults shows results of searching for a given word or words
func displaySearchResults(searchString string) error {
	results, err := searchESV(searchString)
	if err != nil {
		return err
	}

	if len(results.Results) == 0 {
		displayErrorText("No results found")
	} else {
		for _, result := range results.Results {
			fmt.Printf("%s - %s\n\n", result.Reference, result.Content)
		}
	}

	return nil
}

// searchESV sends the searchString to the API and displays the results.
// This is not ideal given that you cannot choose which testament to search in.
// Search results are currently capped at 100 so if you search on "fear" you
// will never receive any NT results.
func searchESV(searchString string) (*SearchResults, error) {
	urlSafeSearchString := strings.ReplaceAll(searchString, " ", "+")

	url := fmt.Sprintf(`%s?q=%s&page-size=100&page=1`, baseSearchUrl, urlSafeSearchString)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, errors.Wrap(err, "Error reading request.")
	}

	req.Header.Set("Authorization", "Token "+apiToken)

	client := &http.Client{Timeout: time.Second * 10}
	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "Error reading response.")
	}
	defer resp.Body.Close()

	jsonBody := SearchResults{}
	err = json.NewDecoder(resp.Body).Decode(&jsonBody)
	if err != nil {
		return nil, errors.Wrap(err, "Error reading response.")
	}

	return &jsonBody, nil
}
