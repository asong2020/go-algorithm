package service

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"asong.cloud/go-algorithm/leaf/dao"
	"asong.cloud/go-algorithm/leaf/model"
)

const (
	MAXSTEP  = 10e6
	Default  = 2000
	MAXRETRY = 3 // 最大重试次数
)

type LeafService struct {
	dao     *dao.LeafDao
	leafSeq *model.LeafSeq
	mutex   sync.Mutex
}

func NewLeafService(dao *dao.LeafDao, leafSeq *model.LeafSeq) *LeafService {
	return &LeafService{
		dao:     dao,
		leafSeq: leafSeq,
	}
}

func (l *LeafService) GetID(ctx context.Context, bizTag string) (uint64, error) {
	// 先去内存中看一下是否已经初始了，未初始化则开启兜底策略初始化一下。
	l.mutex.Lock()
	var err error
	seqs := l.leafSeq.Get(bizTag)
	if seqs == nil {
		// 不存在初始化一下
		seqs, err = l.InitCache(ctx, bizTag)
		if err != nil {
			l.mutex.Unlock()
			return 0, err
		}
	}
	l.mutex.Unlock()

	var id uint64
	id, err = l.NextID(seqs)
	if err != nil {
		return 0, err
	}
	l.leafSeq.Update(bizTag, seqs)

	return id, nil
}

// 第一次使用要初始化也就是把DB中的数据存到内存中,非必须操作，直接使用的话有兜底策略
func (l *LeafService) InitCache(ctx context.Context, bizTag string) (*model.LeafAlloc, error) {
	leaf, err := l.dao.NextSegment(ctx, bizTag)
	if err != nil {
		fmt.Printf("initCache failed; err:%v\n", err)
		return nil, err
	}
	alloc := model.NewLeafAlloc(leaf)
	alloc.Buffer = append(alloc.Buffer, model.NewLeafSegment(leaf))

	_ = l.leafSeq.Add(alloc)
	return alloc, nil
}

func (l *LeafService) PreloadBuffer(ctx context.Context, bizTag string, current *model.LeafAlloc) error {
	for i := 0; i < MAXRETRY; i++ {
		leaf, err := l.dao.NextSegment(ctx, bizTag)
		if err != nil {
			fmt.Printf("preloadBuffer failed; bizTag:%s;err:%v", bizTag, err)
			continue
		}
		segment := model.NewLeafSegment(leaf)
		current.Buffer = append(current.Buffer, segment) // 追加
		l.leafSeq.Update(bizTag, current)
		current.Wakeup()
		break
	}
	current.IsPreload = false
	return nil
}

func (l *LeafService) Create(ctx context.Context, leaf *model.Leaf) error {
	if leaf.Step > MAXSTEP {
		return errors.New("step limit exceeded")
	}
	if len(leaf.Description) == 0 || len(leaf.BizTag) == 0 {
		return errors.New("param invalid")
	}
	// 等于0 则代表使用默认step
	if leaf.Step == 0 {
		leaf.Step = Default
	}
	if leaf.MaxID == 0 {
		leaf.MaxID = 1
	}
	return l.dao.Add(ctx, leaf)
}

func (l *LeafService) UpdateStep(ctx context.Context, step int32, bizTag string) error {
	if step == 0 || len(bizTag) == 0 {
		return errors.New("param invalid")
	}
	return l.dao.UpdateStep(ctx, step, bizTag)
}

func (l *LeafService) NextID(current *model.LeafAlloc) (uint64, error) {
	current.Lock()
	defer current.Unlock()
	var id uint64
	currentBuffer := current.Buffer[current.CurrentPos]
	// 判断当前buffer是否是可用的
	if current.HasSeq() {
		id = atomic.AddUint64(&current.Buffer[current.CurrentPos].Cursor, 1)
		current.UpdateTime = time.Now()
	}

	// 当前号段已下发10%时，如果下一个号段未更新加载，则另启一个更新线程去更新下一个号段
	if currentBuffer.Max-id < uint64(0.9*float32(current.Step)) && len(current.Buffer) <= 1 && !current.IsPreload {
		current.IsPreload = true
		cancel, _ := context.WithTimeout(context.Background(), 3*time.Second)
		go l.PreloadBuffer(cancel, current.Key, current)
	}

	// 第一个buffer的segment使用完成 切换到下一个buffer 并移除现在的buffer
	if id == currentBuffer.Max {
		// 判断第二个buffer是否准备好了(因为上面开启协程去更新下一个号段会出现失败)，准备好了切换  currentPos 永远是0 不管怎么切换
		if len(current.Buffer) > 1 && current.Buffer[current.CurrentPos+1].InitOk {
			current.Buffer = append(current.Buffer[:0], current.Buffer[1:]...)
		}
		// 如果没准备好，直接返回就好了，因为现在已经分配id了, 后面会进行补偿
	}
	// 有id直接返回就可以了
	if current.HasID(id) {
		return id, nil
	}

	// 当前buffer已经没有id可用了，此时补偿线程一定正在运行，我们等待一会
	waitChan := make(chan byte, 1)
	current.Waiting[current.Key] = append(current.Waiting[current.Key], waitChan)
	// 释放锁 等待让其他客户端进行走前面的步骤
	current.Unlock()

	timer := time.NewTimer(500 * time.Millisecond) // 等待500ms最多
	select {
	case <-waitChan:
	case <-timer.C:
	}

	current.Lock()
	// 第二个缓冲区仍未初始化好
	if len(current.Buffer) <= 1 {
		return 0, errors.New("get id failed")
	}
	// 切换buffer
	current.Buffer = append(current.Buffer[:0], current.Buffer[1:]...)
	if current.HasSeq() {
		id = atomic.AddUint64(&current.Buffer[current.CurrentPos].Cursor, 1)
		current.UpdateTime = time.Now()
	}
	return id, nil

}
