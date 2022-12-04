package dao

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"asong.cloud/go-algorithm/leaf/config"
	"asong.cloud/go-algorithm/leaf/model"
)

type LeafDBTest struct {
	suite.Suite
	dao *LeafDB
}

func Test_LeafDBTest(t *testing.T) {
	suite.Run(t, new(LeafDBTest))
}

func (l *LeafDBTest) SetupTest() {
	conf := &config.Server{}
	err := conf.Load("../conf/config.yaml")
	if err != nil {
		log.Panic("load conf file failed", err)
	}

	connInfo := fmt.Sprintf("%s:%s@(%s)/%s?charset=utf8&parseTime=True&loc=Local", conf.Mysql.Username, conf.Mysql.Password, conf.Mysql.Host, conf.Mysql.Db)
	db, err := sql.Open("mysql", connInfo)
	err = db.Ping()
	assert.Equal(l.T(), nil, err)
	l.dao = NewLeafDB(db)
}

func (l *LeafDBTest) Test_Create() {
	ctx := context.Background()
	leaf := &model.Leaf{
		BizTag:      "asong-leaf-segment-test",
		MaxID:       1,
		Step:        2000,
		Description: "this is a test for leaf segment",
		UpdateTime:  uint64(time.Now().Unix()),
	}
	err := l.dao.Create(ctx, leaf)
	assert.Equal(l.T(), nil, err)
}

func (l *LeafDBTest) Test_Get() {
	ctx := context.Background()
	bizTag := "asong-leaf-segment-test"
	leaf, err := l.dao.Get(ctx, bizTag, nil)
	assert.Equal(l.T(), nil, err)
	l.T().Log(leaf)
}

func (l *LeafDBTest) Test_UpdateMaxID() {
	ctx := context.Background()
	bizTag := "asong-leaf-segment-test"
	err := l.dao.UpdateMaxID(ctx, bizTag, nil)
	assert.Equal(l.T(), nil, err)
}

func (l *LeafDBTest) Test_UpdateMaxIDByCustomStep() {
	ctx := context.Background()
	bizTag := "asong-leaf-segment-test"
	step := int32(4000)
	err := l.dao.UpdateMaxIdByCustomStep(ctx, step, bizTag, nil)
	assert.Equal(l.T(), nil, err)
}

func (l *LeafDBTest) Test_GetAll() {
	list, err := l.dao.GetAll(context.Background())
	assert.Equal(l.T(), nil, err)
	l.T().Log(list[0])
}
