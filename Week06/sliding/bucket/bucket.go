package bucket

import (
	"context"
	"errors"
	"log"
	"sync"
	"time"
)

type Bucket struct {
	counter uint64
	prev    *Bucket
	next    *Bucket
}

type BucketList struct {
	size uint64
	head *Bucket
	tail *Bucket
}

func NewBucketList() *BucketList {
	return &BucketList{}
}

func (l *BucketList) Append(counter uint64) {
	// 创建一个新节点
	newBucket := &Bucket{counter: counter}
	// 只有一个节点
	if l.size == 0 {
		l.head = newBucket
		l.tail = newBucket
	} else {
		newBucket.prev = l.tail
		l.tail.next = newBucket
		l.tail = newBucket
	}
	l.size++
}

func (l *BucketList) RemoveFirst() error {
	if l.size == 0 {
		return errors.New("list is empty")
	} else if l.size == 1 {
		l.head = nil
		l.size = 0
	} else {
		l.head = l.head.next
		l.size--
	}
	return nil
}

type RollingTimeWindow struct {
	Mutex sync.RWMutex

	Counter    uint64
	Limit      uint64
	Split      uint64
	IsLimit    bool
	BucketList *BucketList
}

func NewRollingTimeWindow(limit uint64) *RollingTimeWindow {
	return &RollingTimeWindow{
		Split:   10,
		Limit:   limit,
		IsLimit: false,

		BucketList: NewBucketList(),
	}
}

func (w *RollingTimeWindow) Check(ctx context.Context) error {
	for {
		w.BucketList.Append(w.Counter)
		if w.BucketList.size > w.Split {
			err := w.BucketList.RemoveFirst()
			if err != nil {
				return err
			}
		}
		log.Println("tail counter: ", w.BucketList.tail.counter, "head counter: ", w.BucketList.head.counter)
		if w.BucketList.tail.counter-w.BucketList.head.counter > w.Limit {
			log.Println("limiting...")
			w.IsLimit = true
			return errors.New("limit reached")
		} else {
			w.IsLimit = false
		}
		duration := time.Duration(1000 / w.Split)
		time.Sleep(duration * time.Millisecond)
	}
}
