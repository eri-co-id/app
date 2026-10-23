package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/eri-co-id/app/internal/auth"
	"github.com/eri-co-id/app/internal/dashboard"
	"github.com/eri-co-id/app/internal/db"
	"github.com/eri-co-id/app/internal/equipment"
	"github.com/eri-co-id/app/internal/file"
	"github.com/eri-co-id/app/internal/horse"
	"github.com/eri-co-id/app/internal/lessoncode"
	"github.com/eri-co-id/app/internal/mailer"
	"github.com/eri-co-id/app/internal/record"
	"github.com/eri-co-id/app/internal/schema"
	"github.com/eri-co-id/app/internal/storage"
	"github.com/eri-co-id/app/internal/user"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Database connection
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		databaseURL = "postgres://ericoid_user:ericoid_password@localhost:5433/ericoid_db?sslmode=disable"
	}

	database, err := db.Connect(ctx, databaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	accessSecret := os.Getenv("JWT_SECRET")
	refreshSecret := os.Getenv("REFRESH_TOKEN_SECRET")
	if len(accessSecret) < 32 || len(refreshSecret) < 32 {
		log.Fatal("JWT_SECRET and REFRESH_TOKEN_SECRET must be set to at least 32 characters")
	}

	smtpConfig := mailer.SMTPConfig{
		Host:     os.Getenv("SMTP_HOST"),
		Port:     os.Getenv("SMTP_PORT"),
		User:     os.Getenv("SMTP_USER"),
		Password: os.Getenv("SMTP_PASSWORD"),
		From:     os.Getenv("SMTP_FROM"),
	}
	var appMailer interface {
		auth.Mailer
		user.Mailer
	}
	if smtpConfig.Configured() {
		appMailer = mailer.NewSMTPMailer(smtpConfig)
	} else {
		log.Println("WARNING: SMTP not configured — OTP codes and invite emails will be logged, not sent")
		appMailer = mailer.LogMailer{}
	}

	// Initialize services
	authRepo := auth.NewRepository(database)
	authSvc := auth.NewService(authRepo, accessSecret, refreshSecret, appMailer)
	authHandler := auth.NewHandler(authSvc)

	horseRepo := horse.NewRepository(database)
	horseSvc := horse.NewService(horseRepo)
	horseHandler := horse.NewHandler(horseSvc)

	equipmentRepo := equipment.NewRepository(database)
	equipmentSvc := equipment.NewService(equipmentRepo)
	equipmentHandler := equipment.NewHandler(equipmentSvc)

	lessonCodeRepo := lessoncode.NewRepository(database)
	lessonCodeSvc := lessoncode.NewService(lessonCodeRepo)
	lessonCodeHandler := lessoncode.NewHandler(lessonCodeSvc)

	userRepo := user.NewRepository(database)
	if err := user.SeedSuperAdmin(ctx, userRepo); err != nil {
		log.Fatalf("seed super admin: %v", err)
	}
	userSvc := user.NewService(userRepo, appMailer)
	userHandler := user.NewHandler(userSvc)

	schemaRepo := schema.NewRepository(database)
	schemaSvc := schema.NewService(schemaRepo)
	schemaHandler := schema.NewHandler(schemaSvc)

	recordRepo := record.NewRepository(database)
	recordSvc := record.NewService(recordRepo, schemaSvc)
	recordHandler := record.NewHandler(recordSvc)

	uploadDir := os.Getenv("UPLOAD_DIR")
	if uploadDir == "" {
		uploadDir = "./uploads"
	}
	// ponytail: local disk storage, no Garage/S3 credentials configured yet.
	// Swap for an S3-compatible client against Garage when those exist.
	fileStorage := storage.NewLocalStorage(uploadDir, "/uploads")
	fileRepo := file.NewRepository(database)
	fileSvc := file.NewService(fileRepo, recordSvc, schemaSvc, fileStorage)
	fileHandler := file.NewHandler(fileSvc)

	dashboardRepo := dashboard.NewRepository(database)
	dashboardSvc := dashboard.NewService(dashboardRepo, schemaSvc)
	dashboardHandler := dashboard.NewHandler(dashboardSvc)

	adminOnly := auth.RequireRole("admin", "super_admin")
	superAdminOnly := auth.RequireRole("super_admin")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Routes
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"status":"ok","database":"connected"}`)
	})

	http.HandleFunc("/api/auth/send-otp", authHandler.SendOTP)
	http.HandleFunc("/api/auth/verify-otp", authHandler.VerifyOTP)
	http.HandleFunc("/api/auth/refresh", authHandler.RefreshToken)
	http.HandleFunc("/api/auth/logout", authHandler.RequireAuth(authHandler.Logout))
	http.HandleFunc("/api/auth/me", authHandler.RequireAuth(authHandler.GetMe))

	http.HandleFunc("GET /api/horses", authHandler.RequireAuth(horseHandler.List))
	http.HandleFunc("POST /api/horses", authHandler.RequireAuth(adminOnly(horseHandler.Create)))
	http.HandleFunc("GET /api/horses/{id}", authHandler.RequireAuth(horseHandler.Get))
	http.HandleFunc("PATCH /api/horses/{id}", authHandler.RequireAuth(adminOnly(horseHandler.Update)))
	http.HandleFunc("DELETE /api/horses/{id}", authHandler.RequireAuth(adminOnly(horseHandler.Delete)))

	http.HandleFunc("GET /api/equipment", authHandler.RequireAuth(equipmentHandler.List))
	http.HandleFunc("POST /api/equipment", authHandler.RequireAuth(adminOnly(equipmentHandler.Create)))
	http.HandleFunc("GET /api/equipment/horse/{horseId}", authHandler.RequireAuth(equipmentHandler.ListByHorse))
	http.HandleFunc("GET /api/equipment/{id}", authHandler.RequireAuth(equipmentHandler.Get))
	http.HandleFunc("PATCH /api/equipment/{id}", authHandler.RequireAuth(adminOnly(equipmentHandler.Update)))
	http.HandleFunc("DELETE /api/equipment/{id}", authHandler.RequireAuth(adminOnly(equipmentHandler.Delete)))

	http.HandleFunc("GET /api/lesson-codes", authHandler.RequireAuth(lessonCodeHandler.List))
	http.HandleFunc("POST /api/lesson-codes", authHandler.RequireAuth(adminOnly(lessonCodeHandler.Create)))
	http.HandleFunc("GET /api/lesson-codes/{id}", authHandler.RequireAuth(lessonCodeHandler.Get))
	http.HandleFunc("PATCH /api/lesson-codes/{id}", authHandler.RequireAuth(adminOnly(lessonCodeHandler.Update)))
	http.HandleFunc("DELETE /api/lesson-codes/{id}", authHandler.RequireAuth(adminOnly(lessonCodeHandler.Delete)))

	http.HandleFunc("GET /api/users", authHandler.RequireAuth(superAdminOnly(userHandler.List)))
	http.HandleFunc("POST /api/users", authHandler.RequireAuth(superAdminOnly(userHandler.Create)))
	http.HandleFunc("GET /api/users/{id}", authHandler.RequireAuth(superAdminOnly(userHandler.Get)))
	http.HandleFunc("PATCH /api/users/{id}", authHandler.RequireAuth(superAdminOnly(userHandler.Update)))
	http.HandleFunc("PATCH /api/users/{id}/role", authHandler.RequireAuth(superAdminOnly(userHandler.UpdateRole)))
	http.HandleFunc("DELETE /api/users/{id}", authHandler.RequireAuth(superAdminOnly(userHandler.Delete)))

	http.HandleFunc("GET /api/objects", authHandler.RequireAuth(schemaHandler.List))
	http.HandleFunc("GET /api/objects/{name}", authHandler.RequireAuth(schemaHandler.Get))

	http.HandleFunc("GET /api/objects/{name}/records", authHandler.RequireAuth(recordHandler.List))
	http.HandleFunc("POST /api/objects/{name}/records", authHandler.RequireAuth(recordHandler.Create))
	http.HandleFunc("GET /api/objects/{name}/records/{id}", authHandler.RequireAuth(recordHandler.Get))
	http.HandleFunc("PATCH /api/objects/{name}/records/{id}", authHandler.RequireAuth(adminOnly(recordHandler.Update)))
	http.HandleFunc("DELETE /api/objects/{name}/records/{id}", authHandler.RequireAuth(adminOnly(recordHandler.Delete)))

	http.HandleFunc("POST /api/objects/{name}/records/{id}/files", authHandler.RequireAuth(fileHandler.Upload))
	http.HandleFunc("GET /api/objects/{name}/records/{id}/files", authHandler.RequireAuth(fileHandler.List))
	http.HandleFunc("DELETE /api/files/{id}", authHandler.RequireAuth(adminOnly(fileHandler.Delete)))
	// ponytail: requires auth but not per-record ownership (any signed-in user
	// can fetch any uploaded file by URL). Real per-file authorization comes
	// with Garage/imgproxy signed URLs (SPEC.md §6); add it then, not here.
	uploadsHandler := http.StripPrefix("/uploads/", http.FileServer(http.Dir(uploadDir)))
	http.HandleFunc("/uploads/", authHandler.RequireAuth(uploadsHandler.ServeHTTP))

	http.HandleFunc("GET /api/dashboards", authHandler.RequireAuth(dashboardHandler.List))
	http.HandleFunc("GET /api/dashboards/{name}", authHandler.RequireAuth(dashboardHandler.Get))
	http.HandleFunc("GET /api/dashboards/{name}/widgets/{widgetId}", authHandler.RequireAuth(dashboardHandler.ExecuteWidget))

	addr := fmt.Sprintf(":%s", port)
	log.Printf("Server starting on %s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
