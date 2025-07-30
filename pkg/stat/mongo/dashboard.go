package mongo

import (
	"context"
	"fmt"

	"github.com/leimeng-go/athens/pkg/errors"
	"github.com/leimeng-go/athens/pkg/observ"
	"go.mongodb.org/mongo-driver/bson"
	"github.com/leimeng-go/athens/pkg/stat"
)

func (s *ModuleStat)Dashboard(ctx context.Context)(*stat.DashboardResp,error){
	const op errors.Op = "mongo.Total"
	ctx,span:=observ.StartSpan(ctx,op.String())
	defer span.End()
    
    tctx,cancel:=context.WithTimeout(ctx,s.timeout)
    defer cancel()
    result:=&stat.DashboardResp{}
	total,err:=s.client.Database(s.db).Collection(s.coll).CountDocuments(tctx,bson.M{})
	if err!= nil{
		return nil,errors.E(op,err)
	}
	result.ModuleTotal=int(total)
    // 获取整个数据库的存储大小
	var dbStats bson.M
	err = s.client.Database(s.db).RunCommand(tctx, bson.M{"dbStats": 1}).Decode(&dbStats)
	if err != nil {
	   return nil, errors.E(op, err)
	}
	// 将字节大小转换为GB并转为字符串格式
	var storageSize float64
	if totalSize, ok := dbStats["totalSize"]; ok {
		switch v := totalSize.(type) {
		case int64:
			storageSize = float64(v) / (1024 * 1024 * 1024) // 转换为GB
		case int32:
			storageSize = float64(v) / (1024 * 1024 * 1024)
		case float64:
			storageSize = v / (1024 * 1024 * 1024)
		default:
			return nil, errors.E(op, "unexpected type for totalSize field")
		}
		result.StorageSize = fmt.Sprintf("%.2fGB", storageSize) // 格式化为带两位小数的GB字符串
	} else {
		return nil, errors.E(op, "totalSize field not found in dbStats result")
	}
	
	return result,nil 
}
