package safe

import (
	"bytes"
	"sync"
)

type Buffer struct {
	b  bytes.Buffer
	rw sync.RWMutex
}

func NewBuffer() *Buffer {
	return &Buffer{}
}

func (b *Buffer) Read(p []byte) (n int, err error) {
	b.rw.RLock()
	defer b.rw.RUnlock()
	return b.b.Read(p)
}

func (b *Buffer) String() string {
	b.rw.RLock()
	defer b.rw.RUnlock()
	return b.b.String()
}

func (b *Buffer) Write(p []byte) (n int, err error) {
	b.rw.Lock()
	defer b.rw.Unlock()
	return b.b.Write(p)
}

func (b *Buffer) WriteString(s string) (n int, err error) {
	b.rw.Lock()
	defer b.rw.Unlock()
	return b.b.WriteString(s)
}
