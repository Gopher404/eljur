package subjects

import "eljur/internal/storage"

type SubjectService struct {
	subjectStorage storage.Subjects
}

func New(subjectStorage storage.Subjects) *SubjectService {
	return &SubjectService{
		subjectStorage: subjectStorage,
	}
}
