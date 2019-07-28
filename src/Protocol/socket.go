package protocol

import (
	"Register"
	"io"
	"log"
	"net"
	"sync"
	"time"
)

var listen net.Listener

type Conn struct {
	conn    net.Conn
	logger  *log.Logger
	lockerW sync.Mutex
}

func (c *Conn) SetLogger(logger *log.Logger) {
	c.logger = logger
}

func (c *Conn) AccessGetHash(logger *log.Logger) (hash string) {
	if !c.SendAccess() {
		logger.Printf("(%s)hash check: %s", c.GetAddr(), "Send Access Failed!")
		return ""
	}
	c.conn.SetReadDeadline(time.Now().Add(time.Duration(1) * time.Second))
	buf := make([]byte, PackheadSize)
	if h, s := c.RecvCmd(buf); h != 'C' || s == 0 {
		logger.Printf("(%s)hash check: Recv Access Head Failed:%c, %d!", c.GetAddr(), h, s)
		return ""
	} else {
		buf = make([]byte, s)
		if _, err := io.ReadFull(c.conn, buf); err != nil {
			logger.Printf("(%s)hash check: %s", c.GetAddr(), "Recv Access Clip Failed!")
			return ""
		} else {
			version, key := DepackClip(buf)
			logger.Printf("Access Recv: V(%d)%v\n", version, key)
			var ok int
			if ok, hash = register.CheckVandK(logger, version, key); ok != 0 {
				switch ok {
				case 1:
					c.SendUpdate()
					logger.Printf("(%s)hash check for Version<3157553 Failed: V(%d)!", c.GetAddr(), version)
				case 2:
					c.SendKeyErr()
					logger.Printf("(%s)hash check for Key Failed!", c.GetAddr())
				}
				return ""
			}
		}
	}
	c.conn.SetReadDeadline(time.Time{})
	logger.Printf("(%s)hash check: \"%v\"", c.GetAddr(), hash)
	return hash
}

func (c *Conn) GetAddr() string {
	return c.conn.RemoteAddr().String()
}

func (c *Conn) Close(logger *log.Logger) {
	c.lockerW.Lock()
	if c.conn != nil {
		logger.Printf("Connection Closed: %s--", c.GetAddr())
		c.conn.Close()
		c.conn = nil
	}
	c.lockerW.Unlock()
}

func (c *Conn) RecvCmd(buf []byte) (flag byte, size int) {
	flag = 'Z'
	_, err := io.ReadFull(c.conn, buf)
	if err != nil {
		c.logger.Printf("(%s)Read cmd break, err:%v", c.conn.RemoteAddr().String(), err)
		return
	}
	flag, size = Depack(buf)
	if c.logger != nil {
		c.logger.Printf("(%s)recv cmd: flag=%c, size=%d", c.conn.RemoteAddr(), flag, size)
	}
	return
}

func (c *Conn) GetClip(size int) (bool, []byte) {
	buf := make([]byte, size)

	c.conn.SetReadDeadline(time.Now().Add(time.Duration(1) * time.Second))
	if _, err := io.ReadFull(c.conn, buf); err != nil {
		c.logger.Printf("Recv clip break: IP=%s, err:%v", c.conn.RemoteAddr().String(), err)
		return false, nil
	}
	c.conn.SetReadDeadline(time.Time{})

	ftype, clip := DepackClip(buf)
	c.logger.Printf("recv clip: %s: type:%d %s", c.conn.RemoteAddr().String(), ftype, clip)
	return true, buf
}

func (c *Conn) SendClip(clip []byte) bool {
	return c.sendAny(EnpackClip(clip))
}

func (c *Conn) SendActive() bool {
	return c.sendAny(Enpackhead('P', []byte{200}))
}

func (c *Conn) SendAccess() bool {
	return c.sendAny(Enpackhead('C', []byte{'V'}))
}

func (c *Conn) SendUpdate() bool {
	return c.sendAny(Enpackhead('C', []byte{'U'}))
}

func (c *Conn) SendKeyErr() bool {
	return c.sendAny(Enpackhead('C', []byte{'S'}))
}

func (c *Conn) sendAny(date []byte) bool {
	c.lockerW.Lock()
	for lenth, sended := len(date), 0; sended < lenth; {
		c.conn.SetWriteDeadline(time.Now().Add(time.Duration(1) * time.Second))
		if ret, err := c.conn.Write(date); err != nil {
			if c.logger != nil {
				c.logger.Printf("SendActive break: IP=%s, err:%v", c.conn.RemoteAddr().String(), err)
			}
			c.lockerW.Unlock()
			return false
		} else {
			c.conn.SetReadDeadline(time.Time{})
			sended += ret
		}
	}
	c.lockerW.Unlock()
	return true
}

func Listen() (err error) {
	listen, err = net.Listen("tcp", ":7223")
	return err
}

func Accept() (*Conn, error) {
	conn, err := listen.Accept()
	if err != nil {
		time.Sleep(time.Second)
		return nil, err
	}
	return &Conn{conn: conn}, nil
}
