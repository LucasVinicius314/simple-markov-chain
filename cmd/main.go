package main

import (
	"fmt"
	"math/rand"
	"reflect"
	"regexp"
	"strings"

	"simple-markov-chain/pkg/model"
	"simple-markov-chain/pkg/util"
)

func main() {
	inputPath := "../resources/input.txt"
	outputPath := "../resources/chain.json"

	lines, err := util.ReadLines(inputPath)
	if err != nil {
		fmt.Printf("Error reading file: %v", err)
		return
	}

	tokens := []string{}

	for _, line := range lines {
		trimmed := strings.TrimSpace(strings.ToLower(regexp.MustCompile(`\s+`).ReplaceAllString(line, " ")))
		if trimmed == "" {
			continue
		}

		tokens = append(tokens, strings.Split(trimmed, " ")...)
	}

	chain := map[string]*model.ChainEntry{}

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

	err = util.SaveChain(chain, outputPath)
	if err != nil {
		fmt.Printf("Error saving chain: %v", err)
		return
	}

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

	fmt.Println(strings.TrimSpace(out))
}
