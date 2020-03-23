package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	//"https://github.com/mitchellh/go-homedir"
	"github.com/gookit/color"
)

const (
	translationMapFileName = "TTESV.txt"
	strongsGreekFileName   = "strongsgreek.dat"
	strongsHebrewFileName  = "strongshebrew.dat"

	translationMapUrl = "https://github.com/tyndale/STEPBible-Data/raw/master/TTESV%20-%20Tyndale%20Translation%20tags%20for%20ESV%20-%20TyndaleHouse.com%20STEPBible.org%20CC%20BY-NC.txt"
	strongsHebrewUrl  = "https://raw.githubusercontent.com/openscriptures/strongs/master/hebrew/strongshebrew.dat"
	strongsGreekUrl   = "https://raw.githubusercontent.com/openscriptures/strongs/master/greek/strongsgreek.dat"
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

		// Download TTESV, the translation file
		if err := DownloadFile(translationMapFile, translationMapUrl); err != nil {
			color.Red.Printf("Error downloading url:\n  %s\n  %v\n", translationMapUrl, err)
			os.Exit(1)
		}

		// Download the Strongs Greek definitions
		if err := DownloadFile(strongsGreekFile, strongsGreekUrl); err != nil {
			color.Red.Printf("Error downloading url:\n  %s\n  %v\n", strongsGreekUrl, err)
			os.Exit(1)
		}

		// Download the Strongs Hebrew definitions
		if err := DownloadFile(strongsHebrewFile, strongsHebrewUrl); err != nil {
			color.Red.Printf("Error downloading url:\n  %s\n  %v\n", strongsHebrewUrl, err)
			os.Exit(1)
		}

	} else {
		fmt.Printf("Using data directory: %s\n", dataDirPath)
	}
}

func main() {

	// Loop on the main prompt
	for {
		mainPrompt()
	}
}

func mainPrompt() {

	// home, err := os.UserHomeDir()
	if len(previousPassageRef) > 0 {
		color.FgDarkGray.Printf("Current verse: %s  (t)ranslate or (s)how\n", previousPassageRef)
	}
	color.Cyan.Println("Enter verse reference, strongs# (i.e. G1234 or H5678), (p)roverb, (h)elp or (q)uit.")
	color.Cyan.Print("Command: ")
	text, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}

	// Remove newline from the end of text
	text = strings.Replace(text, "\n", "", 1)
	text = strings.ToLower(text)

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

	// Do we have a Greek Strongs number?
	strongsGreek, _ := regexp.MatchString(`^g\d+$`, text)
	if strongsGreek {
		displayStrongs(text, strongsGreekFile)
		return
	}

	// Do we have a Hebrew Strongs number?
	strongsHebrew, _ := regexp.MatchString(`^h\d+$`, text)
	if strongsHebrew {
		displayStrongs(text, strongsHebrewFile)
		return
	}

	if len(previousPassageRef) > 0 {
		// Translate the latest verse to strongs numbers
		translatePrev, _ := regexp.MatchString(`^(t|tr|tra|tran|trans|translate)$`, text)
		if translatePrev {
			translate(previousPassageRef)
			return
		}

		// Show the latest verse again
		showPrev, _ := regexp.MatchString(`^(s|show)$`, text)
		if showPrev {
			showVerse(previousPassageRef)
			return
		}
	}

	// Show a random proverb
	proverb, _ := regexp.MatchString(`^(p|pr|pro|prov|proverb|proverbs)$`, text)
	if proverb {
		previousPassageRef = randomProverb()
		return
	}

	// Assume this must be a verse reference
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
	fmt.Println("Need Help? You can enter:")
	fmt.Println("  - a verse (like Ps1.1 or James 4:5)")
	fmt.Println("  - (t)ranslate the latest verse displayed")
	fmt.Println("  - (s)how the latest verse again")
	fmt.Println("  - a strongs number prefixed by 'g' (for greek)   e.g. g2222")
	fmt.Println("  - a strongs number prefixed by 'h' (for hebrew)  e.g. h5555")
	fmt.Println("  - (p)roverb prints a random proverb")
	fmt.Println("  - (q)uit or e(x)it")
	fmt.Println()
}
