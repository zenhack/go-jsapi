package promise

import (
	"syscall/js"

	"zenhack.net/go/jsapi"
	"zenhack.net/go/util/orerr"
	"zenhack.net/go/util/thunk"
)

type Promise[T ~jsapi.ValueKind] struct {
	Value js.Value
}

func NewPromise[T ~jsapi.ValueKind](use func(resolve func(T), reject func(error))) Promise[T] {
	var useJs js.Func
	useJs = js.FuncOf(func(this js.Value, args []js.Value) any {
		defer useJs.Release()
		use(
			func(result T) {
				args[0].Invoke(jsapi.Value(result).Value)
			},
			func(err error) {
				args[1].Invoke(jsapi.WrapError(err))
			},
		)
		return nil
	})
	return Promise[T]{
		Value: js.Global().Get("Promise").New(useJs),
	}
}

func Ready[T ~jsapi.ValueKind](v T) Promise[T] {
	return Promise[T]{
		Value: js.Global().Get("Promise").Call("ready", jsapi.Value(v).Value),
	}
}

func Then[A, B ~jsapi.ValueKind](
	pa Promise[A],
	onOk func(A) Promise[B],
	onError func(error) Promise[B],
) Promise[B] {
	var (
		onOkJs    js.Func
		onErrorJs js.Func
	)
	onOkJs = js.FuncOf(func(this js.Value, args []js.Value) any {
		defer onOkJs.Release()
		return onOk(A(jsapi.Value{args[0]})).Value
	})
	if onError == nil {
		return Promise[B]{
			Value: pa.Value.Call("then", onOkJs),
		}
	}

	onErrorJs = js.FuncOf(func(this js.Value, args []js.Value) any {
		defer onErrorJs.Release()
		return onError(jsapi.Error{Value: args[0]}).Value
	})
	return Promise[B]{
		Value: pa.Value.Call("then", onOkJs, onErrorJs),
	}
}

// Wait blocks on the promise, returning the value and any error when it resolves.
func (p Promise[T]) Wait() (T, error) {
	return p.Thunk().Force().Get()
}

func (p Promise[T]) Thunk() *thunk.Thunk[orerr.OrErr[T]] {
	result, fulfill := thunk.Promise[orerr.OrErr[T]]()
	Then(
		p,
		func(v T) Promise[jsapi.Value] {
			fulfill(orerr.New(v, nil))
			return Ready(jsapi.Value{})
		},
		func(err error) Promise[jsapi.Value] {
			fulfill(orerr.New(T{}, err))
			return Ready(jsapi.Value{})
		},
	)
	return result
}
