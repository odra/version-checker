package config

import (
	"io"
	"io/ioutil"
)

func DefaultLoader(data []byte) *Loader {
	return &Loader{
		data: data,
	}
}

func FileLoader(path string) (*Loader, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return &Loader{}, err
	}

	return DefaultLoader(b), nil
}

func (l *Loader) Read(p []byte) (n int, err error) {
	if l.index >= int64(len(l.data)) {
		err = io.EOF
		return
	}

	n = copy(p, l.data[l.index:])
	l.index += int64(n)
	return
}
