package models

type SheetTempl struct {
	Id                string
	Title             string
	Highlight         string
	CodeEditor        string
	SheetContent      string
	ExerciseContent   string
	CorrectionContent string
	NbPage            int
	MaxPage           int
}

func NewSheetTempl(
	id string,
	title string,
	highlight string,
	codeEditor string,
	sheetContent string,
	exerciseContent string,
	correctionContent string,
	nbPage, maxPage int,
) SheetTempl {
	return SheetTempl{
		Id:                id,
		Title:             title,
		Highlight:         highlight,
		CodeEditor:        codeEditor,
		SheetContent:      sheetContent,
		ExerciseContent:   exerciseContent,
		CorrectionContent: correctionContent,
		NbPage:            nbPage,
		MaxPage:           maxPage,
	}
}
