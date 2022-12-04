package model

type Leaf struct {
	ID          uint64 `json:"id" form:"id"`                   // 主键id
	BizTag      string `json:"biz_tag" form:"biz_tag"`         // 区分业务
	MaxID       uint64 `json:"max_id" form:"max_id"`           // 该biz_tag目前所被分配的ID号段的最大值
	Step        int32  `json:"step" form:"step"`               // 每次分配ID号段长度
	Description string `json:"description" form:"description"` // 描述
	UpdateTime  uint64 `json:"update_time" form:"update_time"` // 更新时间
}
