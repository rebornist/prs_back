package file

import "prs/models"

type OpinResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Opin    models.Opin `json:"opin"`
	Files   []string    `json:"files"`
	Token   string      `json:"csrf"`
}
