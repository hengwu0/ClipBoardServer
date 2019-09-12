package protocol

import (
	"encoding/binary"
	"unsafe"
)

const (
	constHeaderLength   = 1
	constModesizeLength = 4
	PackheadSize        = constHeaderLength + constModesizeLength
)

//封包
func EnpackClip(content []byte) (packet []byte) {
	head := make([]byte, constHeaderLength+constModesizeLength)
	head[0] = 'M'
	binary.LittleEndian.PutUint32(head[constHeaderLength:constHeaderLength+constModesizeLength], uint32(len(content)))
	return append(head, content...)
}

func DepackClip(content []byte) (int, []byte) {
	fmode := int(binary.LittleEndian.Uint32(content[:4]))
	clip := make([]byte, len(content)-4)
	copy(clip, content[4:])
	fastReverse(clip)
	return fmode, clip
}

const wordSize = int(unsafe.Sizeof(uintptr(0)))

func fastReverse(src []byte) {
	n := len(src)

	if w := n / wordSize; w > 0 {
		dw := *(*[]uintptr)(unsafe.Pointer(&src))
		for i := 0; i < w; i++ {
			dw[i] = ^dw[i]
		}
	}

	for i := n - n%wordSize; i < n; i++ {
		src[i] = ^src[i]
	}
}

func DepackCmd(content []byte) string {
	return string(content)
}

//自动截断，只取前4位
func Enpackhead(head byte, stat []byte) []byte {
	packet := make([]byte, PackheadSize)
	packet[0] = head
	copy(packet[1:], stat)
	return packet
}

//解包
func Depack(buffer []byte) (flag byte, size int) {
	flag = buffer[0]
	switch buffer[0] {
	case 'P', 'M', 'C':
		size = int(binary.LittleEndian.Uint32(buffer[constHeaderLength:PackheadSize]))
	default:
		size = 0
	}
	return
}
