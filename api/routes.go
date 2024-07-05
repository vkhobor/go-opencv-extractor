package api

import (
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/vkhobor/go-opencv/config"
	"github.com/vkhobor/go-opencv/db"
)

func NewRouter(
	queries *db.Queries,
	wakeJobs chan<- struct{},
	config config.DirectoryConfig,
	programConfig config.ProgramConfig,
) chi.Router {
	router := chi.NewMux()
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	api := humachi.New(router, huma.DefaultConfig("My API", "1.0.0"))

	huma.Register(api, huma.Operation{
		Method:        "POST",
		Path:          "/api/jobs",
		DefaultStatus: 201,
		Summary:       "Create a new job",
	}, HandleCreateJob(queries, wakeJobs, programConfig))

	huma.Register(api, huma.Operation{
		Method:        "POST",
		Path:          "/api/jobs/{id}/actions/restart",
		DefaultStatus: 202,
		Summary:       "Restart the job pipeline",
	}, HandleRestartJobPipeline(wakeJobs))

	huma.Get(api, "/api/jobs", HandleListJobs(queries))
	huma.Get(api, "/api/jobs/{id}", HandleJobDetails(queries))
	huma.Get(api, "/api/jobs/{id}/videos", HandleJobVideosFound(queries))

	router.Route("/api", func(r chi.Router) {

		r.Post("/references", HandleReferenceUpload(queries, config))
		r.Get("/references", HandleGetReferences(queries))
		r.Delete("/references", HandleDeleteAllReferences(queries))

		r.Get("/filters", HandleGetFilters(queries))

		r.Get("/images", HandleImages(queries))

		r.Get("/files/{id}", HandleFileServeById(queries))
		r.Get("/zipped", ExportWorkspace(config))
	})

	router.NotFound(HandleCatchAll())

	return router
}
