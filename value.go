package jsapi

import "syscall/js"

type ValueKind = struct {
	Value js.Value
}

type Value ValueKind
