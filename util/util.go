package util

type Config struct {
	Pubkey  string `arg:"" name:"address" help:"Nym address to send test message to." default:""`
	PrivKey bool   `help:"Send a binary file as the test message."`
}

func Contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
