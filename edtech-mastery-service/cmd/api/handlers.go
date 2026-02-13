package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"

	"github.com/edtech-mastery/student-progress-service/internal/domain"
	"github.com/edtech-mastery/student-progress-service/internal/dashboard"
	"github.com/edtech-mastery/student-progress-service/internal/events"
	"github.com/edtech-mastery/student-progress-service/pkg/logging"
	"github.com/edtech-mastery/student-progress-service/pkg/metrics"
)

func eventsHandler(log zerolog.Logger, svc *events.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		reqID := r.Header.Get("X-Request-ID")
		if reqID == "" {
			reqID = chi.URLParam(r, "requestID")
		}
		l := logging.WithRequestID(log, reqID)

		var in domain.IncomingEvent
		if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
			l.Warn().Err(err).Msg("decode event")
			http.Error(w, `{"error":"invalid json"}`, http.StatusBadRequest)
			return
		}
		eventType, eventDBID, err := svc.Ingest(r.Context(), &in)
		if err != nil {
			l.Warn().Err(err).Msg("ingest event")
			metrics.IngestionLatency.Observe(time.Since(start).Seconds())
			http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusBadRequest)
			return
		}
		metrics.IngestionLatency.Observe(time.Since(start).Seconds())
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusAccepted)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"event_id":   in.EventID,
			"type":       eventType,
			"event_db_id": eventDBID,
		})
	}
}

func dashboardHandler(log zerolog.Logger, svc *dashboard.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		defer func() { metrics.DashboardQueryLatency.WithLabelValues("teacher_dashboard").Observe(time.Since(start).Seconds()) }()
		teacherID := chi.URLParam(r, "teacherID")
		classID := chi.URLParam(r, "classID")
		dash, err := svc.TeacherDashboard(r.Context(), teacherID, classID)
		if err != nil {
			log.Warn().Err(err).Msg("dashboard")
			http.Error(w, `{"error":"not found"}`, http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(dash)
	}
}

func masteryHandler(log zerolog.Logger, svc *dashboard.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		defer func() { metrics.DashboardQueryLatency.WithLabelValues("student_mastery").Observe(time.Since(start).Seconds()) }()
		studentID := chi.URLParam(r, "studentID")
		m, err := svc.StudentMastery(r.Context(), studentID)
		if err != nil {
			log.Warn().Err(err).Msg("mastery")
			http.Error(w, `{"error":"not found"}`, http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(m)
	}
}

func timelineHandler(log zerolog.Logger, svc *dashboard.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		defer func() { metrics.DashboardQueryLatency.WithLabelValues("timeline").Observe(time.Since(start).Seconds()) }()
		studentID := chi.URLParam(r, "studentID")
		classID := chi.URLParam(r, "classID")
		t, err := svc.StudentTimeline(r.Context(), studentID, classID, 50)
		if err != nil {
			log.Warn().Err(err).Msg("timeline")
			http.Error(w, `{"error":"not found"}`, http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(t)
	}
}
