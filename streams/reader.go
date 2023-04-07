package streams

import (
	"zenhack.net/go/jsapi"
)

type ReadableStreamDefaultReader jsapi.Value

type ReadResult struct {
	Value jsapi.Value // TODO: type for the chunk.
	Done  bool
}

/*
func (r ReadableStreamDefaultReader) Read() promise.Promise[ReadResult] {
	panic("TODO")
}
*/
