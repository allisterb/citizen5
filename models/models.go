package models

type Location struct {
	Lat  float32
	Long float32
}
type Report struct {
	Id          string
	Reporter    string
	Name        string
	Description string
	Location    Location
}
