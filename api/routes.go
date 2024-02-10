package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/vkhobor/go-opencv/db_sql"
)

func NewRouter(
	queries *db_sql.Queries,
) chi.Router {
	router := chi.NewRouter()
	router.Get("/slow", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(60 * time.Second)
		w.Write([]byte(fmt.Sprintf("all done.\n")))
	})
	router.Post("/jobs", HandleCreateJob(queries))
	return router
}
