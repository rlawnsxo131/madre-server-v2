package middleware

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	chi_middleware "github.com/go-chi/chi/v5/middleware"
	"github.com/rlawnsxo131/madre-server-v2/lib/logger"
	"github.com/rlawnsxo131/madre-server-v2/lib/response"
)

func HTTPLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ww := chi_middleware.NewWrapResponseWriter(w, r.ProtoMajor)
		hl := logger.NewHTTPLogger(r, ww)
		err := hl.ReadBody()
		if err != nil {
			d, _ := json.Marshal(map[string]interface{}{
				"status":  http.StatusInternalServerError,
				"message": response.Http_Msg_InternalServerError,
			})
			ww.WriteHeader(http.StatusInternalServerError)
			ww.Write(d)
			return
		}

		t := time.Now()
		defer func() {
			log.Println("defer")
			hl.Write(t)
		}()

		next.ServeHTTP(
			ww,
			r.WithContext(
				logger.SetHTTPLoggerCtx(
					r.Context(),
					hl,
				),
			),
		)
	})
}

// func RequestWithHTTPLogger(r *http.Request, hl logger.HTTPLogger) *http.Request {
// 	r = r.WithContext(r, hl))
// 	return r
// }
