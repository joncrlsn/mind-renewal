package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/pkg/errors"
)

var (
	verseRegexp            = regexp.MustCompile(`([0-9]*\s*[A-Za-z]+)\s*([0-9]+)[:.]?([0-9]*)`)
	baseApiUrl             = "https://api.esv.org/v3/passage/text/"
	proverbsChapterLengths = []int{
		33, 22, 35, 27, 23, 35, 27, 36, 18, 32,
		31, 28, 25, 35, 33, 33, 28, 24, 29, 30,
		31, 29, 35, 34, 28, 28, 27, 28, 27, 33,
		31,
	}
)

// randomProverb prints a random verse from Proverbs
func randomProverb() string {
	chapter := rand.Intn(len(proverbsChapterLengths))
	//fmt.Printf("numChapters:%d ix:%d\n", len(proverbsChapterLengths), ix)
	verse := rand.Intn(proverbsChapterLengths[chapter] + 1)
	reference := fmt.Sprintf("Proverbs %d:%d", chapter, verse)
	return displayPassage(reference,
		false, /*includeHeadings*/
		false, /*includeFootnotes*/
		false, /*indentPoetry*/
		false /*includeVerseNumbers*/)
}

// parseVerseRef parses the input string into book and chapterAndVerse
func parseVerseRef(text string) (book string, chapterAndVerse string) {
	if verseRegexp.MatchString(text) {
		book = verseRegexp.ReplaceAllString(text, "$1")
		chapterAndVerse = verseRegexp.ReplaceAllString(text, "$2")
		return
	}
	return "", ""
}

// Print out the passage from the reference given
func displayPassage(passageRef string, includeHeadings, includeFootnotes, indentPoetry, includeVerseNumbers bool) (cleanPassageRef string) {
	passage, err := lookupVerse(passageRef, 80,
		includeHeadings,
		includeFootnotes,
		indentPoetry,
		includeVerseNumbers)
	if err != nil {
		displayError(err)
		return
	}

	if len(passage.Passages) == 0 {
		displayErrorText("Passage not found")
	} else {
		cleanPassageRef = passage.VerseRef
		for _, passageText := range passage.Passages {
			//fmt.Println("============================================================")
			fmt.Println(passageText)
		}
	}

	return cleanPassageRef
}

// lookupVerse returns the result of an HTTP REST request to the ESV scriptures
// API. The ESV API supports returning multiple verses or even multiple passages
// in one lookup request.
//
// The API is pretty flexible and does its best to decipher what you are looking
// for.
//
// Example verseRef values that will work:
//    romans 12:1
//    2Tim1:13
//    Psalm 3:3,Isaiah 53:5
//    ps 119:9, 11
//    1 Thess 5:16-18
func lookupVerse(verseRef string, lineLength int, includeHeadings, includeFootnotes, indentPoetry, includeVerseNumbers bool) (*Passage, error) {
	urlSafeVerseRef := strings.ReplaceAll(verseRef, " ", "+")

	url := fmt.Sprintf(`%s?q=%s&line-length=%d&include-headings=%t&include-footnotes=%t&indent-poetry=%t&include-verse-numbers=%t`,
		baseApiUrl,
		urlSafeVerseRef,
		lineLength,
		includeHeadings,
		includeFootnotes,
		indentPoetry,
		includeVerseNumbers)

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

	jsonBody := Passage{}
	err = json.NewDecoder(resp.Body).Decode(&jsonBody)
	if err != nil {
		return nil, errors.Wrap(err, "Error reading response.")
	}

	return &jsonBody, nil
}
