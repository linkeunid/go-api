package pagination

import (
	"math"
	"net/http"
	"strconv"
)

// DefaultLimit is the default number of items per page
const DefaultLimit = 10

// MaxLimit is the maximum number of items per page
const MaxLimit = 100

// Params represents pagination parameters
type Params struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	TotalItems int64 `json:"total_items"`
	TotalPages int   `json:"total_pages"`
}

// PagedData represents a paginated data response
type PagedData struct {
	Items interface{} `json:"items"`
	Meta  Params      `json:"meta"`
}

// NewParams creates a new pagination parameters from HTTP request
func NewParams(r *http.Request) Params {
	query := r.URL.Query()

	// Parse page number
	page, err := strconv.Atoi(query.Get("page"))
	if err != nil || page < 1 {
		page = 1
	}

	// Parse items per page
	limit, err := strconv.Atoi(query.Get("limit"))
	if err != nil || limit < 1 {
		limit = DefaultLimit
	}

	// Enforce maximum limit
	if limit > MaxLimit {
		limit = MaxLimit
	}

	return Params{
		Page:  page,
		Limit: limit,
	}
}

// GetOffset returns the offset for database queries
func (p Params) GetOffset() int {
	return (p.Page - 1) * p.Limit
}

// CalculatePages calculates total pages based on total items
func (p *Params) CalculatePages(totalItems int64) {
	p.TotalItems = totalItems
	p.TotalPages = int(math.Ceil(float64(totalItems) / float64(p.Limit)))
}

// HasPreviousPage returns true if there is a previous page
func (p Params) HasPreviousPage() bool {
	return p.Page > 1
}

// HasNextPage returns true if there is a next page
func (p Params) HasNextPage() bool {
	return p.Page < p.TotalPages
}

// GetPreviousPage returns the previous page number
func (p Params) GetPreviousPage() int {
	if !p.HasPreviousPage() {
		return p.Page
	}
	return p.Page - 1
}

// GetNextPage returns the next page number
func (p Params) GetNextPage() int {
	if !p.HasNextPage() {
		return p.Page
	}
	return p.Page + 1
}
