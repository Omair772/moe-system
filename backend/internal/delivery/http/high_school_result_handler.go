package http

import (
	"net/http"

	"moe-system/internal/usecase"

	"github.com/gin-gonic/gin"
)

type HighSchoolResultHandler struct {
	resultUC usecase.HighSchoolResultUseCase
}

func NewHighSchoolResultHandler(r *gin.Engine, resultUC usecase.HighSchoolResultUseCase) {
	handler := &HighSchoolResultHandler{resultUC: resultUC}

	api := r.Group("/api/v1")
	{
		api.POST("/results", handler.CreateResult)
		api.GET("/results/seat/:seat", handler.GetResultBySeat)
		api.GET("/results", handler.GetResults)
	}
}

func (h *HighSchoolResultHandler) CreateResult(c *gin.Context) {
	var input struct {
		SeatNumber   string                 `json:"seat_number" binding:"required"`
		StudentName  string                 `json:"student_name" binding:"required"`
		Section      string                 `json:"section"`
		Governorate  string                 `json:"governorate"`
		AcademicYear string                 `json:"academic_year" binding:"required"`
		TotalScore   float64                `json:"total_score"`
		Percentage   float64                `json:"percentage"`
		GradesJSON   map[string]interface{} `json:"grades_json"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.resultUC.AddResult(
		c.Request.Context(),
		input.SeatNumber,
		input.StudentName,
		input.Section,
		input.Governorate,
		input.AcademicYear,
		input.TotalScore,
		input.Percentage,
		input.GradesJSON,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create result: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status": "success",
		"data":   result,
	})
}

func (h *HighSchoolResultHandler) GetResultBySeat(c *gin.Context) {
	seatNumber := c.Param("seat")

	result, err := h.resultUC.GetResultBySeat(c.Request.Context(), seatNumber)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Result not found for this seat number"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   result,
	})
}

func (h *HighSchoolResultHandler) GetResults(c *gin.Context) {
	results, err := h.resultUC.GetAllResults(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch results"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   results,
	})
}