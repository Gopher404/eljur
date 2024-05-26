package models

type Lesson struct {
	Id         int    `json:"id"`
	Week       int8   `json:"week"`
	Number     int8   `json:"number"`
	WeekDay    int8   `json:"weekDay"`
	Auditorium string `json:"auditorium"`
	Name       string `json:"name"`
	Teacher    string `json:"teacher"`
}

type LessonForUpdate struct {
	Id         int8   `json:"id"`
	Auditorium string `json:"auditorium"`
	Name       string `json:"name"`
	Teacher    string `json:"teacher"`
}
