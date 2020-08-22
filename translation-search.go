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

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	strongsWordSearchRegex = regexp.MustCompile(`^([gh]\d+) +search`)
	numberedBookRegex      = regexp.MustCompile(` ([123]) ([a-z])`)
	testamentRegex         = regexp.MustCompile(` (old|new) tes[a-z]*`)
)

// searchStrongsWord finds verses that use the specified strongs word
// The input string will look something like this:
// g4982 search gospels
// h4982 search history
func searchStrongsWord(searchString string) error {
	// Convert 1 Kings to 1Kings, 2 Peter to 2Peter, etc.
	searchString = numberedBookRegex.ReplaceAllString(searchString, "$1$2")

	// Convert "old testament" to "old"
	searchString = testamentRegex.ReplaceAllString(searchString, "$1")

	// Split string on one or more whitespace chars
	words := strings.Fields(searchString)
	strongsWord := words[0] //

	var allBooks = true
	var bookNames *[]string
	if len(words) > 2 {
		//
		// Find the book names that we'll use to limit the search results
		//
		allBooks = false
		bookNames = findBooksMatchingWords(words[2:])
		debug("Searching strongs %s in these books: %v\n", words[0], bookNames)
	}

	//
	// Traverse the translation file looking for verses that use the given strongsWord
	// The translation file has lines  that look like this:
	// $Gen 2:11 02=<08034> 05=<00259> 08=<06376> 09=<01931> ...
	//

	strongsNumber := strongsWord[1:] // Remove the prefix: "g" or "h"

	// Greek Strongs words have 4 digits, Hebrew has 5 digits (left-zero padded)
	var format = "[<+]%04s[+>]" // greek is the default (4 digits)
	if strings.HasPrefix(strongsWord, "h") {
		format = "[<+]%05s[+>]" // switch to 5 digits
	}

	lookupRegex := regexp.MustCompile(fmt.Sprintf(format, strongsNumber))
	// Grep the file
	c, err := grep(translationMapFile, lookupRegex)
	if err != nil {
		displayErrorText(fmt.Sprintf("Unable to read file: %s, %v\n", translationMapFile, err))
		return err
	}

	var bookPattern *regexp.Regexp
	if !allBooks {
		// Pattern is a regular expression that contains all the book names we are allowing in our search
		patternStr := `^\$(` + strings.Join(*bookNames, "|") + `) `
		bookPattern, err = regexp.Compile(patternStr)
		if err != nil {
			displayError("Error compiling regex "+patternStr, err)
			return err
		}
		debug("Regex of book names to limit search results: %v\n", patternStr)
	}

	var verses []string
	for line := range c {
		// If the pattern matches then we can accept this line
		if allBooks || bookPattern.MatchString(line) {
			// The verse reference is the first part of the line until the first tab
			tabIx := strings.Index(line, "\t")
			verses = append(verses, line[1:tabIx])
		}
	}

	debug("Found %d verses %v\n", len(verses), verses)

	if len(verses) > 20 {
		fmt.Printf("Showing only 20 of %d verses found.\n", len(verses))
		verses = verses[:19]
	}

	// TODO:  Allow user to page through results

	//
	// Lookup all the verses that use the given strongs number
	//

	versesLookupString := strings.Join(verses, " ")
	//versesLookupString = strings.ReplaceAll(versesLookupString, ":", ".")
	//versesLookupString = strings.ReplaceAll(versesLookupString, " ", "")
	fmt.Printf("Verses lookup string: %s\n", versesLookupString)
	passage, err := lookupVerse(versesLookupString, 0,
		false, /*includeHeadings*/
		false, /*includeFootnotes*/
		false, /*indentPoetry*/
		false /*includeVerseNumbers*/)
	if err != nil {
		displayError("Error looking up verse", err)
		return err
	}

	debug("Passage: %v\n", passage)

	//
	// Print out the passages we found
	//

	found := false
	newLineRegex := regexp.MustCompile(`[\n]`)
	for _, passageText := range passage.Passages {
		// Put verse on one line
		newText := newLineRegex.ReplaceAllString(passageText, " ")
		newText = strings.ReplaceAll(newText, "(ESV)", "")
		fmt.Println(newText)
		found = true
	}
	if found {
		fmt.Println("(ESV)")
	}
	fmt.Println()

	return nil
}

//
// findBooksMatchingWords returns a list of bible book names (TranslationName) for the
// given set of zero or more filtering words (i.e. gospels, nt, prophesy, psalms, etc.)
//
func findBooksMatchingWords(words []string) *[]string {

	// This is a slice of book translation names that we will be returning
	// See Book.TranslationName
	var keepBooks []string

	for _, word := range words {
		books, present := filters[word]
		debug("=== filtering on word %s  valid: %t\n", word, present)
		if !present {
			continue
		}
		if len(keepBooks) == 0 {
			// This is the first found word
			keepBooks = make([]string, len(books))
			// Copy books into keepBooks
			copy(keepBooks, books)

		} else {

			// Find the intersection of keepBooks and books
			tempKeep := []string{}
			for _, b := range keepBooks {
				for _, b2 := range books {
					if b == b2 {
						// Keep the "i"th element in books
						tempKeep = append(tempKeep, b)
						break
					}
				}
			}
			keepBooks = tempKeep
		}

		debug("After %s, keepBooks is now length %d\n", word, len(keepBooks))
		debug("books: %v\n", keepBooks)

		if len(keepBooks) == 0 {
			// If it's zero at this point, there is no point in going further
			break
		}
	}
	return &keepBooks
}
