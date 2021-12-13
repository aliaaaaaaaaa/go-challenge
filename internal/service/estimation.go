package service

import (
	"es/internal/model"
	"es/internal/repo"
)

type EstimationService struct {
	estimationRepo repo.EstimationRepo
}

func NewEstimationService(estimationRepo repo.EstimationRepo) *EstimationService {
	return &EstimationService{estimationRepo: estimationRepo}
}

func (e *EstimationService) SaveSegmentTagForUser(userid uint32, segment string) error {
	var estimation = model.Estimation{
		UserId:  userid,
		Segment: segment,
	}

	e.estimationRepo.SaveSegmentTagForUser(estimation)
	return nil
}
func (e *EstimationService) GetSegmentTagFor14dLastDays(segment string) (uint32, error) {
	count, err := e.estimationRepo.GetSegmentTagFor14dLastDays(segment)
	return count, err

}
