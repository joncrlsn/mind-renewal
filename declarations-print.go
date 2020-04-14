package main

//
// Generate a PDF of declarations from a file that has a declaration per line.
//

import (
	"bufio"
	"fmt"
	"os"

	"github.com/jung-kurt/gofpdf"
)

const (
	pdfLineHeight = 4   /* in mm */
	pdfLineWidth  = 195 /* in mm */
)

// GeneratePdf generates our pdf by adding text to the page
// then saving it to a file.
func GeneratePdf(inputFilename, outputFilename string) error {

	pdf := gofpdf.New("P", /* P=Portrait */
		"mm",     /* mm = millimeters */
		"Letter", /* Letter vs. A4 paper size*/
		"")       /* font directory */
	pdf.SetMargins(12, 15, -1)
	pdf.AddPage()
	pdf.SetFont("Arial", "", 10)

	// Create a channel that will supply each line in the file
	c, err := ReadLinesChannel(inputFilename)
	if err != nil {
		return err
	}

	// Loop over each line in the declarations file
	i := 0
	for declaration := range c {
		if i > 0 {
			// Provide an extra 2mm vertical space between declarations
			pdf.Ln(2)
		}

		// Write the declaration on however many lines are needed
		pdf.MultiCell(pdfLineWidth, /* width */
			pdfLineHeight, /* height */
			declaration,
			"0",   /* border 0=no-border */
			"LM",  /* align LM=middle-left */
			false) /* fill */

		i++
	}
	fmt.Printf("Saved %d page(s) to ./%s\n", pdf.PageCount(), outputFilename)

	return pdf.OutputFileAndClose(outputFilename)
}

// ReadLinesChannel reads a text file line by line into a channel.
//
//   c, err := fileutil.ReadLinesChannel(fileName)
//   if err != nil {
//      log.Fatalf("readLines: %s\n", err)
//   }
//   for line := range c {
//      fmt.Printf("  Line: %s\n", line)
//   }
//
// nil is returned (with the error) if there is an error opening the file
//
func ReadLinesChannel(filePath string) (<-chan string, error) {
	c := make(chan string)
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	go func() {
		defer file.Close()
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			c <- scanner.Text()
		}
		close(c)
	}()
	return c, nil
}
