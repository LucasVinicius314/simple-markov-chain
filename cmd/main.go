package main

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"reflect"
	"regexp"
	"strings"

	"simple-markov-chain/pkg/model"
	"simple-markov-chain/pkg/util"
)

func main() {
	dataFilePath := ``
	if len(dataFilePath) > 0 {
		err := parseDataFile(dataFilePath)
		if err != nil {
			fmt.Printf("Error parsing data file: %v\n", err)
			return
		}
	}

	inputFileName := "input.txt"

	inputBytes, err := util.ReadFile(fmt.Sprintf("../resources/input/%s", inputFileName))
	if err != nil {
		fmt.Printf("Error reading input file: %v\n", err)
		return
	}

	hashBytes := []byte{}
	for _, v := range sha256.Sum256(inputBytes) {
		hashBytes = append(hashBytes, v)
	}

	chain := map[string]*model.ChainEntry{}

	chainPath := fmt.Sprintf("../resources/output/%s.json", hex.EncodeToString(hashBytes))
	if util.FileExists(chainPath) {
		fmt.Printf("Loading chain [%s].\n", chainPath)

		chainBytes, err := util.ReadFile(chainPath)
		if err != nil {
			fmt.Printf("Error reading chain file: %v\n", err)
			return
		}

		err = json.Unmarshal(chainBytes, &chain)
		if err != nil {
			fmt.Printf("Error unmarshaling chain: %v\n", err)
			return
		}
	} else {
		fmt.Printf("Generating chain from [%s].\n", inputFileName)

		tokens := []string{}

		lines := strings.Split(string(inputBytes), "\n")
		for _, line := range lines {
			trimmed := strings.TrimSpace(strings.ToLower(regexp.MustCompile(`\s+`).ReplaceAllString(line, " ")))
			if trimmed == "" {
				continue
			}

			tokens = append(tokens, strings.Split(trimmed, " ")...)
		}

		lastToken := ""
		for _, token := range tokens {
			if lastToken == "" {
				lastToken = token
				continue
			}

			if _, ok := chain[lastToken]; !ok {
				chain[lastToken] = &model.ChainEntry{
					Next:  map[string]int{},
					Token: lastToken,
				}
			}

			chain[lastToken].Next[token]++

			lastToken = token
		}

		for _, v := range chain {
			v.ComputeProbabilities()
		}

		err = util.SaveChain(chain, chainPath)
		if err != nil {
			fmt.Printf("Error saving chain [%s]: %v\n", chainPath, err)
			return
		}
	}

	fmt.Println(generateOutput(chain))
	fmt.Println(generateOutput(chain))
	fmt.Println(generateOutput(chain))
	fmt.Println(generateOutput(chain))
	fmt.Println(generateOutput(chain))
}

type DataFileEntry struct {
	Contents string `json:"Contents"`
}

func parseDataFile(filePath string) error {
	fileBytes, err := util.ReadFile(filePath)
	if err != nil {
		return err
	}

	dataFileEntries := []DataFileEntry{}
	err = json.Unmarshal(fileBytes, &dataFileEntries)
	if err != nil {
		return err
	}

	file, err := os.Create("../resources/output.txt")
	if err != nil {
		return err
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	for _, entry := range dataFileEntries {
		trimmedLine := strings.TrimSpace(entry.Contents)
		if trimmedLine == "" {
			continue
		}

		_, err := w.WriteString(trimmedLine + "\n")
		if err != nil {
			return err
		}
	}
	w.Flush()

	return nil
}

func generateOutput(chain map[string]*model.ChainEntry) string {
	out := ""

	keys := reflect.ValueOf(chain).MapKeys()
	currentEntry := chain[keys[rand.Intn(len(keys))].String()]

	stoppingChance := .05

	for {
		out = fmt.Sprintf("%s %s", out, currentEntry.Token)

		nextToken := currentEntry.PickNext()
		if nextToken == nil {
			break
		}

		nextEntry := chain[*nextToken]
		if nextEntry == nil {
			break
		}

		currentEntry = nextEntry

		if rand.Float64() < stoppingChance {
			break
		}
	}

	return strings.TrimSpace(out)
}
