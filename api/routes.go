package api

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/vkhobor/go-opencv/config"
	"github.com/vkhobor/go-opencv/db"
)

func NewRouter(
	queries *db.Queries,
	wakeJobs chan<- struct{},
	config config.DirectoryConfig,
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

		r.Post("/references", HandleReferenceUpload(queries, config))
		r.Get("/references", HandleGetReferences(queries))
		r.Delete("/references", HandleDeleteAllReferences(queries))

		r.Get("/filters", HandleGetFilters(queries))

		r.Get("/images", HandleImages(queries))

		r.Get("/files/{id}", HandleFileServeById(queries))
		r.Get("/zipped", ExportWorkspace(config))
	})

	router.NotFound(HandleNotFound())

	return router
}
