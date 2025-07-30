package stat

type ModuleRecord struct{
	Path string `json:"path"`
	Versions []string `json:"versions"`
}

type ModuleListReq struct{
	Key string `json:"key"`
	PageNum int `json:"page_num"`
	PageSize int `json:"page_size"`
}

type ModuleListResp struct{
	Total int `json:"total"`
	List []*ModuleRecord `json:"list"`
}

type DashboardResp struct{
	ModuleTotal int `json:"module_total"` //模块总数
    DownLoadTotal int `json:"download_total"` //下载总次数
	StorageSize string `json:"storage_size"` //存储大小
}