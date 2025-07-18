package actions

import (
	"net/http"

	"github.com/leimeng-go/athens/pkg/storage"
)

func getReadinessHandler(s storage.Backend) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if _, err := s.List(r.Context(), "github.com/leimeng-go/athens"); err != nil {
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}
