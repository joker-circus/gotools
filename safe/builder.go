package safe

import (
	"strings"
	"sync"
)

type Builder struct {
	b  strings.Builder
	rw sync.RWMutex
}

func NewBuilder() *Builder {
	return &Builder{}
}

func (b *Builder) String() string {
	b.rw.RLock()
	defer b.rw.RUnlock()
	return b.b.String()
}

func (b *Builder) Len() int {
	b.rw.RLock()
	defer b.rw.RUnlock()
	return b.b.Len()
}

func (b *Builder) Cap() int {
	b.rw.RLock()
	defer b.rw.RUnlock()
	return b.b.Cap()
}

func (b *Builder) Reset() {
	b.rw.Lock()
	defer b.rw.Unlock()
	b.b.Reset()
}

func (b *Builder) Write(p []byte) (n int, err error) {
	b.rw.Lock()
	defer b.rw.Unlock()
	return b.b.Write(p)
}

func (b *Builder) WriteByte(c byte) error {
	b.rw.Lock()
	defer b.rw.Unlock()
	return b.b.WriteByte(c)
}

func (b *Builder) WriteRune(r rune) (int, error) {
	b.rw.Lock()
	defer b.rw.Unlock()
	return b.b.WriteRune(r)
}

func (b *Builder) WriteString(s string) (int, error) {
	b.rw.Lock()
	defer b.rw.Unlock()
	return b.b.WriteString(s)
}
