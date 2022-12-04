package handler

import (
	"asong.cloud/go-algorithm/leaf/model"
)

type InitLeafReq struct {
	BizTag string `json:"biz_tag"`
}

type CreateLeafReq struct {
	BizTag      string  `json:"biz_tag"`
	MaxID       *uint64 `json:"max_id"` // 可以不传 默认为1
	Step        *int32  `json:"step"`   // 可以不传 默认为2000
	Description string  `json:"description"`
}

func (c CreateLeafReq) toCreate() *model.Leaf {
	leaf := model.Leaf{}
	if c.MaxID == nil {
		leaf.MaxID = 0
	} else {
		leaf.MaxID = *c.MaxID
	}
	if c.Step == nil {
		leaf.Step = 0
	} else {
		leaf.Step = 0
	}
	leaf.BizTag = c.BizTag
	leaf.Description = c.BizTag
	return &leaf
}

type UpdateStepReq struct {
	Step   int32  `json:"step"`
	BizTag string `json:"biz_tag"`
}
