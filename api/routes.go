package api

import (
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/kevincobain2000/gol"
	"github.com/vkhobor/go-opencv/api/filters"
	"github.com/vkhobor/go-opencv/api/jobs"
	"github.com/vkhobor/go-opencv/config"
	"github.com/vkhobor/go-opencv/db"
)

func NewRouter(
	queries *db.Queries,
	wakeJobs chan<- struct{},
	config config.DirectoryConfig,
	programConfig config.ServerConfig,
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
		Method:        "GET",
		Path:          "/api/videos",
		DefaultStatus: 200,
		Summary:       "List downloaded videos",
	}, HandleListVideos(queries))

	huma.Register(api, huma.Operation{
		Method:        "POST",
		Path:          "/api/jobs",
		DefaultStatus: 201,
		Summary:       "Create a new job",
	}, jobs.HandleCreateJob(queries, wakeJobs, programConfig))

	huma.Register(api, huma.Operation{
		Method:        "POST",
		Path:          "/api/jobs/video",
		DefaultStatus: 201,
		Summary:       "Create a direct video job",
	}, HandleImportJob(queries, programConfig, wakeJobs))

	huma.Register(api, huma.Operation{
		Method:        "POST",
		Path:          "/api/jobs/{id}/actions/restart",
		DefaultStatus: 202,
		Summary:       "Restart the job pipeline",
	}, jobs.HandleRestartJobPipeline(wakeJobs))

	huma.Post(api, "/api/jobs/{id}/actions/update-limit", jobs.HandleUpdateJobLimit(queries, wakeJobs))

	huma.Get(api, "/api/jobs", jobs.HandleListJobs(queries))
	huma.Get(api, "/api/jobs/{id}", jobs.HandleJobDetails(queries))
	huma.Get(api, "/api/jobs/{id}/videos", jobs.HandleJobVideosFound(queries))
	huma.Get(api, "/api/images", HandleImages(queries))

	huma.Register(api, huma.Operation{
		Method:        "POST",
		Path:          "/api/references",
		DefaultStatus: 201,
		Summary:       "Upload reference images",
	}, filters.HandleReferenceUpload(queries, config))

	// TODO migrate legacy routes
	router.Get("/api/references", filters.HandleGetReferences(queries))
	router.Delete("/api/references", filters.HandleDeleteAllReferences(queries))

	router.Get("/api/filters", filters.HandleGetFilters(queries))

	router.Get("/api/files/{id}", HandleFileServeById(queries))
	router.Get("/api/zipped", ExportWorkspace(config))

	router.NotFound(HandleCatchAll())

	return router
}
