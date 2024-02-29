package models

type Grade struct {
	Id        int    `json:"id"`
	UserId    int    `json:"user_id"`
	SubjectId int    `json:"subject_id"`
	Value     string `json:"value"`
	Day       int8   `json:"day"`
	Month     int8   `json:"month"`
	Course    int8   `json:"course"`
}

type MinGrade struct {
	Id    int    `json:"id"`
	Value string `json:"value"`
}

type GradesFindOpts struct {
	Id        *int
	UserId    *int
	SubjectId *int
	Value     *string
	Day       *int8
	Month     *int8
	Course    *int8
}
