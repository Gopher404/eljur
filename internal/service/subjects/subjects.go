package subjects

import (
	"eljur/internal/domain/models"
	"eljur/internal/storage"
	"eljur/pkg/tr"
)

type GradesSubjects interface {
	NewResGradesBySubject(subjectId int, semester int8, course int8) error
	DeleteBySubject(subjectId int) error
}

type SubjectService struct {
	subjectStorage storage.Subjects
	grades         GradesSubjects
}

func New(subjectStorage storage.Subjects, grades GradesSubjects) *SubjectService {
	return &SubjectService{
		subjectStorage: subjectStorage,
		grades:         grades,
	}
}

func (s *SubjectService) GetSubject(id int) (string, error) {
	name, err := s.subjectStorage.GetById(id)
	if err != nil {
		return "", tr.Trace(err)
	}
	return name, nil
}

func (s *SubjectService) GetAllSubjects() (*[4][2][]models.MinSubject, error) {
	subjects, err := s.subjectStorage.GetAll()
	var structSubjects [4][2][]models.MinSubject
	for _, subject := range subjects {
		structSubjects[subject.Course-1][subject.Semester-1] = append(structSubjects[subject.Course-1][subject.Semester-1],
			models.MinSubject{
				Id:   subject.Id,
				Name: subject.Name,
			})
	}
	if err != nil {
		return nil, tr.Trace(err)
	}
	return &structSubjects, nil
}

func (s *SubjectService) GetBySemester(semester int8, course int8) ([]models.MinSubject, error) {
	subjects, err := s.subjectStorage.GetBySemester(semester, course)
	if err != nil {
		return nil, tr.Trace(err)
	}
	return subjects, nil
}

type SaveSubjectIn struct {
	Action   string `json:"action"`
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Semester int8   `json:"semester"`
	Course   int8   `json:"course"`
}

func (s *SubjectService) Save(subjects []SaveSubjectIn) error {
	for _, subject := range subjects {
		switch subject.Action {
		case "new":
			if err := s.newSubject(models.Subject{
				Name:     subject.Name,
				Semester: subject.Semester,
				Course:   subject.Course,
			}); err != nil {
				return tr.Trace(err)
			}
			break
		case "update":
			if err := s.update(models.MinSubject{
				Id:   subject.Id,
				Name: subject.Name,
			}); err != nil {
				return tr.Trace(err)
			}
			break
		case "del":
			if err := s.delete(subject.Id); err != nil {
				return tr.Trace(err)
			}
			break
		}
	}
	return nil
}

func (s *SubjectService) newSubject(subject models.Subject) error {
	id, err := s.subjectStorage.NewSubject(subject)
	if err != nil {
		return tr.Trace(err)
	}
	if err := s.grades.NewResGradesBySubject(id, subject.Semester, subject.Course); err != nil {
		return tr.Trace(err)
	}
	return nil
}

func (s *SubjectService) update(subject models.MinSubject) error {
	if err := s.subjectStorage.Update(subject); err != nil {
		return tr.Trace(err)
	}
	return nil
}

func (s *SubjectService) delete(id int) error {
	if err := s.subjectStorage.Delete(id); err != nil {
		return tr.Trace(err)
	}
	if err := s.grades.DeleteBySubject(id); err != nil {
		return tr.Trace(err)
	}
	return nil
}
