package models

type Pagination struct {
	Page      int    `form:"page"`
	Size      int    `form:"size"`
	OrderBy   string `form:"order_by"`
	Ascending bool   `form:"asc"`
}

func NewPagination() *Pagination {
	return &Pagination{Page: 1, Size: 20, OrderBy: "", Ascending: false}
}

func (p *Pagination) Next() {
	p.Page++
}
