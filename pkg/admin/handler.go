package admin

import (
	"github.com/gorilla/mux"

)

func RegisterHandlers(r *mux.Router) {
    // 主管理页面（原有的）
    r.Handle(Admin, AdminHandler())
    
    // 仪表盘页面
    r.Handle(Admin+"/dashboard", DashboardHandler())
    
    // 仓库管理页面
    r.Handle(Admin+"/repositories", RepositoriesHandler())
    
    // 上传模块页面
    r.Handle(Admin+"/upload", UploadHandler())
    
    // 下载模块页面
    r.Handle(Admin+"/download", DownloadHandler())
}