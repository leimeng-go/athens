package mongo

import (
	"context"
	"fmt"

	"github.com/leimeng-go/athens/pkg/errors"
	"github.com/leimeng-go/athens/pkg/observ"
	"github.com/leimeng-go/athens/pkg/stat"
	"go.mongodb.org/mongo-driver/bson"
  )

// List 根据查询条件获取模块列表
// 实现了stat.Stat接口的List方法
func (s *ModuleStat) List(ctx context.Context, req *stat.ModuleListReq) (*stat.ModuleListResp, error) {
	const op errors.Op = "mongo.List" 
	ctx, span := observ.StartSpan(ctx, op.String())
	defer span.End()

	tctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	// 计算分页偏移量
	offset := (req.PageNum - 1) * req.PageSize
	limit := req.PageSize
	// 构建聚合管道
	pipeline := []bson.M{
		{
			"$group": bson.M{
				"_id":      "$module",
				"versions": bson.M{"$push": "$version"},
			},
		},
		{
			"$project": bson.M{
				"_id":      0,
				"module":   "$_id",
				"versions": 1,
			},
		},
		{
			"$skip": offset,
		},
		{
			"$limit": limit,
		},
	}

	// 如果有搜索关键字，添加匹配条件
	if req.Key != "" {
		matchStage := bson.M{
			"$match": bson.M{
				"module": bson.M{
					"$regex":   fmt.Sprintf(".*%s.*", req.Key),
					"$options": "i",
				},
			},
		}
		pipeline = append([]bson.M{matchStage}, pipeline...)
	}

	cursor, err := s.client.Database(s.db).Collection(s.coll).Aggregate(tctx, pipeline)
	if err != nil {
		return nil, errors.E(op, err)
	}

	// 定义一个临时结构体来接收MongoDB中的数据
	type moduleDoc struct {
		Module   string   `bson:"module"`	
		Versions []string `bson:"versions"`
	}

	var docs []moduleDoc
	err = cursor.All(tctx, &docs)
	if err != nil {
		return nil, errors.E(op, err)
	}

	// 将查询结果转换为ModuleRecord格式
	moduleMap := make(map[string][]string)
	for _, doc := range docs {
		moduleMap[doc.Module] = append(moduleMap[doc.Module], doc.Versions...)
	}

	// 构建返回结果
	result := &stat.ModuleListResp{
		List: make([]*stat.ModuleRecord, 0, len(moduleMap)),
	}

	for path, versions := range moduleMap {
		result.List = append(result.List, &stat.ModuleRecord{
			Path:     path,
			Versions: versions,
		})
	}
	// 构建计算总数的聚合管道
	countPipeline := []bson.M{
		{
			"$group": bson.M{
				"_id": "$module",
			},
		},
		{
			"$count": "total",
		},
	}

	// 如果有搜索关键字，添加匹配条件
	if req.Key != "" {
		matchStage := bson.M{
			"$match": bson.M{
				"module": bson.M{
					"$regex":   fmt.Sprintf(".*%s.*", req.Key),
					"$options": "i",
				},
			},
		}
		countPipeline = append([]bson.M{matchStage}, countPipeline...)
	}

	// 执行聚合查询获取总数
	countCursor, err := s.client.Database(s.db).Collection(s.coll).Aggregate(tctx, countPipeline)
	if err != nil {
		return nil, errors.E(op, err)
	}

	type countResult struct {
		Total int64 `bson:"total"`
	}
	var counts []countResult
	if err = countCursor.All(tctx, &counts); err != nil {
		return nil, errors.E(op, err)
	}

	// 设置总数
	if len(counts) > 0 {
		result.Total = int(counts[0].Total)
	} else {
		result.Total = 0
	}

	return result, nil
}