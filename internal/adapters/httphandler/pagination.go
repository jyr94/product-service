package httphandler

import (
	"net/http"
	"strconv"
)

type Pagination struct {
	Limit  int
	Offset int
}

func ParsePagination(r *http.Request, defaultLimit int) Pagination {
	query := r.URL.Query()

	limit := defaultLimit
	offset := 0

	if v := query.Get("limit"); v != "" {
		if l, err := strconv.Atoi(v); err == nil && l > 0 {
			limit = l
		}
	}

	if v := query.Get("offset"); v != "" {
		if o, err := strconv.Atoi(v); err == nil && o >= 0 {
			offset = o
		}
	}

	return Pagination{
		Limit:  limit,
		Offset: offset,
	}
}
