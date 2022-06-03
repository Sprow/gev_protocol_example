package main

import (
	"bytes"
	"github.com/Allenxuxu/gev"
	"github.com/Allenxuxu/ringbuffer"
	"github.com/gobwas/pool/pbytes"
)

type Protocol struct{}

func (p *Protocol) UnPacket(c *gev.Connection, buffer *ringbuffer.RingBuffer) (interface{}, []byte) {
	buf := pbytes.GetLen(buffer.VirtualLength())  // get []bytes with len 'buffer.VirtualLength()' from pool to get all buffer
	defer pbytes.Put(buf)                         // return slice to pool
	_, _ = buffer.VirtualRead(buf)                // read all buffer
	firstDelimIndex := bytes.IndexByte(buf, '\n') // find first delim index
	if firstDelimIndex != -1 {                    // delimiter '\n' present in buffer
		buffer.VirtualRevert()         // revert virtual pointer to start position
		_, _ = buffer.VirtualRead(buf) // take from buffer bytes until first delim '\n'
		buffer.VirtualFlush()
		return nil, buf
	}
	if firstDelimIndex == -1 { // in delim not present in buffer
		buffer.VirtualRevert() // revert virtual pointer to start position
	}
	return nil, nil
}

func (p *Protocol) Packet(c *gev.Connection, data interface{}) []byte {
	return data.([]byte)
}
