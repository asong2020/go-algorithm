package model

import (
	"fmt"
	"sync"
	"time"
)

const ExpiredTime = time.Minute * 15 //清理超过15min没更新的缓存

// 全局分配器
// key: biz_tag value: SegmentBuffer
type LeafSeq struct {
	cache sync.Map
}

func NewLeafSeq() *LeafSeq {
	seq := &LeafSeq{}
	go seq.clear()
	return seq
}

// 获取
func (l *LeafSeq) Get(bizTag string) *LeafAlloc {
	if seq, ok := l.cache.Load(bizTag); ok {
		return seq.(*LeafAlloc)
	}
	return nil
}

// 添加
func (l *LeafSeq) Add(seq *LeafAlloc) string {
	l.cache.Store(seq.Key, seq)
	return seq.Key
}

// 更新
func (l *LeafSeq) Update(key string, bean *LeafAlloc) {
	if element, ok := l.cache.Load(key); ok {
		alloc := element.(*LeafAlloc)
		alloc.Buffer = bean.Buffer
		alloc.UpdateTime = bean.UpdateTime
	}
}

// 清理超过15min没用过的内存
func (l *LeafSeq) clear() {
	for {
		now := time.Now()
		// 15分钟后
		mm, _ := time.ParseDuration("15m")
		next := now.Add(mm)
		next = time.Date(next.Year(), next.Month(), next.Day(), next.Hour(), next.Minute(), 0, 0, next.Location())
		t := time.NewTimer(next.Sub(now))
		<-t.C
		fmt.Println("start clear goroutine")
		l.cache.Range(func(key, value interface{}) bool {
			alloc := value.(*LeafAlloc)
			if next.Sub(alloc.UpdateTime) > ExpiredTime {
				fmt.Printf("clear biz_tag: %s cache", key)
				l.cache.Delete(key)
			}
			return true
		})
	}
}
