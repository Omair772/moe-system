package handler

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "moe-system/internal/usecase"
)

type StudentHandler struct {
    studentUsecase usecase.StudentUsecase
}

func NewStudentHandler(su usecase.StudentUsecase) *StudentHandler {
    return &StudentHandler{studentUsecase: su}
}

type CreateStudentReq struct {
    Name    string `json:"name" binding:"required"`
    Email   string `json:"email" binding:"required"`
    Grade   string `json:"grade"`
    Section string `json:"section"`
}

func (h *StudentHandler) CreateStudent(c *gin.Context) {
    var req CreateStudentReq
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "بيانات غير صالحة"})
        return
    }

    student, err := h.studentUsecase.CreateStudent(c.Request.Context(), req.Name, req.Email, req.Grade, req.Section)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusCreated, gin.H{
        "message": "تم إضافة الطالب بنجاح",
        "data":    student,
    })
}

func (h *StudentHandler) ListStudents(c *gin.Context) {
    students, err := h.studentUsecase.ListStudents(c.Request.Context())
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "data": students,
    })
}
