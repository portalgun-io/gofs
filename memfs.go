package gofs

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
)

type MemFS struct {
	fs map[string]*bytes.Buffer
}

func (m *MemFS) Open(path string) (io.ReadCloser, error) {
	contents := bytes.NewBuffer(m.fs[path].Bytes())
	return ioutil.NopCloser(contents), nil
}

func (m *MemFS) ReadAll(path string) ([]byte, error) {
	data, ok := m.fs[path]
	if !ok {
		return nil, fmt.Errorf("unable to find file[%s]", path)
	}
	return data.Bytes(), nil
}

func (m *MemFS) Create(path string) (io.WriteCloser, error) {
	buf := &bytes.Buffer{}
	m.fs[path] = buf
	return newWriterNopCloser(buf), nil
}

func (m *MemFS) WriteAll(path string, contents []byte) error {
	buf := &bytes.Buffer{}
	io.Copy(buf, bytes.NewReader(contents))
	m.fs[path] = buf
	return nil
}

func (m *MemFS) Remove(path string) error {
	delete(m.fs, path)
	return nil
}

func NewMemFS() *MemFS {
	return &MemFS{
		fs: map[string]*bytes.Buffer{},
	}
}

type writerNopCloser struct {
	io.Writer
}

func (writerNopCloser) Close() error { return nil }

func newWriterNopCloser(w io.Writer) io.WriteCloser {
	return writerNopCloser{w}
}
