package mongo

import (
	"context"
	"fmt"
	"testing"
	"time"
	"github.com/leimeng-go/athens/pkg/config"
)

func TestDashboard(t *testing.T) {
	// 从环境变量获取MongoDB连接URL
	mongoURL := "mongodb://admin:password123@mongo-rs-1:27017,mongo-rs-2:27017,mongo-rs-3:27017/?authSource=admin&replicaSet=rs0"
	t.Logf("mongoURL: %s", mongoURL)

	ctx := context.Background()
	stat, err := NewStat(&config.MongoConfig{
		URL: mongoURL,
	}, 30*time.Second)
	if err != nil {
		t.Skip(fmt.Sprintf("无法连接到MongoDB: %v，跳过测试", err))
	}

	result, err := stat.Dashboard(ctx)
	if err != nil {
		t.Skip(fmt.Sprintf("无法获取仪表板统计信息: %v，跳过测试", err))
	}

	if result.ModuleTotal == 0 {
		t.Skip("模块总数为0，跳过测试")
	}

	t.Logf("%+v", *result)
}
