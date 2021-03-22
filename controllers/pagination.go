package controllers

import (
	"net/http"
	"strconv"

	. "notes/config"

	"gorm.io/gorm"
)

//Paginate scope
func Paginate(r *http.Request) func(db *gorm.DB) *gorm.DB {
	page := GetPage(r)
	return func(db *gorm.DB) *gorm.DB {
		return db.Offset((page - 1) * Cfg.PerPage).Limit(Cfg.PerPage)
	}
}

// PaginationData ...
func PaginationData(req *http.Request, db *gorm.DB) map[string]interface{} {
	var count int64
	err := db.Count(&count).Error
	if err != nil {
		panic(err)
	}
	pag := map[string]interface{}{
		"current_page": GetPage(req),
		"max_page":     (int(count)-1)/Cfg.PerPage + 1,
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
