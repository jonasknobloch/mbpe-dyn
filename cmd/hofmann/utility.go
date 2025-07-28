package main

import (
	"encoding/json"
	"golang.org/x/exp/constraints"
	"os"
	"slices"
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

func getKeys[K constraints.Ordered, V any](m map[K]V) []K {
	keys := make([]K, 0, len(m))

	for k := range m {
		keys = append(keys, k)
	}

	slices.Sort(keys)

	return keys
}

func toSet[T comparable](slice []T) map[T]struct{} {
	s := make(map[T]struct{}, len(slice))

	for _, v := range slice {
		s[v] = struct{}{}
	}

	return s
}
