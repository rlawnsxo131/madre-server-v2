package user

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rlawnsxo131/madre-server-v2/database"
	"github.com/rlawnsxo131/madre-server-v2/lib/response"
)

func ApplyRoutes(v1 *mux.Router) {
	route := v1.NewRoute().PathPrefix("/user").Subrouter()

	route.HandleFunc("/{id}", get()).Methods("GET")
	route.HandleFunc("/{id}", put()).Methods("PUT", "OPTIONS")
}

func get() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		writer := response.NewHttpWriter(w, r)
		vars := mux.Vars(r)
		id := vars["id"]

		db, err := database.GetDBConn(r.Context())
		if err != nil {
			writer.WriteError(err, "get /user/{id}")
			return
		}

		userService := NewUserService(db)
		user, err := userService.FindOneById(id)
		if err != nil {
			writer.WriteError(err, "get /user/{id}")
			return
		}

		writer.WriteCompress(user)
	}
}

func put() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {}
}
