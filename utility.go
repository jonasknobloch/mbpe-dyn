package main

import (
	"bufio"
	"encoding/json"
	"os"
	"sync"
)

func fromFile(name string, callback func(scanner *bufio.Scanner) error) error {
	file, err := os.Open(name)

	if err != nil {
		return err
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	if err := callback(scanner); err != nil {
		return err
	}

	return nil
}

func toFile(name string, callback func(writer *bufio.Writer) error) error {
	file, err := os.Create(name)

	if err != nil {
		return err
	}

	defer file.Close()

	writer := bufio.NewWriter(file)

	if err := callback(writer); err != nil {
		return err
	}

	if err = writer.Flush(); err != nil {
		return err
	}

	return nil
}

func fromJSON(name string, data interface{}) error {
	file, err := os.Open(name)

	if err != nil {
		return err
	}

	defer file.Close()

	decoder := json.NewDecoder(file)

	if err := decoder.Decode(data); err != nil {
		return err
	}

	return nil
}

func toJSON(name string, data interface{}) error {
	file, err := os.Create(name)

	if err != nil {
		return err
	}

	defer file.Close()

	encoder := json.NewEncoder(file)

	if err := encoder.Encode(data); err != nil {
		return err
	}

	return nil
}

func countLines(names ...string) (int, error) {
	var wg sync.WaitGroup

	results := make(chan int, len(names))
	errors := make(chan error, len(names))

	// Start a goroutine for each file
	for _, name := range names {
		wg.Add(1)
		go func() {
			defer wg.Done()

			var scanner *bufio.Scanner

			if file, err := os.Open(name); err != nil {
				results <- 0
				errors <- err

				return
			} else {
				scanner = bufio.NewScanner(file)

				defer file.Close()
			}

			count := 0

			for scanner.Scan() {
				count++
			}

			if err := scanner.Err(); err != nil {
				results <- 0
				errors <- err

				return
			}

			results <- count
		}()
	}

	wg.Wait()

	close(results)
	close(errors)

	total := 0

	for count := range results {
		total += count
	}

	return total, <-errors
}
