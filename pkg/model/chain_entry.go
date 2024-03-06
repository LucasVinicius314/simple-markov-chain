package model

import "math/rand"

type ChainEntry struct {
	Conditions []MatchCondition `json:"conditions"`
	Next       map[string]int   `json:"next"`
	Token      string           `json:"token"`
	Total      int              `json:"total"`
}

func (e *ChainEntry) ComputeProbabilities() {
	for _, v := range e.Next {
		e.Total += v
	}

	conditions := []MatchCondition{}

	acc := 0
	for k, v := range e.Next {
		condition := MatchCondition{
			Max:   acc + v,
			Min:   acc,
			Token: k,
		}

		conditions = append(conditions, condition)

		acc = condition.Max
	}

	e.Conditions = conditions
}

func (e *ChainEntry) PickNext() *string {
	if e.Total == 0 {
		return nil
	}

	index := rand.Intn(e.Total)

	for _, condition := range e.Conditions {
		if index >= condition.Min && index < condition.Max {
			return &condition.Token
		}
	}

	return nil
}
