package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	. "notes/config"

	"gorm.io/gorm"
)

// Page return certain page of notes
func Page(notes *gorm.DB, page int) ([]map[string]interface{}, error) {
	res := []map[string]interface{}{}
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

// GetPage extracts page from request
func GetPage(r *http.Request) int {
	pageStr := r.URL.Query().Get("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		page = 1
	}

	return page
}
