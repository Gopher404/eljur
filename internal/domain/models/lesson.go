package models

type Lesson struct {
	Id         int    `json:"id"`
	Week       int8   `json:"week"`
	WeekDay    int8   `json:"week_day"`
	Group      int8   `json:"group"`
	Number     int8   `json:"number"`
	Auditorium string `json:"auditorium"`
	Name       string `json:"name"`
	Teacher    string `json:"teacher"`
}
