package models

type Tutorial struct {
	Language    string      `toml:"language"`
	Sheets      []Paragraph `toml:"sheets"`
	Highlighter string      `toml:"highlighter"`
}

type Paragraph struct {
	Content string `toml:"content"`
	IsCode  bool   `toml:"is_code"`
}
