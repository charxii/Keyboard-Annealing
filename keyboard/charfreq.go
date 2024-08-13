package keyboard

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type CharFreq struct {
	Chars    map[rune]int
	Bigrams  map[string]int
	Trigrams map[string]int
}

func NewCharFreq(path string) (*CharFreq, error) {
	chars := make(map[rune]int)
	bigrams := make(map[string]int)
	trigrams := make(map[string]int)

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		if len(line) < 2 {
			continue
		}

		for i := 0; i < len(line)-1; i++ {
			bigram := line[i : i+2]
			bigrams[bigram]++

			if i+3 <= len(line) {
				trigram := line[i : i+3]
				trigrams[trigram]++
			}
		}

		for _, r := range line {
			chars[r]++
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
	}

	return &CharFreq{chars, bigrams, trigrams}, nil
}

func CharFreqFromFolder(path string) (*CharFreq, error) {
	cf := &CharFreq{
		Chars:    make(map[rune]int),
		Bigrams:  make(map[string]int),
		Trigrams: make(map[string]int),
	}

	readFile := func(filename string, m interface{}) error {
		file, err := os.Open(filename)
		if err != nil {
			return err
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			parts := strings.Fields(scanner.Text())
			if len(parts) != 2 {
				continue
			}

			key := parts[0]
			value, err := strconv.Atoi(parts[1])
			if err != nil {
				continue
			}

			switch v := m.(type) {
			case map[rune]int:
				if len(key) == 1 {
					v[rune(key[0])] = value
				}
			case map[string]int:
				v[key] = value
			}
		}

		return scanner.Err()
	}

	if err := readFile(path+"/monograms.txt", cf.Chars); err != nil {
		return nil, err
	}
	if err := readFile(path+"/bigrams.txt", cf.Bigrams); err != nil {
		return nil, err
	}
	if err := readFile(path+"/trigrams.txt", cf.Trigrams); err != nil {
		return nil, err
	}

	return cf, nil
}
