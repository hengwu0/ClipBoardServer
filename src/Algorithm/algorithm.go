//要求所有函数、方法线程安全
package algorithm

import (
	"Protocol"
	"sync"
	"time"
)

var clientMap = make(map[string]*Alist)
var mapLocker sync.Mutex

func NewClient(conn *protocol.Conn) *Node {
	var netlist *Alist
	var ok bool
	var hash string

	logger.Printf("Connected in: %s--", conn.GetAddr())
	if hash = conn.AccessGetHash(logger); hash == "" {
		logger.Printf("Access %s: Rejected!!!\n", conn.GetAddr())
		time.Sleep(time.Second * 60)
		conn.Close(logger)
		return nil
	}
	logger.Printf("Access %s Pass!!!\n", conn.GetAddr())

	mapLocker.Lock()
	if netlist, ok = clientMap[hash]; !ok {
		netlist = netlist.New(hash)
		clientMap[hash] = netlist
	}
	mapLocker.Unlock()
	node := new(Node)
	node.cycleRoot = netlist
	conn.SetLogger(netlist.GetLogger())
	node.conn = conn
	netlist.PushBack(node)
	return node
}

type Node struct {
	cycleRoot *Alist
	conn      *protocol.Conn
}

//go ParseCmd
func (node *Node) ParseCmd() {
	tOld := time.Now()
	tNow := tOld
	buf := make([]byte, protocol.PackheadSize)
	for {
		head, size := node.conn.RecvCmd(buf)
		switch head {
		case 'P':
			if ok := node.conn.SendActive(); !ok {
				node.Remove()
				return
			}
		case 'M':
			if ok, clip := node.conn.GetClip(size); !ok {
				node.Remove()
				return
			} else if len(clip) > 0 {
				//跳过小于0.5秒的操作
				tNow = time.Now()
				if tNow.Sub(tOld) > time.Second {
					node.deliver(clip)
				} else {
					node.conn.GetLogger().Printf("Ignored! %s: Duration < 1s!!!\n", node.conn.GetAddr())
				}
				tOld = tNow
			}
		default: //数据错误，关闭连接
			node.Remove()
			return
		}
	}
}

func (node *Node) deliver(clip []byte) {
	e := node.cycleRoot.Front()
	next := e
	for ; e != nil; e = next {
		next = e.Next() //防止在删除节点操作
		if e.Value == node {
			continue
		}
		go func(node *Node) {
			if ok := node.conn.SendClip(clip); !ok {
				logger.Printf("[Debug](%s) node.Remove in deliver\n", node.conn.GetAddr())
				node.cycleRoot.Remove(e)
			}
			logger.Printf("[Debug] node.delived to %s\n", node.conn.GetAddr())
		}(e.Value.(*Node))
	}
}

func (node *Node) Remove() {
	node.conn.Close(logger)
	for e := node.cycleRoot.Front(); e != nil; e = e.Next() {
		if e.Value == node {
			node.cycleRoot.Remove(e)
			if node.cycleRoot.Len() == 0 {
				logger.Printf("[Debug] mapLocker(%s) deleted\n", node.cycleRoot.hash)
				mapLocker.Lock()
				delete(clientMap, node.cycleRoot.hash)
				mapLocker.Unlock()
				node.cycleRoot.CloseLog()
			}
			return
		}
	}
}
