package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"asong.cloud/go-algorithm/leaf/service"
)

type LeafHandler struct {
	engine  *gin.Engine
	service *service.LeafService
}

func NewLeafHandler(engine *gin.Engine, service *service.LeafService) *LeafHandler {
	return &LeafHandler{
		engine:  engine,
		service: service,
	}
}

func (h *LeafHandler) Run() {
	// Force log's color
	gin.ForceConsoleColor()
	h.engine.Use(gin.Logger())
	h.engine.Use(gin.Recovery())
	h.registerRouter()

	err := h.engine.Run()
	if err != nil {
		log.Fatalln("server start failed")
	}
}

func (h *LeafHandler) registerRouter() {
	r := h.engine.Group("api/leaf")
	{
		r.PUT("/init/cache", h.Init)
		r.GET("", h.Get)
		r.POST("", h.Create)
		r.PUT("/step", h.UpdateStep)
	}
}

func (h *LeafHandler) Init(ctx *gin.Context) {
	req := InitLeafReq{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code": 1000, "msg": "param Invalid"})
		return
	}
	alloc, err := h.service.InitCache(ctx, req.BizTag)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code": 1001, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "msg": "success", "data": alloc})
}

func (h *LeafHandler) Get(ctx *gin.Context) {
	bizTag := ctx.Query("biz_tag")
	id, err := h.service.GetID(ctx, bizTag)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code": 1001, "msg": err.Error()})
		return
	}
	resp := struct {
		ID uint64
	}{
		ID: id,
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "msg": "success", "data": resp})
}

func (h *LeafHandler) Create(ctx *gin.Context) {
	req := CreateLeafReq{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code": 1000, "msg": "param Invalid"})
		return
	}
	err := h.service.Create(ctx, req.toCreate())
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code": 1001, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "msg": "success", "data": nil})
}

func (h *LeafHandler) UpdateStep(ctx *gin.Context) {
	req := UpdateStepReq{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code": 1000, "msg": "param Invalid"})
		return
	}
	err := h.service.UpdateStep(ctx, req.Step, req.BizTag)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code": 1001, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "msg": "success", "data": nil})
}
