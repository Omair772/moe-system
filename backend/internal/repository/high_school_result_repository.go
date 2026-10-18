package repository

import (
	"context"

	"moe-system/ent"
	"moe-system/ent/highschoolresult"
)

type HighSchoolResultRepository interface {
	CreateResult(ctx context.Context, seatNumber, studentName, section, governorate, academicYear string, totalScore, percentage float64, gradesJSON map[string]interface{}) (*ent.HighSchoolResult, error)
	GetResultBySeatNumber(ctx context.Context, seatNumber string) (*ent.HighSchoolResult, error)
	ListResults(ctx context.Context) ([]*ent.HighSchoolResult, error)
}

type highSchoolResultRepository struct {
	client *ent.Client
}

func NewHighSchoolResultRepository(client *ent.Client) HighSchoolResultRepository {
	return &highSchoolResultRepository{client: client}
}

func (r *highSchoolResultRepository) CreateResult(ctx context.Context, seatNumber, studentName, section, governorate, academicYear string, totalScore, percentage float64, gradesJSON map[string]interface{}) (*ent.HighSchoolResult, error) {
	return r.client.HighSchoolResult.
		Create().
		SetSeatNumber(seatNumber).
		SetStudentName(studentName).
		SetSection(section).
		SetGovernorate(governorate).
		SetAcademicYear(academicYear).
		SetTotalScore(totalScore).
		SetPercentage(percentage).
		SetGradesJSON(gradesJSON). // تم تصحيحها إلى الأحرف الكبيرة JSON
		Save(ctx)
}

func (r *highSchoolResultRepository) GetResultBySeatNumber(ctx context.Context, seatNumber string) (*ent.HighSchoolResult, error) {
	return r.client.HighSchoolResult.
		Query().
		Where(highschoolresult.SeatNumberEQ(seatNumber)).
		Only(ctx)
}

func (r *highSchoolResultRepository) ListResults(ctx context.Context) ([]*ent.HighSchoolResult, error) {
	return r.client.HighSchoolResult.Query().All(ctx)
}