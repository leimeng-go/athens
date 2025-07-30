package mongo

import (
	"context"
	"testing"
	"time"

	"github.com/leimeng-go/athens/pkg/config"
)

func TestDashboard(t *testing.T){
	ctx:=context.Background()
	stat,err:=NewStat(&config.MongoConfig{
		URL: "mongodb://admin:password123@mongo-rs-1:27017,mongo-rs-2:27017,mongo-rs-3:27017/?authSource=admin&replicaSet=rs0",
	},30*time.Second)
	if err!=nil{
		t.Fatal(err)
	}
	result,err:=stat.Dashboard(ctx)
	if err!= nil{
		t.Fatal(err)
	}
	if result.ModuleTotal==0{
		t.Fatal("module total is 0")
	}
}