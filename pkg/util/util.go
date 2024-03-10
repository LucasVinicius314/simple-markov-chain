package util

import (
	"encoding/json"
	"fmt"
	"os"
	"simple-markov-chain/pkg/model"
)

func FileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func ReadFile(path string) ([]byte, error) {
	return os.ReadFile(path)
}

func SaveChain(chain map[string]*model.ChainEntry, path string) error {
	data, err := json.Marshal(chain)
	if err != nil {
		fmt.Printf("Error marshaling chain: %v\n", err)
		return err
	}

	err = os.WriteFile(path, data, 0644)
	if err != nil {
		fmt.Printf("Error writing file: %v\n", err)
		return err
	}

	return nil
}
