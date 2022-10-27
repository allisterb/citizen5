package models

import (
	"github.com/allisterb/citizen5/expertai/hatespeech"
	"github.com/allisterb/citizen5/expertai/nlapi"
	"github.com/allisterb/citizen5/expertai/pii"
)

type Config struct {
	Pubkey  string
	PrivKey string
}

type Location struct {
	Lat  float32
	Long float32
}

type NLUAnalysis struct {
	Pii        pii.Response
	NL         nlapi.AnalyzeResponse
	HateSpeech hatespeech.HateSpeechDetectResponse
}

type VictimReport struct {
	Id            string
	DateSubmitted string
	Reporter      string
	Title         string
	Description   string
	Location      Location
}

type WitnessReport struct {
	Id               string
	DateSubmitted    string
	Reporter         string
	Title            string
	Text             string
	Date             string
	Location         Location
	GroupResponsible string
	Analysis         NLUAnalysis
}

type MediaReport struct {
	Id            string
	DateSubmitted string
	Reporter      string
	Text          string
	// D.C properties
	Contributor string
	Coverage    string
	Creator     string
	Date        string
	Description string
	Format      string
	Identifier  string
	Language    string
	Publisher   string
	Relation    string
	Rights      string
	Source      string
	Subject     string
	Title       string
	Type        string
	Analysis    NLUAnalysis
}
