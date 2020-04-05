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
// Bible Study CLI tool
//

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/gookit/color"
)

const (
	translationMapFileName = "TTESV.txt"
	strongsGreekFileName   = "strongsgreek.dat"
	strongsHebrewFileName  = "strongshebrew.dat"

	translationMapURL = "https://github.com/tyndale/STEPBible-Data/raw/master/TTESV%20-%20Tyndale%20Translation%20tags%20for%20ESV%20-%20TyndaleHouse.com%20STEPBible.org%20CC%20BY-NC.txt"
	strongsHebrewURL  = "https://raw.githubusercontent.com/openscriptures/strongs/master/hebrew/strongshebrew.dat"
	strongsGreekURL   = "https://raw.githubusercontent.com/openscriptures/strongs/master/greek/strongsgreek.dat"
)

var (
	reader             = bufio.NewReader(os.Stdin)
	nonNumericRegexp   = regexp.MustCompile(`[^0-9]`)
	previousPassageRef = ""
	dataDirName        = ".biblestudy-data"
	dataDirPath        string

	translationMapFile string
	strongsGreekFile   string
	strongsHebrewFile  string

	debugFlag bool
)

func init() {
	rand.Seed(time.Now().UnixNano())

	//
	// Possibly create the data directory
	//
	home, err := os.UserHomeDir()
	if err != nil {
		color.Red.Printf("Error finding home directory: %v\n", err)
		os.Exit(1)
	}

	dataDirPath = filepath.Join(home, dataDirName)
	translationMapFile = filepath.Join(dataDirPath, translationMapFileName)
	strongsGreekFile = filepath.Join(dataDirPath, strongsGreekFileName)
	strongsHebrewFile = filepath.Join(dataDirPath, strongsHebrewFileName)

	if _, err := os.Stat(dataDirPath); os.IsNotExist(err) {
		os.Mkdir(dataDirPath, 0774)
		fmt.Printf("Downloading data files to: %s\n", dataDirPath)
	}

	if _, err := os.Stat(translationMapFile); os.IsNotExist(err) {
		// Download TTESV, the translation file
		if err := DownloadFile(translationMapFile, translationMapURL); err != nil {
			color.Red.Printf("Error downloading url:\n  %s\n  %v\n", translationMapURL, err)
			os.Exit(1)
		}
	}

	if _, err := os.Stat(strongsGreekFile); os.IsNotExist(err) {
		// Download the Strongs Greek definitions
		if err := DownloadFile(strongsGreekFile, strongsGreekURL); err != nil {
			color.Red.Printf("Error downloading url:\n  %s\n  %v\n", strongsGreekURL, err)
			os.Exit(1)
		}
	}

	if _, err := os.Stat(strongsHebrewFile); os.IsNotExist(err) {
		// Download the Strongs Hebrew definitions
		if err := DownloadFile(strongsHebrewFile, strongsHebrewURL); err != nil {
			color.Red.Printf("Error downloading url:\n  %s\n  %v\n", strongsHebrewURL, err)
			os.Exit(1)
		}
	}
}

// Keep looping until the user decides to quit
func main() {
	// Loop on the main prompt
	for {
		mainPrompt()
	}
}

func mainPrompt() {

	// home, err := os.UserHomeDir()
	if len(previousPassageRef) > 0 {
		color.FgDarkGray.Printf("Current verse: %s", previousPassageRef)
		color.Cyan.Println("  (t)ranslate or (s)how it again")
	}
	color.Cyan.Println("Enter verse reference, strongs# (i.e. g4982 or h3068), (p)roverb, (h)elp or (q)uit.")
	color.Magenta.Print(" > ")
	text, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}

	// Remove leading or trailing spaces and newline from the end of text
	text = strings.TrimSpace(text)
	text = strings.ToLower(text)

	// Try again if no data was entered
	if len(text) == 0 {
		return
	}

	// Shall we turn debug on/off?
	isDebug, _ := regexp.MatchString(`^debug\s+(on|off)\s*$`, text)
	if isDebug {
		debugFlag = strings.Contains(text, "on")
		fmt.Printf("Set debug to %t\n", debugFlag)
		return
	}

	// Shall we show help?
	help, _ := regexp.MatchString(`^(help|h)$`, text)
	if help {
		printHelpMainPrompt()
		return
	}

	// Shall we exit?
	exit, _ := regexp.MatchString(`^(exit|x|q|quit)$`, text)
	if exit {
		os.Exit(0)
	}

	// Search on the ESV text?
	if strings.HasPrefix(text, "search ") {
		displaySearchResults(text[7:])
		return
	}

	// Do we have a Greek Strongs number?
	// Example: 'g4982' or 'G4982'
	strongsGreek, _ := regexp.MatchString(`^g\d+$`, text)
	if strongsGreek {
		displayStrongs(text, strongsGreekFile)
		return
	}

	// Do we have a Hebrew Strongs number?
	// Example: 'h7654' or 'H7654'
	strongsHebrew, _ := regexp.MatchString(`^h\d+$`, text)
	if strongsHebrew {
		displayStrongs(text, strongsHebrewFile)
		return
	}

	// WORK IN PROGRESS
	// Display verses that use the given strongs number
	// Example:  'g4982 search epistle'
	// Example:  'h4982 search prophecy' (search the old testament)
	strongsSearch, _ := regexp.MatchString(`^[gh]\d+ .*search`, text)
	if strongsSearch {
		searchStrongsWord(text)
		return
	}

	// Translate the latest verse to strongs numbers
	translatePrev, _ := regexp.MatchString(`^(t|tr|tran|trans|translate)$`, text)
	if translatePrev {
		if len(previousPassageRef) == 0 {
			displayErrorText("You have not looked up a verse to translate.")
		} else {
			translate(previousPassageRef)

			fmt.Println()
			fmt.Println("Find other verses that include a strongs number.  Example: g4982 search")
			fmt.Println()
		}
		return
	}

	// Show the latest verse again
	showPrev, _ := regexp.MatchString(`^(s|show)$`, text)
	if showPrev {
		if len(previousPassageRef) == 0 {
			displayErrorText("You have not looked up a verse to show.")
		} else {
			showVerse(previousPassageRef)
		}
		return
	}

	// Show a random proverb
	proverb, _ := regexp.MatchString(`^(p|prov|proverb|proverbs)$`, text)
	if proverb {
		previousPassageRef = randomProverb()
		return
	}

	// Show a random declaration.
	// A declaration is a biblical truth or verse that you have reworded to help you
	// internalize it as applying directly to you.
	// i.e.
	declaration, _ := regexp.MatchString(`^(d|declaration)$`, text)
	if declaration {
		displayRandomDeclaration()
		return
	}

	// Assume this is a verse reference
	showVerse(text)
}

// showVerse looks up the reference and displays it on system out
func showVerse(verseRef string) {
	// Show the verse
	book, _ := parseVerseRef(verseRef)
	if book != "" {
		previousPassageRef = displayPassage(verseRef,
			true, /*includeHeadings*/
			true, /*includeFootnotes*/
			true, /*indentPoetry*/
			true /*includeVerseNumbers*/)
		return
	}
}

func printHelpMainPrompt() {
	fmt.Println("Need Help?  You can enter:")
	fmt.Println("  - a verse (like Ps3.3 or James 4:11)")
	fmt.Println("  - (t)ranslate the latest verse requested")
	fmt.Println("  - (s)how text for the latest verse again")
	fmt.Println("  - search rabble - search for the given word 'rabble'")
	fmt.Println("                  - may not return all matches if too many verses")
	fmt.Println("                  - use quotes around phrases to limit search results")
	fmt.Println("  - a strongs number prefixed by 'g' (for greek)   e.g. g2222")
	fmt.Println("  - a strongs number prefixed by 'h' (for hebrew)  e.g. h5555")
	fmt.Println("  - <strongs-num> search epistles    - searches on strongs num")
	fmt.Println("  - (p)roverb prints a random proverb")
	fmt.Println("  - (d)eclaration - WIP - displays a random line from your declarations file")
	fmt.Println("  - (q)uit or e(x)it")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  > 2Tim 1.7             (shows text for 2 Tim 1:7)")
	fmt.Println("  > search rabble        (shows verses with the English word rabble)")
	fmt.Println("  > g4982                (shows definition of Strongs Greek 4982)")
	fmt.Println("  > g4982 search gospels (WIP: shows verses that use Strongs Greek 4982)") // WIP
	fmt.Println()
}
