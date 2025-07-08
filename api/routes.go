package api

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/kevincobain2000/gol"
	"github.com/vkhobor/go-opencv/api/filters"
	"github.com/vkhobor/go-opencv/api/jobs"
	"github.com/vkhobor/go-opencv/api/testsurf"
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
		Tags:          []string{"Videos"},
		Path:          "/api/videos",
		DefaultStatus: 200,
		Summary:       "List downloaded videos",
	}, HandleListVideos(queries))

	huma.Register(api, huma.Operation{
		Method:        "POST",
		Tags:          []string{"Jobs"},
		Path:          "/api/jobs/video",
		DefaultStatus: 201,
		Summary:       "Create a direct video job",
	}, jobs.HandleImportJob(queries, programConfig, wakeJobs))

	huma.Register(api, huma.Operation{
		Method:        "POST",
		Tags:          []string{"References"},
		Path:          "/api/references",
		DefaultStatus: 201,
		Summary:       "Upload reference images",
	}, filters.HandleReferenceUpload(queries, config))

	huma.Register(api, huma.Operation{
		Method:        "POST",
		Tags:          []string{"TestSurf"},
		Path:          "/testsurf/video",
		DefaultStatus: 201,
		Summary:       "Upload a video file",
	}, testsurf.HandleUploadVideo(config))

	huma.Register(api, huma.Operation{
		Method:        "POST",
		Tags:          []string{"TestSurf"},
		Path:          "/testsurf/reference",
		DefaultStatus: 201,
		Summary:       "Upload a reference file",
	}, testsurf.HandleUploadReference())

	huma.Register(api, huma.Operation{
		Method:        "GET",
		Tags:          []string{"TestSurf"},
		Path:          "/testsurf/test",
		DefaultStatus: 200,
		Summary:       "Test frame matching",
	}, testsurf.HandleFrameMatchingTest())

	huma.Register(api, huma.Operation{
		Method:        "GET",
		Tags:          []string{"TestSurf"},
		Path:          "/testsurf/frame",
		DefaultStatus: 200,
		Summary:       "Retrieve a specific frame image",
	}, testsurf.HandleRetrieveFrameImage(config))

	huma.Register(api, huma.Operation{
		Method:        "GET",
		Tags:          []string{"TestSurf"},
		Path:          "/testsurf",
		DefaultStatus: 200,
		Summary:       "Get video metadata",
	}, testsurf.HandleVideoMetadata(config))

	huma.Register(api, huma.Operation{
		Method:        "GET",
		Tags:          []string{"Jobs"},
		Path:          "/api/jobs",
		DefaultStatus: 200,
		Summary:       "List all jobs",
	}, jobs.HandleListJobs(queries))

	huma.Register(api, huma.Operation{
		Method:        "GET",
		Tags:          []string{"Images"},
		Path:          "/api/images",
		DefaultStatus: 200,
		Summary:       "List images",
	}, HandleImages(queries))

	huma.Register(api, huma.Operation{
		Method:        "GET",
		Tags:          []string{"References"},
		Path:          "/api/references/{id}",
		DefaultStatus: 200,
		Summary:       "Get reference by ID",
	}, filters.HandleReferenceGet(queries))

	registerLegacyChiRoutes(router, queries, config)

	router.NotFound(HandleCatchAll())

	router.Get("/docs", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(`<!doctype html>
<html>
  <head>
    <title>API Reference</title>
    <meta charset="utf-8" />
    <meta
      name="viewport"
      content="width=device-width, initial-scale=1" />
  </head>
  <body>
    <script
      id="api-reference"
      data-url="/openapi.json"></script>
    <script src="https://cdn.jsdelivr.net/npm/@scalar/api-reference"></script>
  </body>
</html>`))
	})

	return router
}

// registerLegacyChiRoutes registers routes directly using Chi router methods
func registerLegacyChiRoutes(router chi.Router, queries *db.Queries, config config.DirectoryConfig) {
	// TODO migrate legacy routes
	router.Get("/api/references", filters.HandleGetReferences(queries))
	router.Get("/api/filters", filters.HandleGetFilters(queries))
	router.Get("/api/files/{id}", HandleFileServeById(queries)) // Files tag would be added when migrated to huma
	router.Get("/api/zipped", ExportWorkspace(config))          // Export tag would be added when migrated to huma
}
