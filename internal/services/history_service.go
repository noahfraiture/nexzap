package services

import (
	"context"
	"nexzap/internal/db"
	generated "nexzap/internal/db/generated"
)

type HistoryService struct {
	db *db.Database
}

func NewHistoryService(database *db.Database) *HistoryService {
	return &HistoryService{
		db: database,
	}
}

type ListTutorials = generated.ListTutorialsRow

func (h *HistoryService) ListTutorials() ([]ListTutorials, error) {
	return h.db.GetRepository().ListTutorials(context.Background())
}
