package admin

import (
	"net/http"

	"github.com/leimeng-go/athens/pkg/download"
	"github.com/leimeng-go/athens/pkg/errors"
	"github.com/leimeng-go/athens/pkg/log"
)

const Dashboard = "/dashboard/stat"

func DashboardStatHandler(dp download.Protocol,lggr log.Entry) http.Handler {
	const op errors.Op = "admin.DashboardStatHandler"
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
         w.Header().Set("Content-Type", "application/json; charset=utf-8")
		//  
    })
}
