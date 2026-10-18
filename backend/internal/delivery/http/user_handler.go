package http

import (
    "net/http"
    "strconv"

    "moe-system/internal/usecase"

    "github.com/gin-gonic/gin"
)

type UserHandler struct {
    userUC usecase.UserUseCase
}

func NewUserHandler(r *gin.Engine, userUC usecase.UserUseCase) {
    handler := &UserHandler{userUC: userUC}

    api := r.Group("/api/v1")
    {
        api.POST("/users", handler.CreateUser)
        api.GET("/users/:id", handler.GetUser)
        api.GET("/users", handler.GetUsers)
    }
}

func (h *UserHandler) CreateUser(c *gin.Context) {
    var input struct {
        FullName string `json:"full_name" binding:"required"`
        Email    string `json:"email" binding:"required,email"`
        Password string `json:"password" binding:"required"`
        Role     string `json:"role" binding:"required"`
    }

    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    user, err := h.userUC.CreateUser(c.Request.Context(), input.FullName, input.Email, input.Password, input.Role)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusCreated, gin.H{
        "status": "success",
        "data":   user,
    })
}

func (h *UserHandler) GetUser(c *gin.Context) {
    idParam := c.Param("id")
    id, err := strconv.Atoi(idParam)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
        return
    }

    user, err := h.userUC.GetUserByID(c.Request.Context(), id)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "status": "success",
        "data":   user,
    })
}

func (h *UserHandler) GetUsers(c *gin.Context) {
    users, err := h.userUC.GetAllUsers(c.Request.Context())
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "status": "success",
        "data":   users,
    })
}
