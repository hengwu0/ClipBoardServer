package algorithm

import (
	"container/list"
	"log"
	"os"
	"sync"
)

type Alist struct {
	list   *list.List //list列表存储Node
	locker *sync.Mutex
	hash   string
	file   *os.File
	logger *log.Logger
}

func (l *Alist) New(hash string) *Alist { return &Alist{list.New(), &sync.Mutex{}, hash, nil, nil} }

func (l *Alist) Front() *list.Element { return l.list.Front() }

func (l *Alist) PushBack(v interface{}) { l.locker.Lock(); l.list.PushBack(v); l.locker.Unlock() }

func (l *Alist) Len() int { l.locker.Lock(); lenth := l.list.Len(); l.locker.Unlock(); return lenth }

func (l *Alist) Remove(e *list.Element) {
	l.locker.Lock()
	l.list.Remove(e)
	e.Value.(*Node).conn.Close(logger)
	l.locker.Unlock()
}

func (l *Alist) CloseLog() {
	l.locker.Lock()
	l.logger = nil
	if l.file != nil {
		l.file.Close()
		l.file = nil
	}
	l.locker.Unlock()
}
