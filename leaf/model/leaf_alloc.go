package model

import (
	"sync"
	"time"
)

type LeafAlloc struct {
	Key        string                 // 也就是`biz_tag`用来区分业务
	Step       int32                  // 记录步长
	CurrentPos int32                  // 当前使用的 segment buffer光标; 总共两个buffer缓存区，循环使用
	Buffer     []*LeafSegment         // 双buffer 一个作为预缓存作用
	UpdateTime time.Time              // 记录更新时间 方便长时间不用进行清理，防止占用内存
	mutex      sync.Mutex             // 互斥锁
	IsPreload  bool                   // 是否正在预加载
	Waiting    map[string][]chan byte // 挂起等待
}

func NewLeafAlloc(leaf *Leaf) *LeafAlloc {
	return &LeafAlloc{
		Key:        leaf.BizTag,
		Step:       leaf.Step,
		CurrentPos: 0, // 初始化使用第一块buffer缓存
		Buffer:     make([]*LeafSegment, 0),
		UpdateTime: time.Now(),
		Waiting:    make(map[string][]chan byte), //初始化
		IsPreload:  false,
	}
}

func (l *LeafAlloc) Lock() {
	l.mutex.Lock()
}

func (l *LeafAlloc) Unlock() {
	l.mutex.Unlock()
}

func (l *LeafAlloc) HasSeq() bool {
	if l.Buffer[l.CurrentPos].InitOk && l.Buffer[l.CurrentPos].Cursor < l.Buffer[l.CurrentPos].Max {
		return true
	}
	return false
}

func (l *LeafAlloc) HasID(id uint64) bool {
	return id != 0
}

func (l *LeafAlloc) Wakeup() {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	for _, waitChan := range l.Waiting[l.Key] {
		close(waitChan)
	}
	l.Waiting[l.Key] = l.Waiting[l.Key][:0]
}
