package models

type SheetTempl struct {
	Id                string
	Title             string
	CodeEditor        string
	SheetContent      string
	ExerciseContent   string
	SubmissionContent string
	NbPage            int
	MaxPage           int
}

func NewSheetTempl(
	id string,
	title string,
	codeEditor string,
	sheetContent string,
	exerciseContent string,
	submissionContent string,
	nbPage, maxPage int,
) SheetTempl {
	return SheetTempl{
		Id:                id,
		Title:             title,
		CodeEditor:        codeEditor,
		SheetContent:      sheetContent,
		ExerciseContent:   exerciseContent,
		SubmissionContent: submissionContent,
		NbPage:            nbPage,
		MaxPage:           maxPage,
	}
}
