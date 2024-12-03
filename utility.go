package main

import (
	"bufio"
	"encoding/json"
	"os"
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
