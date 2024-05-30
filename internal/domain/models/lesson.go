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

type StructLessons struct {
	Weeks [2]struct {
		Days [6]struct {
			Lessons []Lesson `json:"lessons"`
		} `json:"days"`
	} `json:"weeks"`
}

func NewStructLessons() *StructLessons {
	ss := new(StructLessons)
	for w := range ss.Weeks {
		for d := range ss.Weeks[w].Days {
			ss.Weeks[w].Days[d].Lessons = make([]Lesson, 0)
		}
	}
	return ss
}
