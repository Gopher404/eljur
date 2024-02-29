package subjects

import (
	"eljur/internal/domain/models"
	"eljur/internal/storage"
	"eljur/pkg/tr"
)

type SubjectService struct {
	subjectStorage storage.Subjects
}

func New(subjectStorage storage.Subjects) *SubjectService {
	return &SubjectService{
		subjectStorage: subjectStorage,
	}
}

func (s *SubjectService) GetSubject(id int) (string, error) {
	name, err := s.subjectStorage.GetById(id)
	if err != nil {
		return "", tr.Trace(err)
	}
	return name, nil
}

func (s *SubjectService) GetAllSubjects() ([]models.Subject, error) {
	subjects, err := s.subjectStorage.GetAll()
	if err != nil {
		return nil, tr.Trace(err)
	}
	return subjects, nil
}
