package algorithm

import (
	"Protocol"
	"testing"
	//"fmt"
)

func TestNewClient(t *testing.T) {
	conn := new(protocol.Conn)
	if node := NewClient(conn); node == nil {
		t.Errorf("%s failed! return nil!\n", "NewClient")
	} else {
		if node.cycleRoot == nil || node.conn == nil {
			t.Errorf("%s failed! node=%#v!\n", "NewClient", node)
		} else {
			if node.conn != conn {
				t.Errorf("%s failed! conn=%#v, node.conn=%#v!\n", "NewClient", conn, node.conn)
			}
			if node.cycleRoot.list.Front().Value != node {
				t.Errorf("%s failed! node=%#v, node.cycleRoot.Front=%#v!\n", "NewClient", node, node.cycleRoot.list.Front())
			}
		}

	}
	//
}
