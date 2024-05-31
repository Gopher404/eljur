package schedules

import (
	"context"
	"eljur/internal/config"
	"eljur/internal/domain/models"
	"eljur/internal/storage"
	"eljur/pkg/tr"
	"errors"
	"fmt"
	"slices"
	"strings"
	"time"
)

type UserGroupGetter interface {
	GetGroup(ctx context.Context, login string) (int8, error)
}

type ScheduleService struct {
	storage   *storage.Storage
	parser    *Parser
	groupName string
	ownerId   string
	user      UserGroupGetter
}

func New(storage *storage.Storage, user UserGroupGetter, config *config.ScheduleConfig) *ScheduleService {
	vk := newVKAPI(config.VKAPI.Version)
	parser := newParser(vk, config.VKSever, config.VKAPI.CacheTTL)
	return &ScheduleService{
		storage:   storage,
		user:      user,
		parser:    parser,
		groupName: config.GroupName,
		ownerId:   config.VKAPI.GroupId,
	}
}

type Lesson struct {
	Number     int8   `json:"number"`
	Name       string `json:"name"`
	Teacher    string `json:"teacher"`
	Auditorium string `json:"auditorium"`
	IsChange   bool   `json:"is_change"`
}

type DaySchedule struct {
	WeekDay int8   `json:"week_day"`
	Dates   string `json:"date"`
	date    time.Time
	Lessons []Lesson `json:"lessons"`
}

type WeekDay struct {
	Name string
	Num  int8
}

func newWeekSchedule() *WeekSchedule {
	s := new(WeekSchedule)
	for i := range s.Days {
		s.Days[i] = new(DaySchedule)
		s.Days[i].Lessons = make([]Lesson, 0)
	}
	return s
}

type WeekSchedule struct {
	CurrentDay int8            `json:"current_day"`
	WeekType   string          `json:"week_type"`
	Days       [7]*DaySchedule `json:"days"`
}

func (s *WeekSchedule) sort() {
	for _, d := range s.Days {
		slices.SortFunc(d.Lessons, func(a, b Lesson) int {
			return int(a.Number - b.Number)
		})
	}
}

func (s *WeekSchedule) addLessons(lessons []models.Lesson, userGroup int8) {
	for _, lesson := range lessons {
		if lesson.Group != userGroup && lesson.Group != 0 {
			continue
		}
		s.Days[lesson.WeekDay].Lessons = append(s.Days[lesson.WeekDay].Lessons, Lesson{
			Number:     lesson.Number,
			Name:       lesson.Name,
			Teacher:    lesson.Teacher,
			Auditorium: lesson.Auditorium,
			IsChange:   false,
		})
	}
}

func (s *WeekSchedule) Change(dayN int8, ch change) {
	if s == nil {
		return
	}

	for i := range s.Days[dayN].Lessons {
		if s.Days[dayN].Lessons[i].Number == ch.Number {
			updateLesson(&s.Days[dayN].Lessons[i], &Lesson{
				ch.Number, ch.Name, ch.Teacher, ch.Auditorium, true,
			})
			return
		}
	}
	if len(s.Days[dayN].Lessons) == 0 {
		s.Days[dayN].Lessons = append(s.Days[dayN].Lessons, Lesson{
			ch.Number, ch.Name, ch.Teacher, ch.Auditorium, true,
		})
		return
	}
	if ch.Number < s.Days[dayN].Lessons[0].Number {
		s.Days[dayN].Lessons = append([]Lesson{
			{ch.Number, ch.Name, ch.Teacher, ch.Auditorium, true},
		}, s.Days[dayN].Lessons...)
		return
	}
	if ch.Number > s.Days[dayN].Lessons[len(s.Days[dayN].Lessons)-1].Number {
		s.Days[dayN].Lessons = append(s.Days[dayN].Lessons, Lesson{
			ch.Number, ch.Name, ch.Teacher, ch.Auditorium, true,
		})
		return
	}

	for i := range s.Days[dayN].Lessons {
		if s.Days[dayN].Lessons[i].Number < ch.Number {
			s.Days[dayN].Lessons = append(append(s.Days[dayN].Lessons[:i], Lesson{
				ch.Number, ch.Name, ch.Teacher, ch.Auditorium, true,
			}), s.Days[dayN].Lessons[i:]...)
		}
	}
}

const dateLayout = "02.01.2006"

func updateLesson(lesson1, lesson2 *Lesson) {
	if lesson2.Name == "з/а" {
		lesson1.Auditorium = lesson2.Auditorium
		lesson1.IsChange = lesson2.IsChange
	} else if lesson2.Name == "ничего" {
		*lesson1 = Lesson{Number: lesson2.Number, Name: "Ничего", IsChange: lesson2.IsChange}
	} else {
		*lesson1 = *lesson2
	}
}

var weekDays = map[string]WeekDay{
	"Mon": {"Понедельник", 0},
	"Tue": {"Вторник", 1},
	"Wed": {"Среда", 2},
	"Thu": {"Четверг", 3},
	"Fri": {"Пятница", 4},
	"Sat": {"Суббота", 5},
	"Sun": {"Воскресенье", 6},
}

func (s *ScheduleService) GetActualSchedule(ctx context.Context, login string) (*WeekSchedule, error) {
	now := time.Now()

	weekDay := weekDays[now.Format("Mon")]
	timeToGetDoc := now
	if weekDay.Num == 6 {
		timeToGetDoc = now.Add(-1 * time.Hour * 24)
	}

	docsInf, err := s.parser.getListDocuments(s.ownerId)
	if err != nil {
		return nil, tr.Trace(err)
	}

	date := timeToGetDoc.Format(dateLayout)
	var nowDocInf *documentInfo
	for _, docInf := range docsInf {
		if strings.Index(docInf.Title, date) > -1 {
			nowDocInf = docInf
		}
	}
	if nowDocInf == nil {
		return nil, tr.Trace(errors.New("document not found"))
	}

	nowDoc, err := s.parser.getDocument(nowDocInf)
	if err != nil {
		return nil, tr.Trace(err)
	}
	week := s.parser.getWeekFromDocument(nowDoc)
	var weekN int8
	if week == "числитель" {
		weekN = 0
	} else {
		weekN = 1
	}

	userGroup, err := s.user.GetGroup(ctx, login)
	if err != nil {
		return nil, tr.Trace(err)
	}

	lessons, err := s.storage.Schedule.GetByWeek(ctx, weekN)
	if err != nil {
		return nil, tr.Trace(err)
	}

	weekSchedule := newWeekSchedule()
	weekSchedule.WeekType = week //string(unicode.ToUpper([]rune(week)[0])) + string([]rune(week)[1:])
	weekSchedule.CurrentDay = weekDay.Num

	weekSchedule.addLessons(lessons, userGroup)

	for i, day := range weekSchedule.Days {
		day.date = now.Add(time.Duration(int8(i)-weekDay.Num) * time.Hour * 24)
		day.Dates = day.date.Format(dateLayout)
		day.WeekDay = weekDays[day.date.Format("Mon")].Num
		//fmt.Println(weekDays[day.date.Format("Mon")], day.date.Format(time.DateTime))
	}

	lessons, err = s.storage.Schedule.GetByWeekAndDay(ctx, revertWeek(weekN), 0)
	if err != nil {
		return nil, tr.Trace(err)
	}
	for i := range lessons {
		lessons[i].WeekDay = 6
	}
	weekSchedule.addLessons(lessons, userGroup)
	weekSchedule.Days[6].date = weekSchedule.Days[5].date.Add(time.Hour * 24 * 2)
	weekSchedule.Days[6].Dates = weekSchedule.Days[6].date.Format(dateLayout)

	multiErr := newMultiError("set changes errors: ")
	for _, docInf := range docsInf {
		if docInf.Ext != "pdf" {
			continue
		}
		date, ok := s.parser.getDateFromDocInfo(docInf)
		if !ok {
			continue
		}
		dayOfWeek := weekDays[date.Format("Mon")]

		if now.Add(time.Duration(dayOfWeek.Num-weekDay.Num)*time.Hour*24).Format(dateLayout) != date.Format(dateLayout) {
			if date.Format(dateLayout) != weekSchedule.Days[6].Dates {
				continue
			}
			dayOfWeek = WeekDay{"Понедельник", 6}
		}

		doc, err := s.parser.getDocument(docInf)
		if err != nil {
			multiErr.AddError(err)
			continue
		}

		changes := s.parser.getChangesFromDocument(doc, s.groupName)
		for _, ch := range changes {
			weekSchedule.Change(dayOfWeek.Num, ch)
		}
	}
	weekSchedule.sort()
	if !multiErr.IsNil() {
		return weekSchedule, nil
	}
	return weekSchedule, multiErr
}

func revertWeek(week int8) int8 {
	if week == 0 {
		return 1
	}
	return 0
}

type LessonToSave struct {
	Action     string `json:"action"`
	Id         int    `json:"id"`
	Week       int8   `json:"week"`
	WeekDay    int8   `json:"week_day"`
	Number     int8   `json:"number"`
	Auditorium string `json:"auditorium"`
	Name       string `json:"name"`
	Teacher    string `json:"teacher"`
	Group      int8   `json:"group"`
}

func (s *ScheduleService) Save(ctx context.Context, saveList []LessonToSave) error {
	for _, lesson := range saveList {
		switch lesson.Action {
		case "new":
			if err := s.storage.Schedule.New(ctx, &models.Lesson{
				Week:       lesson.Week,
				WeekDay:    lesson.WeekDay,
				Number:     lesson.Number,
				Auditorium: lesson.Auditorium,
				Name:       lesson.Name,
				Teacher:    lesson.Teacher,
				Group:      lesson.Group,
			}); err != nil {
				return tr.Trace(err)
			}
			break
		case "update":
			if err := s.storage.Schedule.Update(ctx, &models.Lesson{
				Id:         lesson.Id,
				Week:       lesson.Week,
				WeekDay:    lesson.WeekDay,
				Number:     lesson.Number,
				Auditorium: lesson.Auditorium,
				Name:       lesson.Name,
				Teacher:    lesson.Teacher,
				Group:      lesson.Group,
			}); err != nil {
				return tr.Trace(err)
			}
			break
		case "del":
			if err := s.storage.Schedule.Delete(ctx, lesson.Id); err != nil {
				return tr.Trace(err)
			}
			break
		}
	}
	return nil
}

func (s *ScheduleService) GetAll(ctx context.Context) (*models.StructLessons, error) {
	ss := models.NewStructLessons()

	lessons, err := s.storage.Schedule.GetAll(ctx)
	if err != nil {
		return nil, tr.Trace(err)
	}
	for _, lesson := range lessons {
		ss.Weeks[lesson.Week].Days[lesson.WeekDay].Lessons =
			append(ss.Weeks[lesson.Week].Days[lesson.WeekDay].Lessons, lesson)
	}
	sortStructLessons(ss)
	return ss, nil
}

func sortStructLessons(ss *models.StructLessons) {
	for w := range ss.Weeks {
		for d := range ss.Weeks[w].Days {
			slices.SortFunc(ss.Weeks[w].Days[d].Lessons, func(a, b models.Lesson) int {
				return int(a.Number - b.Number)
			})
		}
	}
}

func newMultiError(header string) *MultiError {
	return &MultiError{
		value: header,
	}
}

type MultiError struct {
	value  string
	count  int
	header string
}

func (e *MultiError) Error() string {
	return e.value
}

func (e *MultiError) AddError(err error) {
	e.count++
	e.value += fmt.Sprintf("err%d: %s, ", e.count, err.Error())
}
func (e *MultiError) IsNil() bool {
	return e.value == e.header
}
