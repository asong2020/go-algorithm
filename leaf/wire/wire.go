//+build wireinject

package wire

import (
	"github.com/google/wire"

	"asong.cloud/go-algorithm/leaf/common"
	"asong.cloud/go-algorithm/leaf/leaf/config"
	"asong.cloud/go-algorithm/leaf/leaf/dao"
	"asong.cloud/go-algorithm/leaf/leaf/handler"
	"asong.cloud/go-algorithm/leaf/leaf/model"
	"asong.cloud/go-algorithm/leaf/service"
)

func InitializeHandler(conf *config.Server) *handler.LeafHandler {
	wire.Build(
		common.NewMysqlClient,
		common.NewRouterClient,

		dao.NewLeafDB,
		dao.NewLeafDAO,

		service.NewLeafService,
		model.NewLeafSeq,

		handler.NewLeafHandler,
	)
	return &handler.LeafHandler{}
}
