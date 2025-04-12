package models

type Tutorial struct {
	Language string
	Sheets   []Sheet
}

type Sheet struct {
	Tests   string
	Content string
}
