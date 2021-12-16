package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

const FILENAME = "realhuman_phill.txt"

func find(file *os.File, target string) (int, int, error) {
	scanner := bufio.NewScanner(file)

	lineCount := 0
	occurences := 0

	for scanner.Scan() {
		password := scanner.Text()

		if strings.Contains(password, target) {
			fmt.Println(password)
			occurences++
		}

		lineCount++
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("scan file error: %v", err)
		return 0, 0, err
	}

	return lineCount, occurences, nil
}

func main() {

	file, err := os.Open(FILENAME)

	if err != nil {
		fmt.Printf("Could not read file %s\n", FILENAME)
		panic(err)
	}

	defer file.Close()

	start := time.Now()

	for _, targetPassword := range []string{"japierdole", "kurwa", "dupa", "pizda"} {
		lineCount, occurences, err := find(file, targetPassword)
		if err != nil {
			log.Fatal(err)
			return
		}

		fmt.Printf("%s: %d, lines searched, %d occurences\n", targetPassword, lineCount, occurences)
		fmt.Printf("=============\n")
	}

	end := time.Now()

	diff := end.Sub(start)

	fmt.Printf("\nDuration (in seconds): %.2f sec\n", diff.Seconds())
}
