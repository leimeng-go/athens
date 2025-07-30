package actions

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/leimeng-go/athens/pkg/index"
	"github.com/leimeng-go/athens/pkg/log"
	"github.com/sirupsen/logrus"
)

func statHandler(index index.Indexer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), log.CtxKey, log.NewEntry(logrus.StandardLogger()))
		total, err := index.Total(ctx)
		if err != nil {
			log.EntryFromContext(ctx).SystemErr(err)
			http.Error(w, err.Error(), errors.Kind(err))
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		enc := json.NewEncoder(w)
		if err = enc.Encode(total); err != nil {
			log.EntryFromContext(ctx).SystemErr(err)
			fmt.Fprintln(w, err)
			return
		}
	}
}