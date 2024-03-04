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

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	router.Post("/jobs", HandleCreateJob(queries))
	router.Get("/jobs", HandleListJobs(queries))

	router.Post("/references", HandleReferenceUpload(queries))
	router.Get("/references", HandleGetReferences(queries))
	router.Delete("/references", HandleDeleteAllReferences(queries))

	router.Get("/files/{fileID}", HandleFileServeById(queries))
	router.Get("/zipped", ExportWorkspace())

	router.Get("/state", HandleAppState(queries))

	router.NotFound(HandleNotFound())

	return router
}
