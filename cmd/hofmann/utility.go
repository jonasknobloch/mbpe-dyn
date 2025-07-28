package main

import (
	"encoding/json"
	"os"
	"sort"
)

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

func getKeys(m map[string]interface{}) []string {
	keys := make([]string, 0, len(m))

	for k := range m {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	return keys
}

func toSet(slice []string) map[string]struct{} {
	s := make(map[string]struct{}, len(slice))

	for _, v := range slice {
		s[v] = struct{}{}
	}

	return s
}
