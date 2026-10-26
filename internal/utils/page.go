package utils

// paginationToSkipLimit converts pageNumber (1-based) and pageLimit to skip and limit values for MongoDB queries.
// If pageNumber <= 0, it defaults to 1.
// If pageLimit <= 0, limit is set to 0 meaning no limit.
func PaginationToSkipLimit(pageNumber, pageLimit int64) (skip int64, limit int64) {
	if pageNumber <= 0 {
		pageNumber = 1
	}
	if pageLimit <= 0 {
		return 0, 0 // no limit, no skip
	}
	skip = (pageNumber - 1) * pageLimit
	limit = pageLimit
	return skip, limit
}
