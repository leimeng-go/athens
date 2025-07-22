package admin

import (
	"github.com/gorilla/mux"

)

func RegisterHandlers(r *mux.Router) {
    r.Handle(Admin,AdminHandler())
}