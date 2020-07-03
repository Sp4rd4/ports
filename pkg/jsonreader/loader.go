package jsonreader

import (
	"io"

	jsoniter "github.com/json-iterator/go"
	"github.com/sp4rd4/ports/pkg/domain"
	"github.com/sp4rd4/ports/pkg/domain/loader"
)

type loaderJson struct {
	reader     io.Reader
	bufferSize int
	cancel     <-chan struct{}
}

func NewLoader(reader io.Reader, bufferSize int, cancel <-chan struct{}) loader.Ports {
	return loaderJson{
		reader:     reader,
		bufferSize: bufferSize,
		cancel:     cancel,
	}
}

func (lj loaderJson) Load() <-chan *domain.Port {
	iter := jsoniter.Parse(jsoniter.ConfigFastest, lj.reader, lj.bufferSize)
	data := make(chan *domain.Port)

	if lj.cancel == nil {
		go lj.iterate(iter, data)
	} else {
		go lj.iterateCancellable(iter, data)
	}

	return data
}

func (lj loaderJson) iterate(iter *jsoniter.Iterator, data chan *domain.Port) {
	defer close(data)

	for field := iter.ReadObject(); field != ""; field = iter.ReadObject() {
		if iter.Error != nil {
			return
		}
		p := &domain.Port{}
		iter.ReadVal(p)
		if iter.Error != nil {
			return
		}
		p.ID = field
		data <- p
	}
}

func (lj loaderJson) iterateCancellable(iter *jsoniter.Iterator, data chan *domain.Port) {
	defer close(data)

	for field := iter.ReadObject(); field != ""; field = iter.ReadObject() {
		if iter.Error != nil {
			return
		}
		select {
		case <-lj.cancel:
			return
		default:
		}

		p := &domain.Port{}
		iter.ReadVal(p)
		if iter.Error != nil {
			return
		}
		p.ID = field
		data <- p
	}
}
