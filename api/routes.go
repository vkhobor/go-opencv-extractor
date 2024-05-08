package api

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/vkhobor/go-opencv/db_sql"
)

func NewRouter(
	queries *db_sql.Queries,
	wakeJobs chan<- struct{},
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

	router.Route("/api", func(r chi.Router) {
		r.Post("/jobs", HandleCreateJob(queries, wakeJobs))
		r.Get("/jobs", HandleListJobs(queries))
		r.Get("/jobs/{id}", HandleJobDetails(queries))
		r.Get("/jobs/{id}/videos", HandleJobVideosFound(queries))
		r.Get("/jobs/{id}/progress", HandleJobProgress(queries))

		r.Post("/references", HandleReferenceUpload(queries))
		r.Get("/references", HandleGetReferences(queries))
		r.Delete("/references", HandleDeleteAllReferences(queries))

		r.Get("/images", HandleImages(queries))

		r.Get("/files/{id}", HandleFileServeById(queries))
		r.Get("/zipped", ExportWorkspace())

		r.Get("/state", HandleAppState(queries))

		r.Get("/stats", HandleGetStats(queries))
	})

	router.NotFound(HandleNotFound())

	return router
}
