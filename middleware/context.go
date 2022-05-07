package middleware

import (
	"net/http"

	"github.com/rlawnsxo131/madre-server-v2/constants"
	"github.com/rlawnsxo131/madre-server-v2/database"
	"github.com/rlawnsxo131/madre-server-v2/lib/response"
	"github.com/rlawnsxo131/madre-server-v2/lib/syncmap"
)

func SetSyncMapCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := syncmap.GenerateHttpCtx(r.Context())
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func SetDatabaseToSyncMapCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		db, err := database.GetDatabaseInstance()
		if err != nil {
			rw := response.NewWriter(w, r)
			rw.Error(
				err,
				"SetDBContext",
			)
			return
		}

		ctx, err := syncmap.SetNewValueFromHttpCtx(
			r.Context(),
			constants.Key_HttpContextDB,
			db,
		)
		if err != nil {
			rw := response.NewWriter(w, r)
			rw.Error(
				err,
				"SetDBContext",
				"context set error",
			)
			return
		}

		r.Context().Value(ctx)
		next.ServeHTTP(w, r)
	})
}
