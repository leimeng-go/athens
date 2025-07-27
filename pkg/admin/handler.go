package admin

import (
	"github.com/gorilla/mux"

)

const(
    Admin = "/admin"
)

func RegisterHandlers(r *mux.Router) {
    // 注册所有API路由
    RegisterAPIHandlers(r)
}