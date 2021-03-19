package controllers

import (
	"fmt"
	"notes/models"

	. "notes/config"

	"gorm.io/gorm"
)

// Page return certain page of notes
func Page(notes *gorm.DB, page int) ([]models.Note, error) {
	res := []models.Note{}
	err := notes.Offset((page - 1) * Cfg.PerPage).Limit(Cfg.PerPage).Find(&res).Error
	if err != nil {
		return nil, fmt.Errorf("when fetching page: %v", err)
	}
	return res, nil
}

// PaginationData ...
func PaginationData(cur, total int) map[string]int {
	pag := map[string]int{
		"current_page": cur,
		"max_page":     total,
		"per_page":     Cfg.PerPage,
	}
	return pag
}
