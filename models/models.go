package models

type Location struct {
	Lat  float32
	Long float32
}
type Report struct {
	Id            string
	DateSubmitted string
	Reporter      string
	Title         string
	Description   string
	Location      Location
}
