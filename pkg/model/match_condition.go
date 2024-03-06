package model

type MatchCondition struct {
	Max   int    `json:"max"`
	Min   int    `json:"min"`
	Token string `json:"token"`
}
