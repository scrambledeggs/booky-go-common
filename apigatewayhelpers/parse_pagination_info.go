package apigatewayhelpers

import "strconv"

func InitPagination(queryString map[string]string) (int32, int32, map[string]any) {
	page, err := strconv.ParseInt(queryString["page"], 10, 32)
	if err != nil {
		page = 1
	}

	resultsPerPage, err := strconv.ParseInt(queryString["results_per_page"], 10, 32)
	if err != nil {
		resultsPerPage = 10
	}

	offset := (page - 1) * resultsPerPage

	return int32(offset), int32(resultsPerPage), map[string]any{"page": page, "results_per_page": resultsPerPage}
}
