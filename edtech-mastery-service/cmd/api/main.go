package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/edtech-mastery/student-progress-service/internal/dashboard"
	"github.com/edtech-mastery/student-progress-service/internal/events"
	"github.com/edtech-mastery/student-progress-service/internal/storage"
	"github.com/edtech-mastery/student-progress-service/pkg/logging"
)

func main() {
	env := os.Getenv("ENV")
	if env == "" {
		env = "development"
	}
	log := logging.Init(env)

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, getDSN())
	if err != nil {
		log.Fatal().Err(err).Msg("db connect")
	}
	defer pool.Close()

	eventRepo := storage.NewEventRepo(pool)
	eventsSvc := events.NewService(eventRepo)

	rollupsRepo := storage.NewRollupsRepo(pool)
	riskRepo := storage.NewRiskRepo(pool)
	recentRepo := storage.NewRecentActivityRepo(pool)
	masteryRepo := storage.NewMasteryRepo(pool)
	timelineRepo := storage.NewTimelineRepo(pool)
	dashboardSvc := dashboard.NewService(rollupsRepo, riskRepo, recentRepo, masteryRepo, timelineRepo)

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Post("/events", eventsHandler(log, eventsSvc))
	r.Get("/teachers/{teacherID}/classes/{classID}/dashboard", dashboardHandler(log, dashboardSvc))
	r.Get("/students/{studentID}/mastery", masteryHandler(log, dashboardSvc))
	r.Get("/classes/{classID}/students/{studentID}/timeline", timelineHandler(log, dashboardSvc))
	r.Handle("/metrics", promhttp.Handler())

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	srv := &http.Server{Addr: ":" + port, Handler: r}

	go func() {
		log.Info().Str("port", port).Msg("api listening")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("server")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info().Msg("shutting down")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Error().Err(err).Msg("shutdown")
	}
}

func getDSN() string {
	if s := os.Getenv("DATABASE_URL"); s != "" {
		return s
	}
	return "postgres://postgres:postgres@localhost:5432/edtech?sslmode=disable"
}
