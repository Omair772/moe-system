package main

import (
    "context"
    "database/sql"
    "log"

    "entgo.io/ent/dialect"
    sqldriver "entgo.io/ent/dialect/sql"
    "github.com/gin-gonic/gin"
    _ "modernc.org/sqlite"
    "github.com/redis/go-redis/v9"

    "moe-system/ent"
    "moe-system/internal/delivery/http"
    "moe-system/internal/repository"
    "moe-system/internal/usecase"
)

func main() {
    // فتح اتصال SQLite مباشرة وتمريره لـ Ent لتجاوز أي قيود في الأسماء
    db, err := sql.Open("sqlite", "file:moe.db?cache=shared&_fk=1")
    if err != nil {
        log.Fatalf("failed opening connection to sqlite: %v", err)
    }
    defer db.Close()

    drv := sqldriver.OpenDB(dialect.SQLite, db)
    client := ent.NewClient(ent.Driver(drv))
    defer client.Close()

    ctx := context.Background()
    if err := client.Schema.Create(ctx); err != nil {
        log.Fatalf("failed creating schema resources: %v", err)
    }
    log.Println("SQLite (Pure-Go) Database connected and schema migrated successfully!")

    rdb := redis.NewClient(&redis.Options{
        Addr:     "localhost:6379",
        Password: "",
        DB:       0,
    })
    if err := rdb.Ping(ctx).Err(); err != nil {
        log.Printf("Warning: Redis connection failed: %v", err)
    } else {
        log.Println("Redis connected successfully!")
    }

    r := gin.Default()

    r.Use(func(c *gin.Context) {
        c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
        c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
        c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
        c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(204)
            return
        }
        c.Next()
    })

    userRepo := repository.NewUserRepository(client)
    studentRepo := repository.NewStudentRepository(client)
    teacherRepo := repository.NewTeacherRepository(client)

    studentUsecase := usecase.NewStudentUsecase(studentRepo, userRepo)
    teacherUsecase := usecase.NewTeacherUsecase(teacherRepo, userRepo)

    http.NewStudentHandler(r, studentUsecase)
    http.NewTeacherHandler(r, teacherUsecase)

    log.Println("Server running on port 9000...")
    if err := r.Run(":9000"); err != nil {
        log.Fatalf("failed to run server: %v", err)
    }
}