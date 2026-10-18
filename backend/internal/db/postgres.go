package db

import (
    "context"
    "database/sql"
    "fmt"
    "log"
    "os"
    "time"

    _ "github.com/jackc/pgx/v5/stdlib"
)

// InitPostgres initializes the PostgreSQL connection pool using pgx
func InitPostgres() (*sql.DB, error) {
    dsn := os.Getenv("DATABASE_URL")
    if dsn == "" {
        // قيمة افتراضية في حال لم يتم ضبط ملف البيئة بعد
        dsn = "postgres://postgres:password@localhost:5432/moe_db?sslmode=disable"
    }

    db, err := sql.Open("pgx", dsn)
    if err != nil {
        return nil, fmt.Errorf("failed to open database: %w", err)
    }

    // ضبط إعدادات تجمع الاتصالات (Connection Pool) لتحقيق أعلى أداء
    db.SetMaxOpenConns(100)
    db.SetMaxIdleConns(50)
    db.SetConnMaxLifetime(30 * time.Minute)
    db.SetConnMaxIdleTime(10 * time.Minute)

    // التحقق من الاتصال بقاعدة البيانات
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    if err := db.PingContext(ctx); err != nil {
        return nil, fmt.Errorf("failed to ping database: %w", err)
    }

    log.Println("Database connection pool established successfully via pgx")
    return db, nil
}