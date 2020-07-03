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

type Config struct {
	BufferSize int `env:"JSON_BUFFER" envDefault:"512"`
}

func NewLoader(reader io.Reader, conf Config, cancel <-chan struct{}) loader.Ports {
	return loaderJson{
		reader:     reader,
		bufferSize: conf.BufferSize,
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
		p := &domain.Port{}
		iter.ReadVal(p)
		p.ID = field
		data <- p
	}
}

func (lj loaderJson) iterateCancellable(iter *jsoniter.Iterator, data chan *domain.Port) {
	defer close(data)

	for field := iter.ReadObject(); field != ""; field = iter.ReadObject() {
		select {
		case <-lj.cancel:
			return
		default:
		}

		p := &domain.Port{}
		iter.ReadVal(p)
		p.ID = field
		data <- p
	}
}
