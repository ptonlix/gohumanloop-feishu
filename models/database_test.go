package models

import (
	"testing"

	"github.com/beego/beego/orm"
	_ "github.com/mattn/go-sqlite3"
)

func TestHumanLoop(t *testing.T) {
	// 初始化ORM
	orm.RegisterModel(new(HumanLoop))
	orm.RegisterDataBase("default", "sqlite3", "./test.db")

	// automatically build table
	orm.RunSyncdb("default", false, true)

	// 创建测试数据
	// 创建记录
	loop1 := &HumanLoop{
		TaskId:         "task_123",
		ConversationId: "conv_456",
		RequestId:      "req_789",
		LoopType:       "review",
		Context:        `{"doc_id": "doc1"}`,
		Platform:       "wechat",
	}

	o := orm.NewOrm()
	id, err := o.Insert(loop1)
	if err != nil {
		t.Fatalf("Failed to insert record: %v", err)
	}
	t.Logf("Inserted record with ID: %d", id)

}
