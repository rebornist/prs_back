package record

import "prs/models"

type RecordsResponse struct {
	TotalItems int             `json:"total_items"`
	Page       int             `json:"page"`
	Message    string          `json:"message"`
	Records    []models.Record `json:"records"`
	QueryParam string          `json:"query_param"`
	ViewNumber int             `json:"view_number"`
}
