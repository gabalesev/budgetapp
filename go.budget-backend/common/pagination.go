package common

import (
	"math"
	"net/http"
	"strconv"

	"gorm.io/gorm"
)

type Pagination struct {
	PageSize   int         `query:"page_size" json:"page_size"`
	PageNumber int         `query:"page" json:"page_number"`
	TotalItems int64       `query:"total_items" json:"total_items"`
	TotalPages int         `query:"total_pages" json:"total_pages"`
	Sort       string      `query:"sort"`
	Items      interface{} `json:"items"`
}

func (p *Pagination) GetPage() int {
	if p.PageNumber <= 0 {
		p.PageNumber = 1
	}

	return p.PageNumber
}

func (p *Pagination) GetPageSize() int {
	if p.PageSize > 100 {
		p.PageSize = 100
	} else if p.PageSize <= 0 {
		p.PageSize = 10
	}

	return p.PageSize
}

func (p *Pagination) GetOffset() int {
	return (p.GetPage() - 1) * p.GetPageSize()
}

func NewPagination(model interface{}, r *http.Request, db *gorm.DB) *Pagination {
	var pagination Pagination
	q := r.URL.Query()

	pageNumber, _ := strconv.Atoi(q.Get("page_number"))
	pageSize, _ := strconv.Atoi(q.Get("page_size"))

	var totalItems int64
	db.Model(model).Count(&totalItems)

	pagination.PageNumber = pageNumber
	pagination.PageSize = pageSize
	pagination.TotalItems = totalItems

	totalPages := int(math.Ceil(float64(totalItems) / float64(pagination.GetPageSize())))

	pagination.TotalPages = totalPages

	return &pagination
}

func (p *Pagination) Paginate() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Offset(p.GetOffset()).Limit(p.GetPageSize())
	}
}
