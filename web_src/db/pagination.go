package db

const maxRowsPerPage uint64 = 100

type PageResult struct {
	Items interface{} `json:"items,omitempty"`
	Count int64       `json:"count"`
}

type Pagination struct {
	Page    uint64 `json:"page,omitempty"` // start from 0
	PerPage uint64 `json:"perPage,omitempty"`
}

func NewPagination(page, perPage uint64) Pagination {
	return Pagination{
		Page:    page,
		PerPage: perPage,
	}
}

func (p *Pagination) filter() (page, perPage uint64) {
	page, perPage = p.Page+1, p.PerPage

	if perPage == 0 || perPage > maxRowsPerPage {
		perPage = maxRowsPerPage
	}

	p.Page, p.PerPage = page, perPage

	return
}
