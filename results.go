package mbpe

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func WalkResults(root, match string) ([]string, error) {
	paths := make([]string, 0)

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.Contains(path, match) {
			paths = append(paths, path)
		}

		return nil
	})

	return paths, err
}

func WalkResultsStatic(format string) ([]string, [][3]string) {
	vocabSizes := []int{8192, 16384, 32768, 50256, 100512}
	prefixes := []string{"m", "mi"}
	alphas := []float64{0.0, 0.1, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8, 0.9, 1.0}

	paths := make([]string, 0, len(vocabSizes)*len(prefixes)*len(alphas))
	stubs := make([][3]string, 0, len(vocabSizes)*len(prefixes)*len(alphas))

	for _, vocab := range vocabSizes {
		for _, prefix := range prefixes {
			for _, alpha := range alphas {
				path := fmt.Sprintf(format, vocab, prefix, fmt.Sprintf("%03d", int(alpha*100)))

				stub := [3]string{strconv.Itoa(vocab), prefix, fmt.Sprintf("%.2f", alpha)}

				paths = append(paths, path)
				stubs = append(stubs, stub)
			}
		}
	}

	return paths, stubs
}
