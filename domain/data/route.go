package data

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/rlawnsxo131/madre-server-v2/database"
	"github.com/rlawnsxo131/madre-server-v2/lib/logger"
	"github.com/rlawnsxo131/madre-server-v2/lib/response"
	"github.com/rlawnsxo131/madre-server-v2/utils"
)

func ApplyRoutes(v1 *mux.Router) {
	route := v1.NewRoute().PathPrefix("/data").Subrouter()

	route.HandleFunc("", getAll()).Methods("GET")
	route.HandleFunc("/{id}", get()).Methods("GET")
}

func getAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rw := response.NewWriter(w, r)
		limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
		if err != nil {
			logger.GetDefaultLogger().
				Warn().Msgf("route: limit Atoi wrong: %v", err)
		}
		limit = utils.IfIsNotExistGetDefaultIntValue(limit, 50)

		db, err := database.LoadFromHttpCtx(r.Context())
		if err != nil {
			rw.Error(err, "get /data")
			return
		}

		dataUseCase := NewUseCase(db)
		dd, err := dataUseCase.FindAll(limit)
		if err != nil {
			rw.Error(err, "get /data")
			return
		}

		rw.Compress(dd)
	}
}

func get() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rw := response.NewWriter(w, r)
		vars := mux.Vars(r)
		id := vars["id"]

		db, err := database.LoadFromHttpCtx(r.Context())
		if err != nil {
			rw.Error(err, "get /data/{id}")
			return
		}

		dataUseCase := NewUseCase(db)
		d, err := dataUseCase.FindOneById(id)
		if err != nil {
			rw.Error(err, "get /data/{id}")
			return
		}

		rw.Compress(d)
	}
}
