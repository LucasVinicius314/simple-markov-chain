package util

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"simple-markov-chain/pkg/model"
)

func ReadLines(fileName string) ([]string, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines, scanner.Err()
}

func SaveChain(chain map[string]*model.ChainEntry, outputPath string) error {
	data, err := json.Marshal(chain)
	if err != nil {
		fmt.Printf("Error marshaling chain: %v", err)
		return err
	}

	err = os.WriteFile(outputPath, data, 0644)
	if err != nil {
		fmt.Printf("Error writing file: %v", err)
		return err
	}

	return nil
}
