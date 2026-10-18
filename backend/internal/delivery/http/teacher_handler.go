package http

import (
	"net/http"
	"strconv"
	"time"

	"moe-system/internal/usecase"

	"github.com/gin-gonic/gin"
)

type TeacherHandler struct {
	teacherUC usecase.TeacherUsecase
}

func NewTeacherHandler(r *gin.Engine, teacherUC usecase.TeacherUsecase) {
	handler := &TeacherHandler{teacherUC: teacherUC}

	api := r.Group("/api/v1")
	{
		api.POST("/teachers", handler.CreateTeacher)
		api.GET("/teachers/:id", handler.GetTeacher)
		api.GET("/teachers", handler.GetTeachers)
	}
}

func (h *TeacherHandler) CreateTeacher(c *gin.Context) {
	var input struct {
		TeacherID  string    `json:"teacher_id" binding:"required"`
		FirstName  string    `json:"first_name" binding:"required"`
		LastName   string    `json:"last_name" binding:"required"`
		Email      string    `json:"email" binding:"required,email"`
		Subject    string    `json:"subject"`
		Department string    `json:"department"`
		Phone      string    `json:"phone"`
		HireDate   time.Time `json:"hire_date"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if input.HireDate.IsZero() {
		input.HireDate = time.Now()
	}

	teacher, err := h.teacherUC.RegisterTeacher(
		c.Request.Context(),
		input.TeacherID,
		input.FirstName,
		input.LastName,
		input.Email,
		input.Subject,
		input.Department,
		input.Phone,
		input.HireDate,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create teacher: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status": "success",
		"data":   teacher,
	})
}

func (h *TeacherHandler) GetTeacher(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid teacher ID"})
		return
	}

	teacher, err := h.teacherUC.GetTeacher(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Teacher not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   teacher,
	})
}

func (h *TeacherHandler) GetTeachers(c *gin.Context) {
	teachers, err := h.teacherUC.GetAllTeachers(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch teachers"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   teachers,
	})
}