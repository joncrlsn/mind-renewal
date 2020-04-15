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
	"strings"

	"github.com/gookit/color"
)

//
// https://biblewithann.wordpress.com/2013/10/21/divisions-classifications-of-bible-books/
//
const (
	law          BookCategory = iota
	history      BookCategory = iota
	poetry       BookCategory = iota
	prophesy     BookCategory = iota
	gospel       BookCategory = iota
	epistle      BookCategory = iota
	oldTestament Testament    = iota
	newTestament Testament    = iota
)

var (
	bookNameMap = map[string]Book{}

	apiToken = "9f29b9475bd3e4b05765c06741cf4c094eef8a8a"

	categoryAliases = map[string][]string{
		"history":      []string{"history", "hist", "historical"},
		"poetry":       []string{"poetry", "poet"},
		"prophesy":     []string{"prophesy", "prophetic"},
		"gospel":       []string{"gospel", "gospels"},
		"epistle":      []string{"epistle", "epistles"},
		"oldTestament": []string{"oldTestament", "ot", "old"},
		"newTestament": []string{"newTestament", "nt", "new"},
	}

	// filters is a map keyed by a descriptor with multiple book.TranslationName values
	// key can be a category or a testament name (old or new), or an alias name for a book.
	filters map[string][]string = make(map[string][]string)
)

// Passage is the result of a Bible passage lookup
type Passage struct {
	VerseRef string   `json:"canonical"`
	Passages []string `json:"passages"`
}

// Book represents a book of the Bible
type Book struct {
	FullName        string
	TranslationName string
	Testament       Testament
	Category        BookCategory
	Aliases         []string
}

// BookCategory is the classification type of the book
type BookCategory int

func (category BookCategory) String() string {
	classNames := [...]string{
		"law",
		"history",
		"poetry",
		"prophesy",
		"gospel",
		"epistle",
	}
	if category < law || category > epistle {
		return "Unknown"
	}
	return classNames[category]
}

// Testament is described as either old or new
type Testament int

func (t Testament) String() string {
	if t == oldTestament {
		return "oldTestament"
	} else if t == newTestament {
		return "newTestament"
	} else {
		return "unknown"
	}
}

func init() {

	for _, book := range books {
		fullNameLower := strings.ToLower(book.FullName)
		bookNameMap[fullNameLower] = book
		for _, aliasName := range book.Aliases {
			bookNameMap[aliasName] = book
			filters[aliasName] = append(filters[aliasName], book.TranslationName)
		}

		// Add the TranslationName of the book to the slices for each alias of the Testament
		for _, alias := range categoryAliases[book.Testament.String()] {
			lowerAlias := strings.ToLower(alias)
			filters[lowerAlias] = append(filters[lowerAlias], book.TranslationName)
		}

		// Add the TranslationName of the book to the slices for each alias of the BookClass
		//filters[book.Category.String()] = append(filters[book.Category.String()], book.TranslationName)
		for _, alias := range categoryAliases[book.Category.String()] {
			filters[alias] = append(filters[alias], book.TranslationName)
		}

		// Make sure the full name is represented as a key in the filters map
		key := strings.ToLower(book.FullName)
		filters[key] = append(filters[key], book.TranslationName)
	}
}

// books is used for the translate command.
var books = []Book{
	Book{"Genesis", "Gen", oldTestament, law, []string{"gen"}},
	Book{"Exodus", "Exo", oldTestament, law, []string{"ex", "exo"}},
	Book{"Leviticus", "Lev", oldTestament, law, []string{"lev"}},
	Book{"Numbers", "Num", oldTestament, law, []string{"nu", "num", "numb"}},
	Book{"Deuteronomy", "Deu", oldTestament, law, []string{"deu", "deut"}},
	Book{"Joshua", "Jos", oldTestament, history, []string{"jos", "josh"}},
	Book{"Judges", "Jdg", oldTestament, history, []string{"jdg", "judg"}},
	Book{"Ruth", "Rut", oldTestament, history, []string{"ru", "rut"}},
	Book{"1 Samuel", "1Sa", oldTestament, history, []string{"1sa", "1sam"}},
	Book{"2 Samuel", "2Sa", oldTestament, history, []string{"2sa", "2sam"}},
	Book{"1 Kings", "1Ki", oldTestament, history, []string{"1ki", "1kin", "1king"}},
	Book{"2 Kings", "2Ki", oldTestament, history, []string{"2ki", "2kin", "2king"}},
	Book{"1 Chronicles", "1Ch", oldTestament, history, []string{"1chro", "1chron"}},
	Book{"2 Chronicles", "2Ch", oldTestament, history, []string{"2chro", "2chron"}},
	Book{"Ezra", "Ezr", oldTestament, history, []string{"ez", "ezr"}},
	Book{"Nehemiah", "Neh", oldTestament, history, []string{"ne", "neh"}},
	Book{"Esther", "Est", oldTestament, history, []string{"es", "est", "esth"}},
	Book{"Job", "Job", oldTestament, poetry, []string{}},
	//	Book{"Psalm", "Psa", oldTestament, poetry, []string{"Ps", "Psalms"}},
	Book{"Psalms", "Psa", oldTestament, poetry, []string{"ps", "psa", "psalm"}},
	Book{"Proverbs", "Pro", oldTestament, poetry, []string{"pr", "pro", "prov"}},
	Book{"Ecclesiastes", "Ecc", oldTestament, poetry, []string{"ecc", "ec", "eccles"}},
	Book{"Song of Solomon", "Song", oldTestament, poetry, []string{"song", "song of songs"}},
	//	Book{"Song of Songs", "Song", oldTestament, poetry, []string{}},
	Book{"Isaiah", "Isa", oldTestament, prophesy, []string{"is", "isa"}},
	Book{"Jeremiah", "Jer", oldTestament, prophesy, []string{"je", "jer", "jere"}},
	Book{"Lamentations", "Lam", oldTestament, prophesy, []string{"la", "lam", "lamen"}},
	Book{"Ezekiel", "Ezek", oldTestament, prophesy, []string{"ezek"}},
	Book{"Daniel", "Dan", oldTestament, prophesy, []string{"dan"}},
	Book{"Hosea", "Hos", oldTestament, prophesy, []string{"hos"}},
	Book{"Joel", "Joel", oldTestament, prophesy, []string{"joe"}},
	Book{"Amos", "Amo", oldTestament, prophesy, []string{"am", "amo"}},
	Book{"Obadiah", "Oba", oldTestament, prophesy, []string{"ob", "oba"}},
	Book{"Jonah", "Jon", oldTestament, prophesy, []string{"jon"}},
	Book{"Micah", "Mic", oldTestament, prophesy, []string{"mic"}},
	Book{"Nahum", "Nah", oldTestament, prophesy, []string{"na", "nah"}},
	Book{"Habakkuk", "Hab", oldTestament, prophesy, []string{"hab"}},
	Book{"Zephaniah", "Zep", oldTestament, prophesy, []string{"zep", "zeph", "zef"}},
	Book{"Haggai", "Hag", oldTestament, prophesy, []string{"hag", "hagg"}},
	Book{"Zechariah", "Zec", oldTestament, prophesy, []string{"zec", "zech", "zek"}},
	Book{"Malachi", "Mal", oldTestament, prophesy, []string{"mal"}},
	Book{"Matthew", "Mat", newTestament, gospel, []string{"mat", "matt"}},
	Book{"Mark", "Mrk", newTestament, gospel, []string{"mrk", "mar"}},
	Book{"Luke", "Luk", newTestament, gospel, []string{"lu", "luk"}},
	Book{"John", "Jhn", newTestament, gospel, []string{"joh", "jhn"}},
	Book{"Acts", "Act", newTestament, history, []string{"ac", "act"}},
	Book{"Romans", "Rom", newTestament, epistle, []string{"ro", "rom"}},
	Book{"1 Corinthians", "1Co", newTestament, epistle, []string{"1co", "1cor"}},
	Book{"2 Corinthians", "2Co", newTestament, epistle, []string{"2co", "2cor"}},
	Book{"Galatians", "Gal", newTestament, epistle, []string{"gal"}},
	Book{"Ephesians", "Eph", newTestament, epistle, []string{"eph"}},
	Book{"Philippians", "Php", newTestament, epistle, []string{"php"}},
	Book{"Colossians", "Col", newTestament, epistle, []string{"col", "colo"}},
	Book{"1 Thessalonians", "1Th", newTestament, epistle, []string{"1th", "1the", "1thess"}},
	Book{"2 Thessalonians", "2Th", newTestament, epistle, []string{"2th", "2the", "2thess"}},
	Book{"1 Timothy", "1Ti", newTestament, epistle, []string{"1ti", "1tim"}},
	Book{"2 Timothy", "2Ti", newTestament, epistle, []string{"2ti", "2tim"}},
	Book{"Titus", "Tit", newTestament, epistle, []string{"tit"}},
	Book{"Philemon", "Phm", newTestament, epistle, []string{"phm", "philem"}},
	Book{"Hebrews", "Heb", newTestament, epistle, []string{"heb"}},
	Book{"James", "Jas", newTestament, epistle, []string{"jam", "jas", "jame"}},
	Book{"1 Peter", "1Pe", newTestament, epistle, []string{"1pe", "1pet"}},
	Book{"2 Peter", "2Pe", newTestament, epistle, []string{"2pe", "2pet"}},
	Book{"1 John", "1Jn", newTestament, epistle, []string{"1jn", "1jo", "1joh"}},
	Book{"2 John", "2Jn", newTestament, epistle, []string{"2jn", "2jo", "2joh"}},
	Book{"3 John", "3Jn", newTestament, epistle, []string{"3jn", "3jo", "3joh"}},
	Book{"Jude", "Jud", newTestament, epistle, []string{"jud"}},
	Book{"Revelation", "Rev", newTestament, prophesy, []string{"rev", "revel"}},
}

// DELETEME
func appendStringSliceMap(aMap map[string][]string, key string, value string) {
	// _, present := filters[key]
	// if !present {
	// 	newSlice := make([]string, 10)
	// 	filters[key] = newSlice
	// }
	aMap[key] = append(aMap[key], value)

}

// displayErrorText writes an error message
func displayErrorText(message string) {
	color.Red.Println(message)
}

// displayError writes an error message
func displayError(message string, err error) {
	color.Red.Printf("%s: %v\n", message, err)
}

// debug only prints if the debugFlag is true
func debug(format string, a ...interface{}) {
	if debugFlag {
		fmt.Printf("DEBUG: "+format, a...)
	}
}
