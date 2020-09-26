package wasman

import (
	"github.com/c0mm4nd/wasman/instr"
)

type wasmContext struct {
	PC         uint64
	Func       *wasmFunc
	Locals     []uint64
	LabelStack *labelStack
}

var instructions = [256]func(ins *Instance){
	instr.OpCodeUnreachable:       func(ins *Instance) { panic("unreachable") }, // TODO: avoid panic
	instr.OpCodeNop:               func(ins *Instance) {},
	instr.OpCodeBlock:             block,
	instr.OpCodeLoop:              loop,
	instr.OpCodeIf:                ifOp,
	instr.OpCodeElse:              elseOp,
	instr.OpCodeEnd:               end,
	instr.OpCodeBr:                br,
	instr.OpCodeBrIf:              brIf,
	instr.OpCodeBrTable:           brTable,
	instr.OpCodeReturn:            func(ins *Instance) {},
	instr.OpCodeCall:              call,
	instr.OpCodeCallIndirect:      callIndirect,
	instr.OpCodeDrop:              drop,
	instr.OpCodeSelect:            selectOp,
	instr.OpCodeLocalGet:          getLocal,
	instr.OpCodeLocalSet:          setLocal,
	instr.OpCodeLocalTee:          teeLocal,
	instr.OpCodeGlobalGet:         getGlobal,
	instr.OpCodeGlobalSet:         setGlobal,
	instr.OpCodeI32Load:           i32Load,
	instr.OpCodeI64Load:           i64Load,
	instr.OpCodeF32Load:           f32Load,
	instr.OpCodeF64Load:           f64Load,
	instr.OpCodeI32Load8s:         i32Load8s,
	instr.OpCodeI32Load8u:         i32Load8u,
	instr.OpCodeI32Load16s:        i32Load16s,
	instr.OpCodeI32Load16u:        i32Load16u,
	instr.OpCodeI64Load8s:         i64Load8s,
	instr.OpCodeI64Load8u:         i64Load8u,
	instr.OpCodeI64Load16s:        i64Load16s,
	instr.OpCodeI64Load16u:        i64Load16u,
	instr.OpCodeI64Load32s:        i64Load32s,
	instr.OpCodeI64Load32u:        i64Load32u,
	instr.OpCodeI32Store:          i32Store,
	instr.OpCodeI64Store:          i64Store,
	instr.OpCodeF32Store:          f32Store,
	instr.OpCodeF64Store:          f64Store,
	instr.OpCodeI32Store8:         i32Store8,
	instr.OpCodeI32Store16:        i32Store16,
	instr.OpCodeI64Store8:         i64Store8,
	instr.OpCodeI64Store16:        i64Store16,
	instr.OpCodeI64Store32:        i64Store32,
	instr.OpCodeMemorySize:        memorySize,
	instr.OpCodeMemoryGrow:        memoryGrow,
	instr.OpCodeI32Const:          i32Const,
	instr.OpCodeI64Const:          i64Const,
	instr.OpCodeF32Const:          f32Const,
	instr.OpCodeF64Const:          f64Const,
	instr.OpCodeI32eqz:            i32eqz,
	instr.OpCodeI32eq:             i32eq,
	instr.OpCodeI32ne:             i32ne,
	instr.OpCodeI32lts:            i32lts,
	instr.OpCodeI32ltu:            i32ltu,
	instr.OpCodeI32gts:            i32gts,
	instr.OpCodeI32gtu:            i32gtu,
	instr.OpCodeI32les:            i32les,
	instr.OpCodeI32leu:            i32leu,
	instr.OpCodeI32ges:            i32ges,
	instr.OpCodeI32geu:            i32geu,
	instr.OpCodeI64eqz:            i64eqz,
	instr.OpCodeI64eq:             i64eq,
	instr.OpCodeI64ne:             i64ne,
	instr.OpCodeI64lts:            i64lts,
	instr.OpCodeI64ltu:            i64ltu,
	instr.OpCodeI64gts:            i64gts,
	instr.OpCodeI64gtu:            i64gtu,
	instr.OpCodeI64les:            i64les,
	instr.OpCodeI64leu:            i64leu,
	instr.OpCodeI64ges:            i64ges,
	instr.OpCodeI64geu:            i64geu,
	instr.OpCodeF32eq:             f32eq,
	instr.OpCodeF32ne:             f32ne,
	instr.OpCodeF32lt:             f32lt,
	instr.OpCodeF32gt:             f32gt,
	instr.OpCodeF32le:             f32le,
	instr.OpCodeF32ge:             f32ge,
	instr.OpCodeF64eq:             f64eq,
	instr.OpCodeF64ne:             f64ne,
	instr.OpCodeF64lt:             f64lt,
	instr.OpCodeF64gt:             f64gt,
	instr.OpCodeF64le:             f64le,
	instr.OpCodeF64ge:             f64ge,
	instr.OpCodeI32clz:            i32clz,
	instr.OpCodeI32ctz:            i32ctz,
	instr.OpCodeI32popcnt:         i32popcnt,
	instr.OpCodeI32add:            i32add,
	instr.OpCodeI32sub:            i32sub,
	instr.OpCodeI32mul:            i32mul,
	instr.OpCodeI32divs:           i32divs,
	instr.OpCodeI32divu:           i32divu,
	instr.OpCodeI32rems:           i32rems,
	instr.OpCodeI32remu:           i32remu,
	instr.OpCodeI32and:            i32and,
	instr.OpCodeI32or:             i32or,
	instr.OpCodeI32xor:            i32xor,
	instr.OpCodeI32shl:            i32shl,
	instr.OpCodeI32shrs:           i32shrs,
	instr.OpCodeI32shru:           i32shru,
	instr.OpCodeI32rotl:           i32rotl,
	instr.OpCodeI32rotr:           i32rotr,
	instr.OpCodeI64clz:            i64clz,
	instr.OpCodeI64ctz:            i64ctz,
	instr.OpCodeI64popcnt:         i64popcnt,
	instr.OpCodeI64add:            i64add,
	instr.OpCodeI64sub:            i64sub,
	instr.OpCodeI64mul:            i64mul,
	instr.OpCodeI64divs:           i64divs,
	instr.OpCodeI64divu:           i64divu,
	instr.OpCodeI64rems:           i64rems,
	instr.OpCodeI64remu:           i64remu,
	instr.OpCodeI64and:            i64and,
	instr.OpCodeI64or:             i64or,
	instr.OpCodeI64xor:            i64xor,
	instr.OpCodeI64shl:            i64shl,
	instr.OpCodeI64shrs:           i64shrs,
	instr.OpCodeI64shru:           i64shru,
	instr.OpCodeI64rotl:           i64rotl,
	instr.OpCodeI64rotr:           i64rotr,
	instr.OpCodeF32abs:            f32abs,
	instr.OpCodeF32neg:            f32neg,
	instr.OpCodeF32ceil:           f32ceil,
	instr.OpCodeF32floor:          f32floor,
	instr.OpCodeF32trunc:          f32trunc,
	instr.OpCodeF32nearest:        f32nearest,
	instr.OpCodeF32sqrt:           f32sqrt,
	instr.OpCodeF32add:            f32add,
	instr.OpCodeF32sub:            f32sub,
	instr.OpCodeF32mul:            f32mul,
	instr.OpCodeF32div:            f32div,
	instr.OpCodeF32min:            f32min,
	instr.OpCodeF32max:            f32max,
	instr.OpCodeF32copysign:       f32copysign,
	instr.OpCodeF64abs:            f64abs,
	instr.OpCodeF64neg:            f64neg,
	instr.OpCodeF64ceil:           f64ceil,
	instr.OpCodeF64floor:          f64floor,
	instr.OpCodeF64trunc:          f64trunc,
	instr.OpCodeF64nearest:        f64nearest,
	instr.OpCodeF64sqrt:           f64sqrt,
	instr.OpCodeF64add:            f64add,
	instr.OpCodeF64sub:            f64sub,
	instr.OpCodeF64mul:            f64mul,
	instr.OpCodeF64div:            f64div,
	instr.OpCodeF64min:            f64min,
	instr.OpCodeF64max:            f64max,
	instr.OpCodeF64copysign:       f64copysign,
	instr.OpCodeI32wrapI64:        i32wrapi64,
	instr.OpCodeI32truncf32s:      i32truncf32s,
	instr.OpCodeI32truncf32u:      i32truncf32u,
	instr.OpCodeI32truncf64s:      i32truncf64s,
	instr.OpCodeI32truncf64u:      i32truncf64u,
	instr.OpCodeI64Extendi32s:     i64extendi32s,
	instr.OpCodeI64Extendi32u:     i64extendi32u,
	instr.OpCodeI64TruncF32s:      i64truncf32s,
	instr.OpCodeI64TruncF32u:      i64truncf32u,
	instr.OpCodeI64Truncf64s:      i64truncf64s,
	instr.OpCodeI64Truncf64u:      i64truncf64u,
	instr.OpCodeF32Converti32s:    f32converti32s,
	instr.OpCodeF32Converti32u:    f32converti32u,
	instr.OpCodeF32Converti64s:    f32converti64s,
	instr.OpCodeF32Converti64u:    f32converti64u,
	instr.OpCodeF32Demotef64:      f32demotef64,
	instr.OpCodeF64Converti32s:    f64converti32s,
	instr.OpCodeF64Converti32u:    f64converti32u,
	instr.OpCodeF64Converti64s:    f64converti64s,
	instr.OpCodeF64Converti64u:    f64converti64u,
	instr.OpCodeF64Promotef32:     f64promotef32,
	instr.OpCodeI32reinterpretf32: func(ins *Instance) {},
	instr.OpCodeI64reinterpretf64: func(ins *Instance) {},
	instr.OpCodeF32reinterpreti32: func(ins *Instance) {},
	instr.OpCodeF64reinterpreti64: func(ins *Instance) {},
}