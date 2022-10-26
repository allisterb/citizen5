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

type NLAnalysis struct {
	Pii        pii.Response
	NL         nlapi.AnalyzeResponse
	HateSpeech hatespeech.HateSpeechDetectResponse
}

type Report struct {
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
	Description      string
	Date             string
	Location         Location
	GroupResponsible string
	Analysis         NLAnalysis
}
