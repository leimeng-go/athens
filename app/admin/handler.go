package admin

import (
	"github.com/gorilla/mux"
    "github.com/leimeng-go/athens/pkg/download"

)

const(
    Admin = "/admin"
)

func RegisterHandlers(r *mux.Router,opts *download.HandlerOpts) {
    if opts==nil || opts.Protocol==nil||opts.Logger==nil{
        panic("absolutely unacceptable handler opts")
    }
    // 注册所有API路由
    RegisterAPIHandlers(r)
}

