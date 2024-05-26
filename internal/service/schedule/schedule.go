package schedules

type ScheduleService struct {
}

func New() {

}

type Lesson struct {
	TimeStart  string // hour:minutes
	TimeEnd    string // hour:minutes
	LessonName string
	Teacher    string
}

type DaySchedule struct {
	WeekDay string
	Lessons []Lesson
}

var weekDays = map[string]string{
	"Mon": "Понедельник",
	"Tue": "Вторник",
	"Wed": "Среда",
	"Thu": "Четверг",
	"Fri": "Пятница",
	"Sat": "Суббота",
	"Sun": "Воскресенье",
}

func (s *ScheduleService) GetActualSchedule() (*[7]DaySchedule, error) {
	//now := time.Now()

	//weekDay := weekDays[now.Format("Mon")]

	return nil, nil
}
