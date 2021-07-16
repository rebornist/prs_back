package models

import "time"

type Opin struct {
	RecordId        string    `json:"record_id"`
	FileId          string    `json:"file_id"`
	FileNm          string    `json:"file_nm"`
	NlpId           string    `json:"nlp_id"`
	RecordTitle     string    `json:"record_title"`
	Title           string    `json:"title"`
	ReclssDivCd     string    `json:"reclss_div_cd"`
	ReclssOpenGrade string    `json:"reclss_open_grade"`
	DivPart         string    `json:"div_part"`
	DivInfo         string    `json:"div_info"`
	Content         string    `json:"content"`
	Report          string    `json:"report"`
	Remarks         string    `json:"remarks"`
	UserId          string    `json:"user_id"`
	DocType         uint8     `json:"doc_type"`
	Permission      uint8     `json:"permission"`
	Complate        uint8     `json:"complate"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}
