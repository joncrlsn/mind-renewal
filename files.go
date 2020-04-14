package main

import (
	"bufio"
	"math/rand"
	"os"
	"regexp"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// grepRandom picks a random line from the file.  It does this by generating
// a random number for each line and, when that number exceeds the previous
// maximum, becomes the new line to be returned.
func grepRandom(fileName string) (string, error) {

	// This is the line to be returned and the random
	// value it was assigned.  When a larger value is
	// encountered the line is replaced with the new
	// one.
	var line string
	var maxInt int

	file, err := os.Open(fileName)
	if err != nil {
		return "", err
	}

	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		randomInt := rand.Intn(1000)
		if randomInt > maxInt {
			// We have a new line
			line = scanner.Text()
			maxInt = randomInt
		}
	}

	return line, nil
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
