package model

// 号段
type LeafSegment struct {
	Cursor uint64 // 当前发放位置
	Max    uint64 // 最大值
	Min    uint64 // 开始值即最小值
	InitOk bool   // 是否初始化成功
}

// 以初始化DB中的MAXID为1举例子，也就是号段从1开始，步长为1000,最大值就是1000，范围是1～1000。下一段范围就是1001～2000
func NewLeafSegment(leaf *Leaf) *LeafSegment {
	return &LeafSegment{
		Cursor: leaf.MaxID - uint64(leaf.Step+1), // 最小值的前一个值
		Max:    leaf.MaxID - 1,                   // DB默认存的是1 所以这里要减1
		Min:    leaf.MaxID - uint64(leaf.Step),   // 开始的最小值
		InitOk: true,
	}
}
