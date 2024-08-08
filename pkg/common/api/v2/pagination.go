package v2

type Pagination struct {
	PageNumber int
	PerPage int
	OrderBy string
	Ascending bool
}

func NewPagination() *Pagination {
	return &Pagination{PageNumber: 1, PerPage: 20, OrderBy: "", Ascending: false}
}

func (p *Pagination) Next() {
	p.PageNumber++
}
