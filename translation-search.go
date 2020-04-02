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
// h4982 search otHistory
func searchStrongsWord(searchString string) {
	// Convert 1 Kings to 1Kings, 2 Peter to 2Peter, etc.
	searchString = numberedBookRegex.ReplaceAllString(searchString, "$1$2")

	// Convert "old testament" to "old"
	searchString = testamentRegex.ReplaceAllString(searchString, "$1")
	// Split string on one or more whitespace chars
	words := strings.Fields(searchString)

	bookNames := findBooksMatchingWords(words[2:])

	if debug {
		fmt.Printf("Would search strongs %s in these books: %v\n", words[0], bookNames)
	}
}

//
// findBooksMatchingWords returns a list of bible book names (TranslationName) for the
// given set of zero or more filtering words (i.e. gospels, nt, prophesy, psalms, etc.)
//
func findBooksMatchingWords(words []string) *[]string {

	// This is a slice of book translation names that we will be returning
	// See Book.TranslationName
	var keepBooks []string

	// strongsWord := words[0]
	// filters := []string{}
	// Start at the 3rd word
	for _, word := range words {
		books, present := filters[word]
		if debug {
			fmt.Printf("=== filtering on word %s  found: %t\n", word, present)
		}
		if !present {
			if debug {
				fmt.Println(word, "not found")
			}
			continue
		}
		if len(keepBooks) == 0 {
			// This is the first found word
			// Copy books into keepBooks
			keepBooks = make([]string, len(books))
			// We will be changing books so we don't want a reference to books2
			copy(keepBooks, books)
		} else {

			// books := []int{1, 2, 3, 5}
			// books2 := []int{2, 3, 4, 8, 10, 19}
			// Remove all books not in books2
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
		if debug {
			fmt.Printf("After %s, keepBooks is now length %d\n", word, len(keepBooks))
			fmt.Printf("books: %v\n", keepBooks)
		}

		if len(keepBooks) == 0 {
			// If it's zero at this point, there is no point in going further
			break
		}
	}
	return &keepBooks
}
