package api

import (
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/kevincobain2000/gol"
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

	g := gol.NewGol(func(o *gol.GolOptions) error {
		o.FilePaths = []string{programConfig.LogFolder + "/" + "*.log"}
		return nil
	})

	router.HandleFunc("/gol/api", g.Adapter(g.NewAPIHandler().Get))
	router.HandleFunc("/gol", g.Adapter(g.NewAssetsHandler().Get))

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

	huma.Post(api, "/api/jobs/{id}/actions/update-limit", HandleUpdateJobLimit(queries, wakeJobs))

	huma.Get(api, "/api/jobs", HandleListJobs(queries))
	huma.Get(api, "/api/jobs/{id}", HandleJobDetails(queries))
	huma.Get(api, "/api/jobs/{id}/videos", HandleJobVideosFound(queries))
	huma.Get(api, "/api/images", HandleImages(queries))

	// TODO separate endpoints to files possibly packages
	// TODO migrate legacy routes
	router.Post("/api/references", HandleReferenceUpload(queries, config))
	router.Get("/api/references", HandleGetReferences(queries))
	router.Delete("/api/references", HandleDeleteAllReferences(queries))

	router.Get("/api/filters", HandleGetFilters(queries))

	router.Get("/api/files/{id}", HandleFileServeById(queries))
	router.Get("/api/zipped", ExportWorkspace(config))

	router.NotFound(HandleCatchAll())

	return router
}
