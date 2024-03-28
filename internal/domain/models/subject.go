package models

type Subject struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Semester int8   `json:"semester"`
	Course   int8   `json:"course"`
}

type MinSubject struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type SubjectFindOpts struct {
	Id       *int    `json:"id"`
	Name     *string `json:"name"`
	Semester *int8   `json:"semester"`
	Course   *int8   `json:"course"`
}
