package wasm

import (
	"math"
	"reflect"

	"github.com/c0mm4nd/wasman/types"
)

// HostFunc is an implement of wasm.Fn,
// which represents all the functions defined under host(golang) environment
type HostFunc struct {
	Signature *types.FuncType // the shape of func (defined by inputs and outputs)

	// Generator is a func defined by other dev which acts as a Generator to the function
	// (generate when NewInstance's func initializing
	Generator func(ins *Instance) interface{}

	// function is the generated func from Generator, should be set at the time of wasm instance creation
	function interface{}
}

func (f *HostFunc) getType() *types.FuncType {
	return f.Signature
}

func (f *HostFunc) call(ins *Instance) error {
	fnVal := reflect.ValueOf(f.function)
	ty := fnVal.Type()
	in := make([]reflect.Value, ty.NumIn())

	for i := len(in) - 1; i >= 0; i-- {
		val := reflect.New(ty.In(i)).Elem()
		raw := ins.OperandStack.Pop()
		kind := ty.In(i).Kind()

		switch kind {
		case reflect.Float64, reflect.Float32:
			val.SetFloat(math.Float64frombits(raw))
		case reflect.Uint32, reflect.Uint64:
			val.SetUint(raw)
		case reflect.Int32, reflect.Int64:
			val.SetInt(int64(raw))
		default:
			return ErrFuncInvalidInputType
		}
		in[i] = val
	}

	for _, val := range fnVal.Call(in) {
		switch val.Kind() {
		case reflect.Float64, reflect.Float32:
			ins.OperandStack.Push(math.Float64bits(val.Float()))
		case reflect.Uint32, reflect.Uint64:
			ins.OperandStack.Push(val.Uint())
		case reflect.Int32, reflect.Int64:
			ins.OperandStack.Push(uint64(val.Int()))
		default:
			return ErrFuncInvalidReturnType
		}
	}

	return nil
}
