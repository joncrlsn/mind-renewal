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
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

var (
	whitespaceRegex  = regexp.MustCompile(`\s+`)
	newlineRegex     = regexp.MustCompile(`\n+`)
	plusWordNumRegex = regexp.MustCompile(`([0-9]+\+)+`)
	strongsRegex     = regexp.MustCompile(`<.+>$`)
	numberRegex      = regexp.MustCompile(`[0-9]+`)
)

func translate(verseRef string) {
	passage, err := lookupVerse(verseRef, 0,
		false, /*includeHeadings*/
		false, /*includeFootnotes*/
		false, /*indentPoetry*/
		true /*includeVerseNumbers*/)
	if err != nil {
		displayError(err)
		return
	}

	if len(passage.Passages) == 0 {
		displayErrorText("Passage not found")
	} else {

		// Parse the book name and chapter-verse sections from the verse reference
		singleVerseFormat := regexp.MustCompile(`^([0-9\s]*[^0-9]+)([0-9]+:?[0-9]*)`)
		book := singleVerseFormat.ReplaceAllString(verseRef, "$1")
		chapterVerse := singleVerseFormat.ReplaceAllString(verseRef, "$2")
		bookTrimmed := strings.TrimSpace(book)

		//fmt.Printf("book: '%s'\n", book)
		//fmt.Printf("chapterVerse: '%s'\n", chapterVerse)
		bookObj, ok := annotationsBookMap[bookTrimmed]
		if !ok {
			displayErrorText("Unable to find book with name " + bookTrimmed)
			return
		}
		translationMapLookupString := bookObj.TranslationName + " " + chapterVerse
		//translationMapLookupString := buildTranslationLookupString(passage.VerseRef)
		isNewTestament := bookObj.NewTestament

		//fmt.Println("annotationLookupString:" + annotationLookupString)

		p := passage.Passages[0]
		//fmt.Printf("Verse text: '%s'\n", p)

		// This regex splits the passage into two lines
		// Line 1 (index 0) is the reference (i.e. Mark 11:24)
		// Line 2 (index 1) is the text of the verse with verse numbers in
		//        square brackets.
		verseLines := newlineRegex.Split(p, -1)
		//fmt.Printf("NumLines: %d\n", len(lines))

		//
		// Find the ESV Strongs translation mapping line for this verse
		//
		lookupRegex, err := regexp.Compile("^[$]" + translationMapLookupString + "\t")
		if err != nil {
			displayError(err)
		}

		// Grep the file
		c, err := grep(translationMapFile, lookupRegex)
		if err != nil {
			displayErrorText(fmt.Sprintf("Unable to read file: %s, %v\n", translationMapFile, err))
			return
		}
		var translationMapLine string
		for line := range c {
			// Remove the lookup reference from the translation mapping line
			translationMapLine = lookupRegex.ReplaceAllString(line, "")
			//fmt.Printf("before:'%s'", line)
			//fmt.Printf("after:'%s'", wordMapping)
		}
		if len(translationMapLine) == 0 {
			displayErrorText("Unable to locate translation map line for " + passage.VerseRef)
			return
		}

		fmt.Println(passage.VerseRef)

		// Verses in Psalms, etc are on more than one line so put them all on the same line
		var text string
		if len(verseLines) > 1 {
			for i := 1; i < len(verseLines); i++ {
				// String concatentation is not very efficient, but
				// this is just for one verse.
				if i > 1 {
					text += " "
				}
				text += verseLines[i]
			}
		}

		// For processing, remove the (ESV) copyright at the end of the line.
		// We'll print it at the end.
		text = strings.Replace(text, "(ESV)", "", 1)

		// Debug:
		//fmt.Printf("Verse text: '%s'\n", p)
		//fmt.Printf("Verse text on one line: '%s'\n", text)
		//fmt.Printf("Translation map line  : '%s'\n", translationMapLine)

		err = printEnglishWithStrongs(text, translationMapLine, isNewTestament)
		if err != nil {
			displayErrorText(fmt.Sprintf("There was an error while annotating your verse. %v\n", err))
		}
		fmt.Println("(ESV)")
	}
}

// buildAnnotationLookupString converts "2 Timothy 1:7" into "2Ti 1:7" which is
// the format needed to find the ESV Strong's mappings
// func buildTranslationLookupString(verseRef string) string {

// 	// Grab the book name and chapter-verse sections from the verse reference
// 	singleVerseFormat := regexp.MustCompile(`^([0-9\s]*[^0-9]+)([0-9]+:?[0-9]*)`)
// 	book := singleVerseFormat.ReplaceAllString(verseRef, "$1")
// 	chapterVerse := singleVerseFormat.ReplaceAllString(verseRef, "$2")
// 	bookTrimmed := strings.TrimSpace(book)

// 	//fmt.Printf("book: '%s'\n", book)
// 	//fmt.Printf("chapterVerse: '%s'\n", chapterVerse)
// 	newBook, ok := annotationsBookMap[bookTrimmed]
// 	if !ok {
// 		displayErrorText("Unable to find book with name " + bookTrimmed)
// 		return ""
// 	}
// 	return newBook.TranslationName + " " + chapterVerse
// }

// printEnglishWithStrongs prints the English words with their Strongs number.
// This takes a bit of vertical space, but is easy to read.
func printEnglishWithStrongs(text string, strongsMap string, isNewTestament bool) error {

	// Ensure we have no leading or trailing spaces
	text = strings.TrimSpace(text)
	strongsMap = strings.TrimSpace(strongsMap)

	// Split both inputs using whitespace as a delimeter
	verseWords := whitespaceRegex.Split(text, -1)
	wordMappings := whitespaceRegex.Split(strongsMap, -1)

	//fmt.Println("verseWords:", verseWords)
	//fmt.Println("wordMappings:", wordMappings)

	maxLineLength := 0
	var lines []string
	var line string
	i := 0
	j := 0

	for {
		i++
		if i < len(verseWords) {
			line = line + " " + verseWords[i]
			//fmt.Printf("Adding word to line: %s\n", verseWords[i])
			//fmt.Printf("%s ", verseWords[i])
		} else {
			//fmt.Println()
			break
		}

		// Expecting wordmap to look like this:
		//  01=<3972>  or  12+13=<2596>
		wordMap := strings.Split(wordMappings[j], "=")

		if len(wordMap) == 2 {
			numberString := wordMap[0]
			strongs := wordMap[1]

			// Remove the first number in this case: "12+13"
			if strings.Contains(numberString, "+") {
				numberString = plusWordNumRegex.ReplaceAllString(numberString, "")
			}

			wordNum, err := strconv.Atoi(numberString)
			if err != nil {
				return errors.Wrap(err, "Error converting to intger: "+numberString)
			}

			// If we have a strongs number for this word in the text, print it
			// and start a new line.
			if wordNum == i {
				// We are done with the line, so add the strongs number
				// fmt.Printf(" %s\n", strongs)
				if len(line) > maxLineLength {
					maxLineLength = len(line)
				}
				line = line + strongs
				lines = append(lines, line)
				//fmt.Printf("line length:%d  %s\n", len(line), line)
				if j < len(wordMappings)-1 {
					j++
				}
				line = ""

			}
		}
	}
	//fmt.Printf("MaxLineLength: %d\n", maxLineLength)
	//fmt.Printf("Number of lines: %d\n", len(lines))

	strongsPrefix := "H"
	if isNewTestament {
		strongsPrefix = "G"
	}

	// Left pad each line to the max line length.  This generates a Printf format to use.
	format := "%" + strconv.Itoa(maxLineLength) + "s %s"
	// Write the padded lines to stdout
	for _, l := range lines {
		tmpStrongs := strongsRegex.FindString(l)
		english := l[:len(l)-len(tmpStrongs)]

		// Convert tmpStrongs <1111> to G1111  (G is for Greek)
		// Convert tmpStrongs <1111+2222> to G1111+G2222
		// Convert tmpStrongs <1111>+<2222> to G1111+G2222

		// Find the individual strong's numbers in the string
		strongsNumbers := numberRegex.FindAllString(tmpStrongs, -1)
		tmpStrongs = ""
		for i, strongs := range strongsNumbers {
			if i > 0 {
				tmpStrongs += " "
			}
			tmpStrongs += strongsPrefix + strongs
		}

		fmt.Printf(format+"\n", english, tmpStrongs)
	}

	//fmt.Printf("(ESV)\n")
	return nil
}

// grep reads a file line by line and adds to the channel only lines
// that match the given regexp.
func grep(fileName string, regex *regexp.Regexp) (<-chan string, error) {
	c := make(chan string, 10)
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}

	go func() {
		defer file.Close()
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			if regex.Match(scanner.Bytes()) {
				c <- scanner.Text()
			}
		}
		close(c)
	}()

	return c, nil
}
