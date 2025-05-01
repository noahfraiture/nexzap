package models

type SheetTempl struct {
	Id                string
	TutorialId        string
	Title             string
	CodeEditor        string
	SheetContent      string
	ExerciseContent   string
	SubmissionContent string
	NbPage            int
	MaxPage           int
	IsLast            bool
}

func NewSheetTempl(
	id string,
	tutorialId string,
	title string,
	codeEditor string,
	sheetContent string,
	exerciseContent string,
	submissionContent string,
	nbPage, maxPage int,
	isLast bool,
) SheetTempl {
	return SheetTempl{
		Id:                id,
		TutorialId:        tutorialId,
		Title:             title,
		CodeEditor:        codeEditor,
		SheetContent:      sheetContent,
		ExerciseContent:   exerciseContent,
		SubmissionContent: submissionContent,
		NbPage:            nbPage,
		MaxPage:           maxPage,
		IsLast:            isLast,
	}
}

type ListTutorialTempl struct {
	ID    string
	Title string
}

func NewListTutorial(id, title string) ListTutorialTempl {
	return ListTutorialTempl{
		ID:    id,
		Title: title,
	}
}
