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
// What is a declaration?
// It is a bible verse that you have reworded to make it more personal, then
// added to your declarations file, one declaration per line.
//
// i.e. I am born of God, the evil one cannot touch me.  - 1 John 5:18

import (
	"fmt"
	"regexp"
	"strings"

	wordwrap "github.com/mitchellh/go-wordwrap"
)

const (
	declarationLineWidth    = 37
	declarationsPdfFilename = "declarations.pdf"
	declarationsFilename    = "/keybase/private/joncrlsn/declarations"
)

// referenceRegex expects the declaration to end with a period followed
// by one or more spaces and then a dash.
// i.e.  I stand in grace.  - Rom 5:2
var referenceRegex = regexp.MustCompile(`\.\s+-`)

// displayRandomDeclaration assumes a file with a declaration per line.
func displayRandomDeclaration() {

	line, err := grepRandom(declarationsFilename)
	if err != nil {
		displayError("Error reading declarations file", err)
	}

	// Convert this:
	// ... cannot not touch me.   - 1 John 5:18
	// into this on two lines:
	// ... cannot not touch me.
	//     - 1 John 5:18
	line = referenceRegex.ReplaceAllString(line, ".\n    -")

	// Now add a border and wrap the lines at the given length
	border := strings.Repeat("=", declarationLineWidth)
	wrapped := wordwrap.WrapString(line, declarationLineWidth)
	fmt.Println(border)
	fmt.Println(wrapped)
	fmt.Println(border)
}
