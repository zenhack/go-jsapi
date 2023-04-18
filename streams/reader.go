package streams

import (
	"io"
	"syscall/js"

	"zenhack.net/go/jsapi"
	"zenhack.net/go/jsapi/promise"
)

type ReadableStreamDefaultReader jsapi.Value

func (r ReadableStreamDefaultReader) Next() ([]byte, error) {
	v, err := promise.Promise[jsapi.Value]{Value: r.Value.Call("read")}.Wait()
	if err != nil {
		return nil, err
	}
	if v.Value.Get("done").Bool() {
		return nil, io.EOF
	}
	chunk := v.Value.Get("value")
	length := chunk.Length()
	result := make([]byte, length)
	copied := js.CopyBytesToGo(result, chunk)
	if copied != length {
		panic("short copy")
	}
	return result, nil
}

func (r ReadableStreamDefaultReader) WriteTo(w io.Writer) (n int64, err error) {
	for {
		chunk, err := r.Next()
		if err == io.EOF {
			return n, nil
		}
		if err != nil {
			return n, err
		}
		m, err := w.Write(chunk)
		n += int64(m)
		if err != nil {
			return n, err
		}
	}
}
