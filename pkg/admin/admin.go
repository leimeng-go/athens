package admin

import(
	"net/http"

)

const Admin="/admin"

func AdminHandler() http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("admin"))
    })
}