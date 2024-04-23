package subjects

import (
	"context"
	"eljur/internal/domain/models"
	"eljur/internal/storage"
	"eljur/pkg/tr"
)

type GradesSubjects interface {
	NewResGradesBySubject(ctx context.Context, subjectId int, course int8) error
	DeleteBySubject(ctx context.Context, subjectId int) error
}

type SubjectService struct {
	storage *storage.Storage
	grades  GradesSubjects
}

func New(storage *storage.Storage, grades GradesSubjects) *SubjectService {
	return &SubjectService{
		storage: storage,
		grades:  grades,
	}
}

var (
	subjectResSemester int8 = 3
)

func (s *SubjectService) GetSubject(ctx context.Context, id int) (string, error) {
	subject, err := s.storage.Subjects.GetById(ctx, id)
	if err != nil {
		return "", tr.Trace(err)
	}
	return subject.Name, nil
}

func (s *SubjectService) GetAllSubjects(ctx context.Context) (*[4][3][]models.MinSubject, error) {
	subjects, err := s.storage.Subjects.GetAll(ctx)
	var structSubjects [4][3][]models.MinSubject
	for _, subject := range subjects {
		minSubject := models.MinSubject{Id: subject.Id, Name: subject.Name}
		structSubjects[subject.Course-1][subject.Semester-1] = append(
			structSubjects[subject.Course-1][subject.Semester-1],
			minSubject,
		)
	}
	if err != nil {
		return nil, tr.Trace(err)
	}
	return &structSubjects, nil
}

func (s *SubjectService) GetBySemester(ctx context.Context, semester int8, course int8) ([]models.MinSubject, error) {
	subjects, err := s.storage.Subjects.GetBySemester(ctx, semester, course)
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

func (s *SubjectService) Save(ctx context.Context, subjects []SaveSubjectIn) error {
	ctx, err := s.storage.Tx.Begin(ctx)
	if err != nil {
		return tr.Trace(err)
	}

	for _, subject := range subjects {
		switch subject.Action {
		case "new":
			if err := s.newSubject(ctx, models.Subject{
				Name:     subject.Name,
				Semester: subject.Semester,
				Course:   subject.Course,
			}); err != nil {
				return tr.Trace(err)
			}
			break
		case "update":
			if err := s.update(ctx, models.MinSubject{
				Id:   subject.Id,
				Name: subject.Name,
			}); err != nil {
				return tr.Trace(err)
			}
			break
		case "del":
			if err := s.delete(ctx, subject.Id); err != nil {
				return tr.Trace(err)
			}
			break
		}
	}
	if err := s.storage.Tx.Commit(ctx); err != nil {
		return tr.Trace(err)
	}
	return nil
}

func (s *SubjectService) newSubject(ctx context.Context, subject models.Subject) error {

	id, err := s.storage.Subjects.NewSubject(ctx, subject)
	if err != nil {
		return tr.Trace(err)
	}

	if subject.Semester == subjectResSemester {
		if err := s.grades.NewResGradesBySubject(ctx, id, subject.Course); err != nil {
			return tr.Trace(err)
		}
	}

	return nil
}

func (s *SubjectService) update(ctx context.Context, subject models.MinSubject) error {
	if err := s.storage.Subjects.Update(ctx, subject); err != nil {
		return tr.Trace(err)
	}
	return nil
}

func (s *SubjectService) delete(ctx context.Context, id int) error {
	if err := s.storage.Subjects.Delete(ctx, id); err != nil {
		return tr.Trace(err)
	}
	if err := s.grades.DeleteBySubject(ctx, id); err != nil {
		return tr.Trace(err)
	}

	return nil
}
