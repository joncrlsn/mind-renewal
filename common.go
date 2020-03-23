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
	"github.com/gookit/color"
)

var (
	annotationsBookMap = map[string]Book{}
)

type Passage struct {
	VerseRef string   `json:"canonical"`
	Passages []string `json:"passages"`
}

type Book struct {
	FullName        string
	TranslationName string
	NewTestament    bool
	Permutations    []string
}

func init() {
	for _, book := range books {
		annotationsBookMap[book.FullName] = book
	}
}

// books is used for the translate command.
var books = []Book{
	Book{"Genesis", "Gen", false, []string{}},
	Book{"Exodus", "Exo", false, []string{}},
	Book{"Leviticus", "Lev", false, []string{}},
	Book{"Numbers", "Num", false, []string{}},
	Book{"Deuteronomy", "Deu", false, []string{}},
	Book{"Joshua", "Jos", false, []string{}},
	Book{"Judges", "Jdg", false, []string{}},
	Book{"Ruth", "Rut", false, []string{}},
	Book{"1 Samuel", "1Sa", false, []string{}},
	Book{"2 Samuel", "2Sa", false, []string{}},
	Book{"1 Kings", "1Ki", false, []string{}},
	Book{"2 Kings", "2Ki", false, []string{}},
	Book{"1 Chronicles", "1Ch", false, []string{}},
	Book{"2 Chronicles", "2Ch", false, []string{}},
	Book{"Ezra", "Ezr", false, []string{}},
	Book{"Nehemiah", "Neh", false, []string{}},
	Book{"Esther", "Est", false, []string{}},
	Book{"Job", "Job", false, []string{}},
	Book{"Psalm", "Psa", false, []string{}},
	Book{"Psalms", "Psa", false, []string{}},
	Book{"Proverbs", "Pro", false, []string{}},
	Book{"Ecclesiastes", "Ecc", false, []string{}},
	Book{"Song of Solomon", "Song", false, []string{}},
	Book{"Song of Songs", "Song", false, []string{}},
	Book{"Isaiah", "Isa", false, []string{}},
	Book{"Jeremiah", "Jer", false, []string{}},
	Book{"Lamentations", "Lam", false, []string{}},
	Book{"Ezekiel", "Ezek", false, []string{}},
	Book{"Daniel", "Dan", false, []string{}},
	Book{"Hosea", "Hos", false, []string{}},
	Book{"Joel", "Joel", false, []string{}},
	Book{"Amos", "Amo", false, []string{}},
	Book{"Obadiah", "Oba", false, []string{}},
	Book{"Jonah", "Jon", false, []string{}},
	Book{"Micah", "Mic", false, []string{}},
	Book{"Nahum", "Nah", false, []string{}},
	Book{"Habakkuk", "Hab", false, []string{}},
	Book{"Zephaniah", "Zep", false, []string{}},
	Book{"Haggai", "Hag", false, []string{}},
	Book{"Zechariah", "Zec", false, []string{}},
	Book{"Malachi", "Mal", false, []string{}},
	Book{"Matthew", "Mat", true, []string{}},
	Book{"Mark", "Mrk", true, []string{}},
	Book{"Luke", "Luk", true, []string{}},
	Book{"John", "Jhn", true, []string{}},
	Book{"Acts", "Act", true, []string{}},
	Book{"Romans", "Rom", true, []string{}},
	Book{"1 Corinthians", "1Co", true, []string{}},
	Book{"2 Corinthians", "2Co", true, []string{}},
	Book{"Galatians", "Gal", true, []string{}},
	Book{"Ephesians", "Eph", true, []string{}},
	Book{"Philippians", "Php", true, []string{}},
	Book{"Colossians", "Col", true, []string{}},
	Book{"1 Thessalonians", "1Th", true, []string{}},
	Book{"2 Thessalonians", "2Th", true, []string{}},
	Book{"1 Timothy", "1Ti", true, []string{}},
	Book{"2 Timothy", "2Ti", true, []string{}},
	Book{"Titus", "Tit", true, []string{}},
	Book{"Philemon", "Phm", true, []string{}},
	Book{"Hebrews", "Heb", true, []string{}},
	Book{"James", "Jas", true, []string{}},
	Book{"1 Peter", "1Pe", true, []string{}},
	Book{"2 Peter", "2Pe", true, []string{}},
	Book{"1 John", "1Jn", true, []string{}},
	Book{"2 John", "2Jn", true, []string{}},
	Book{"3 John", "3Jn", true, []string{}},
	Book{"Jude", "Jud", true, []string{}},
	Book{"Revelation", "Rev", true, []string{}},
}

// annotationsBookMap is used for the translate command.
//var annotationsBookMap = map[string]Book{
// "Genesis":         "Gen",
// "Exodus":          "Exo",
// "Leviticus":       "Lev",
// "Numbers":         "Num",
// "Deuteronomy":     "Deu",
// "Joshua":          "Jos",
// "Judges":          "Jdg",
// "Ruth":            "Rut",
// "1 Samuel":        "1Sa",
// "2 Samuel":        "2Sa",
// "1 Kings":         "1Ki",
// "2 Kings":         "2Ki",
// "1 Chronicles":    "1Ch",
// "2 Chronicles":    "2Ch",
// "Ezra":            "Ezr",
// "Nehemiah":        "Neh",
// "Esther":          "Est",
// "Job":             "Job",
// "Psalm":           "Psa",
// "Psalms":          "Psa",
// "Proverbs":        "Pro",
// "Ecclesiastes":    "Ecc",
// "Song of Solomon": "Song",
// "Song of Songs":   "Song",
// "Isaiah":          "Isa",
// "Jeremiah":        "Jer",
// "Lamentations":    "Lam",
// "Ezekiel":         "Ezek",
// "Daniel":          "Dan",
// "Hosea":           "Hos",
// "Joel":            "Joel",
// "Amos":            "Amo",
// "Obadiah":         "Oba",
// "Jonah":           "Jon",
// "Micah":           "Mic",
// "Nahum":           "Nah",
// "Habakkuk":        "Hab",
// "Zephaniah":       "Zep",
// "Haggai":          "Hag",
// "Zechariah":       "Zec",
// "Malachi":         "Mal",
// "Matthew":         "Mat",
// "Mark":            "Mrk",
// "Luke":            "Luk",
// "John":            "Jhn",
// "Acts":            "Act",
// "Romans":          "Rom",
// "1 Corinthians":   "1Co",
// "2 Corinthians":   "2Co",
// "Galatians":       "Gal",
// "Ephesians":       "Eph",
// "Philippians":     "Php",
// "Colossians":      "Col",
// "1 Thessalonians": "1Th",
// "2 Thessalonians": "2Th",
// "1 Timothy":       "1Ti",
// "2 Timothy":       "2Ti",
// "Titus":           "Tit",
// "Philemon":        "Phm",
// "Hebrews":         "Heb",
// "James":           "Jas",
// "1 Peter":         "1Pe",
// "2 Peter":         "2Pe",
// "1 John":          "1Jn",
// "2 John":          "2Jn",
// "3 John":          "3Jn",
// "Jude":            "Jud",
// "Revelation":      "Rev",
//}

// displayErrorText writes an error message
func displayErrorText(message string) {
	color.Red.Println(message)
}

// displayError writes an error message
func displayError(err error) {
	color.Red.Println(err)
}
