package models

import "github.com/partyhall/partyhall/config"

type PaginatedResponse struct {
	Results      any `json:"results"`
	Page         int `json:"current_page"`
	PerPageCount int `json:"per_page_count"`
	TotalCount   int `json:"total_count"`
	PageCount    int `json:"page_count"`
}

func (pr *PaginatedResponse) CalculateMaxPage() {
	pr.PerPageCount = config.AMT_RESULTS_PER_PAGE

	if pr.TotalCount == 0 {
		pr.PageCount = 0
	} else {
		pr.PageCount = (pr.TotalCount + config.AMT_RESULTS_PER_PAGE - 1) / config.AMT_RESULTS_PER_PAGE
	}
}
