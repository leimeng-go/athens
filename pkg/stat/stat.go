package stat

import "context"

type Stat interface{
	Dashboard(ctx context.Context)(*DashboardResp,error)
	List(ctx context.Context, req *ModuleListReq) (*ModuleListResp, error)	
}

