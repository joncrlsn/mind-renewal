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

	"github.com/pkg/errors"
)

func displayStrongs(text string, file string) {
	// Remove all non-digits (which should be the first character or nothing)
	text = nonNumericRegexp.ReplaceAllString(text, "")

	// Grep the appropriate lines from the file
	c, err := chooseLines(file, text)
	if err != nil {
		displayError("Error reading lines from "+file, err)
		os.Exit(1)
	}

	// Print the lines
	found := false
	for line := range c {
		found = true
		fmt.Println(line)
	}
	if found {
		fmt.Println()
	} else {
		displayErrorText("Definition not found")
	}
}

// chooseLines returns only the lines that apply to the given strongs number.
//
// Example Lines:
//
// $$T0000003
// \00003\
//  3  Abaddon  ab-ad-dohn'
//
//  of Hebrew origin (11); a destroying angel:--Abaddon.
//  see HEBREW for 011
// $$T0000004
//
func chooseLines(fileName string, strongsNum string) (<-chan string, error) {
	c := make(chan string, 10)
	strongsInt, err := strconv.Atoi(strongsNum)
	if err != nil {
		return nil, errors.Wrap(err, "Invalid strongsNum argument to chooseLines")
	}
	paddedStrongs := fmt.Sprintf("%07d", strongsInt)

	startRegex := regexp.MustCompile(`^\$\$T` + paddedStrongs)
	endRegex := regexp.MustCompile(`^\$\$T`)
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}

	go func() {
		defer file.Close()
		scanner := bufio.NewScanner(file)
		var start bool
		for scanner.Scan() {
			if start {
				if endRegex.Match(scanner.Bytes()) {
					// Do not include the endRegex line
					break
				}
				c <- scanner.Text()
			}
			if startRegex.Match(scanner.Bytes()) {
				start = true
			}

		}
		close(c)
	}()

	return c, nil
}
