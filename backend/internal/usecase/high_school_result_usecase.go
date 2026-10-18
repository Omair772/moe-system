package usecase

import (
	"context"

	"moe-system/ent"
	"moe-system/internal/repository"
)

type HighSchoolResultUseCase interface {
	AddResult(ctx context.Context, seatNumber, studentName, section, governorate, academicYear string, totalScore, percentage float64, gradesJSON map[string]interface{}) (*ent.HighSchoolResult, error)
	GetResultBySeat(ctx context.Context, seatNumber string) (*ent.HighSchoolResult, error)
	GetAllResults(ctx context.Context) ([]*ent.HighSchoolResult, error)
}

type highSchoolResultUseCase struct {
	resultRepo repository.HighSchoolResultRepository
}

func NewHighSchoolResultUseCase(resultRepo repository.HighSchoolResultRepository) HighSchoolResultUseCase {
	return &highSchoolResultUseCase{resultRepo: resultRepo}
}

func (uc *highSchoolResultUseCase) AddResult(ctx context.Context, seatNumber, studentName, section, governorate, academicYear string, totalScore, percentage float64, gradesJSON map[string]interface{}) (*ent.HighSchoolResult, error) {
	return uc.resultRepo.CreateResult(ctx, seatNumber, studentName, section, governorate, academicYear, totalScore, percentage, gradesJSON)
}

func (uc *highSchoolResultUseCase) GetResultBySeat(ctx context.Context, seatNumber string) (*ent.HighSchoolResult, error) {
	return uc.resultRepo.GetResultBySeatNumber(ctx, seatNumber)
}

func (uc *highSchoolResultUseCase) GetAllResults(ctx context.Context) ([]*ent.HighSchoolResult, error) {
	return uc.resultRepo.ListResults(ctx)
}