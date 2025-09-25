package game

import (
	"bufio"
	"math/rand"
	"os"
	"strings"
)

func Getword(diff string) string {
	filePath := "./pkg/words/" + diff + ".txt"
	file, err := os.Open(filePath)
	if err != nil {
		return ""
	}
	defer file.Close()

	var words []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		word := strings.TrimSpace(scanner.Text())
		if word != "" {
			words = append(words, word)
		}
	}

	if len(words) == 0 {
		return ""
	}

	return words[rand.Intn(len(words))]
}
