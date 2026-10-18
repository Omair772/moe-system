package http

import (
	"net/http"
	"strconv"
	"time"

	"moe-system/internal/usecase"

	"github.com/gin-gonic/gin"
)

type StudentHandler struct {
	studentUC usecase.StudentUsecase
}

func NewStudentHandler(r *gin.Engine, studentUC usecase.StudentUsecase) {
	handler := &StudentHandler{studentUC: studentUC}

	api := r.Group("/api/v1")
	{
		api.POST("/students", handler.CreateStudent)
		api.GET("/students/:id", handler.GetStudent)
		api.GET("/students", handler.GetStudents)
	}
}

func (h *StudentHandler) CreateStudent(c *gin.Context) {
	var input struct {
		StudentID      string    `json:"student_id" binding:"required"`
		FirstName      string    `json:"first_name" binding:"required"`
		LastName       string    `json:"last_name" binding:"required"`
		Email          string    `json:"email" binding:"required,email"`
		Phone          string    `json:"phone"`
		GradeLevel     int       `json:"grade_level"`
		EnrollmentDate time.Time `json:"enrollment_date"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if input.EnrollmentDate.IsZero() {
		input.EnrollmentDate = time.Now()
	}

	student, err := h.studentUC.RegisterStudent(
		c.Request.Context(),
		input.StudentID,
		input.FirstName,
		input.LastName,
		input.Email,
		input.Phone,
		input.GradeLevel,
		input.EnrollmentDate,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create student: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status": "success",
		"data":   student,
	})
}

func (h *StudentHandler) GetStudent(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid student ID"})
		return
	}

	student, err := h.studentUC.GetStudent(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Student not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   student,
	})
}

func (h *StudentHandler) GetStudents(c *gin.Context) {
	students, err := h.studentUC.GetAllStudents(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch students"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   students,
	})
}