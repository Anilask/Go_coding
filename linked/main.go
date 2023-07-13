package main

import (
	"bufio"
	"fmt"
	"os"
	"sync"
)

func readFile(fileName string, wg *sync.WaitGroup, results chan<- []string, errors chan<- error) {
	// Open the file
	file, err := os.Open(fileName)
	if err != nil {
		errors <- fmt.Errorf("error opening file %s: %v", fileName, err)
		wg.Done()
		return
	}
	defer file.Close()

	// Create a slice to store the lines of the file
	var lines []string

	// Create a new scanner to read the file line by line
	scanner := bufio.NewScanner(file)

	// Read the file line by line
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}

	// Check for any errors encountered during scanning
	if err := scanner.Err(); err != nil {
		errors <- fmt.Errorf("error scanning file %s: %v", fileName, err)
		wg.Done()
		return
	}

	// Send the lines to the results channel
	results <- lines

	wg.Done()
}

func main() {
	// List of file names to read
	fileNames := []string{"test.text"}

	var wg sync.WaitGroup
	results := make(chan []string)
	errors := make(chan error)

	// Spawn a goroutine for each file
	for _, fileName := range fileNames {
		wg.Add(1)
		go readFile(fileName, &wg, results, errors)
	}

	go func() {
		wg.Wait()
		close(results)
		close(errors)
	}()

	// Process results and errors
	for result := range results {
		// Do something with the lines of each file
		fmt.Println(result)
	}

	for err := range errors {
		// Handle errors encountered during file reading
		fmt.Println("Error:", err)
	}
}
