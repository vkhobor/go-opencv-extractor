package api

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/vkhobor/go-opencv/db_sql"
	"github.com/vkhobor/go-opencv/jobs"
)

func NewRouter(
	queries *db_sql.Queries,
	jobCreator *jobs.JobCreator,
) chi.Router {
	router := chi.NewRouter()

	// Routes
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	router.Post("/jobs", HandleCreateJob(queries, jobCreator))
	router.Get("/jobs", HandleListJobs(queries))
	router.NotFound(HandleNotFound())

	return router
}
