package models

type Record struct {
	RecordId    string `json:"record_id"`
	ProcDepNm   string `json:"proc_dep_nm"`
	ClssNm      string `json:"clss_nm"`
	FolderTitle string `json:"folder_title"`
	Title       string `json:"title"`
	OpenDivCd   string `json:"open_div_cd"`
	OpenGrade   string `json:"open_grade"`
	Complate    uint8  `json:"complate"`
}
