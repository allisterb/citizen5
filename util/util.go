package util

import (
	"os"
)

type Config struct {
	Pubkey  string `arg:"" name:"address" help:"Nym address to send test message to." default:""`
	PrivKey bool   `help:"Send a binary file as the test message."`
}

func GetUserHomeDir() string {
	h, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	return h
}

func Contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
