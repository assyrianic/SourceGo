package SPTools

import (
	"os"
	"fmt"
	"unsafe"
	"strconv"
	"unicode/utf8"
	///"time"
)


/*
 * https://github.com/alliedmodders/sourcepawn/blob/master/compiler/types.cpp
 * https://github.com/alliedmodders/sourcepawn/blob/master/compiler/types.h
 * 
 * https://github.com/alliedmodders/sourcepawn/blob/master/compiler/expressions.cpp
 * https://github.com/alliedmodders/sourcepawn/blob/master/compiler/expressions.h
 * 
 * https://github.com/alliedmodders/sourcepawn/blob/master/compiler/type-checker.cpp
 */


/*
 * Check that identifiers are declared before to be used in computations.
 * Check that reserved keywords are not misused.
 * Check that types are correctly declared, if the language is explicitly typed.
 * Check that computations are type-consistent, wherever possible.
 */

/**
type TypeAndVal struct {
	FieldsOrParams map[string]TypeAndVal
	///D Decl
	RefOrArrayOrRet *TypeAndVal
	Len, Offset, StorageClass int
	Kind TypeKind
	HasFields bool
}

type Scope struct {
	Types, Syms map[*Name]Node
	Kids []*Scope
	Parent *Scope
}


func NewScope() *Scope {
	return &Scope{ Types: make(map[*Name]Node), Syms: make(map[*Name]Node) }
}
*/



type (
	TypeAndVal interface {
		aType()
	}
	
	VoidTypeAndVal struct {
		_type
	}
	
	HandleTypeAndVal struct {
		_type
	}
	
	CharTypeAndVal struct {
		Value byte
		_type
	}
	
	// also used as the `any` type.
	IntTypeAndVal struct {
		Value int32
		_type
	}
	
	FloatTypeAndVal struct {
		Value float32
		_type
	}
	
	RefTypeAndVal struct {
		Ref TypeAndVal
		_type
	}
	
	ArrayTypeAndVal struct {
		Elems []TypeAndVal
		Dynamic bool // dynamically stack allocated.
		_type
	}
	
	FuncTypeAndVal struct {
		Params map[string]TypeAndVal // empty/nill if no params.
		Ret TypeAndVal // nil == 'void'
		Variadic bool
		_type
	}
)
type _type struct {}
func (_type) aType() {}


func IsArithmeticTypeAndVal(a TypeAndVal) bool {
	switch a.(type) {
	case IntTypeAndVal, FloatTypeAndVal, CharTypeAndVal:
		return true
	default:
		return false
	}
}

func GetTypeAndValName(a TypeAndVal) string {
	switch t := a.(type) {
	case VoidTypeAndVal:
		return "Void Type"
	case IntTypeAndVal:
		return "Int Type"
	case FloatTypeAndVal:
		return "Float Type"
	case CharTypeAndVal:
		return "Char Type"
	case RefTypeAndVal:
		return "Ref Type of " + GetTypeAndValName(t.Ref)
	case ArrayTypeAndVal:
		return "Array Type"
	case FuncTypeAndVal:
		return "Func Type"
	default:
		return "Unknown Type"
	}
}

/*
 * 1. Two expressions are convertible when their reduced forms are the same. E.g 2 + 2 is convertible to 4
 * 2. Two expressions are coercible when you can safely cast one to the other. E.g 22 : int32 might be coercible to 22 : int64
 */
func AreTypeAndValCoercible(a, b TypeAndVal) bool {
	switch tA := a.(type) {
	case IntTypeAndVal, CharTypeAndVal, FloatTypeAndVal:
		switch b.(type) {
		case IntTypeAndVal, CharTypeAndVal, FloatTypeAndVal:
			return true
		case VoidTypeAndVal, FuncTypeAndVal, ArrayTypeAndVal, RefTypeAndVal:
			return false
		}
	case ArrayTypeAndVal:
		switch tB := b.(type) {
		case ArrayTypeAndVal:
			if len(tA.Elems) != len(tB.Elems) {
				// trying to coerce array with different size.
				return false
			}
			for i := range tA.Elems {
				if !AreTypeAndValCoercible(tA.Elems[i], tB.Elems[i]) {
					return false
				}
			}
			return true
		}
	}
	return false
}

func ConvertToInt(a TypeAndVal) (IntTypeAndVal, bool) {
	if !IsArithmeticTypeAndVal(a) {
		return IntTypeAndVal{}, false
	}
	
	switch tt := a.(type) {
	case IntTypeAndVal:
		return tt, true
	case FloatTypeAndVal:
		return IntTypeAndVal{ Value: int32(tt.Value) }, true
	case CharTypeAndVal:
		return IntTypeAndVal{ Value: int32(tt.Value) }, true
	default:
		return IntTypeAndVal{}, false
	}
}

func ConvertToFloat(a TypeAndVal) (FloatTypeAndVal, bool) {
	if !IsArithmeticTypeAndVal(a) {
		return FloatTypeAndVal{}, false
	}
	
	switch tt := a.(type) {
	case FloatTypeAndVal:
		return tt, true
	case IntTypeAndVal:
		return FloatTypeAndVal{ Value: float32(tt.Value) }, true
	case CharTypeAndVal:
		return FloatTypeAndVal{ Value: float32(tt.Value) }, true
	default:
		return FloatTypeAndVal{}, false
	}
}

func ConvertToChar(a TypeAndVal) (CharTypeAndVal, bool) {
	if !IsArithmeticTypeAndVal(a) {
		return CharTypeAndVal{}, false
	}
	
	switch tt := a.(type) {
	case CharTypeAndVal:
		return tt, true
	case IntTypeAndVal:
		return CharTypeAndVal{ Value: byte(tt.Value) }, true
	case FloatTypeAndVal:
		return CharTypeAndVal{ Value: byte(tt.Value) }, true
	default:
		return CharTypeAndVal{}, false
	}
}


/*
 * Okay so now here's what you do. You need to work out what the type of every expression is. For simple expressions that is a simple operation, 1 is probably an integer, "foo" is probably a string - depends on your language rules.
 *
 * Variables require some tracking, if you see x then you need to lookup the definition of x and see what type it was declared with. That's the type of whatever value is inside of x and therefore the type of whatever value you get by reading the variable x.
 *
 * Finally, you've got expressions like function calls or operators. foo(a, b) and those work by getting the type of all of the sub-expressions (so get the type of a and b separately) then looking up an appropriate definition of foo and the return type of that function is the type of the expression.
 *
 * Several things can go wrong. Simplest example is new FOO:x; new BAR:y; ...; x = y you get to invent the rules on what happens when users write this code but for the sake of example, I'm going to say that this is an error because you're only allowed to assign values of type FOO to variables of type FOO. The way this sort of thing works is that you lookup the type of x and the type of y and figure out if the assignment is allowed to happen.
 *
 * So in this example we lookup x see that it has type FOO, lookup y see that it has type BAR, and oopsie you're not allowed to assign BAR to a FOO. In the compiler you probably want these types to be represented by some structure that gives your type checker all the information it needs to make these decisions, like you mentioned before about some difference in behaviour for named enums vs anonymous enums you would probably want to have a field on the structure your compiler uses to represent the type that tells it if it's an anonymous enum.
 *
 * This is also how you handle things like implicit conversions or whatever else. Let's say that you decide that `new int:a; new float:b; ...; a = b` means to round b down to the nearest integer. When you're processing that assignment statement/expression you see that b has type float and a has type int so you know that you should emit a rounding operation as well as a value copy.
 */


type Interp struct {
	MsgSpan
	Syms, Types map[string]TypeAndVal
}

func MakeInterpreter(p Parser) Interp {
	var i = Interp{ MsgSpan: p.TokenReader.MsgSpan, Types: make(map[string]TypeAndVal) }
	i.Types["int"] = IntTypeAndVal{}
	i.Types["any"] = IntTypeAndVal{}
	i.Types["bool"] = IntTypeAndVal{}
	i.Types["char"] = IntTypeAndVal{}
	i.Types["float"] = FloatTypeAndVal{}
	return i
}

func (interp Interp) DoMessage(n Node, msgtype, color, msg string, args ...any) string {
	t := n.Tok()
	report := interp.MsgSpan.Report(msgtype, "", color, msg, *t.Path, &t.Span.LineStart, &t.Span.ColStart, args...)
	interp.MsgSpan.PurgeNotes()
	return report
}


func (interp Interp) EvalExpr(e Expr) TypeAndVal {
	if e==nil {
		return VoidTypeAndVal{}
	}
	
	switch ast := e.(type) {
	case *BasicLit:
		switch ast.Kind {
		case IntLit:
			int_val, _ := strconv.ParseInt(ast.Value, 0, 32)
			return IntTypeAndVal{ Value: int32(int_val) }
		case BoolLit:
			bool_val, _ := strconv.ParseBool(ast.Value)
			return IntTypeAndVal{ Value: Ternary[int32](bool_val, int32(1), int32(0)) }
		case CharLit:
			r, _ := utf8.DecodeRuneInString(ast.Value)
			return IntTypeAndVal{ Value: int32(r) }
		case FloatLit:
			float_val, _ := strconv.ParseFloat(ast.Value, 32)
			return FloatTypeAndVal{ Value: float32(float_val) }
		case StringLit:
			arr := ArrayTypeAndVal{}
			for i := range ast.Value {
				arr.Elems = append(arr.Elems, CharTypeAndVal{ Value: byte(ast.Value[i]) })
			}
			arr.Elems = append(arr.Elems, CharTypeAndVal{ Value: byte(0) })
			return arr
		}
	case *BracketExpr:
		arr := ArrayTypeAndVal{}
		for i := range ast.Exprs {
			arr.Elems = append(arr.Elems, interp.EvalExpr(ast.Exprs[i]))
		}
		return arr
	case *NullExpr:
		return IntTypeAndVal{ Value: 0 }
	case *ThisExpr: // get type of 'this'.
		return IntTypeAndVal{ Value: 0 }
	case *Name:
		// TODO: syms here.
		// "load" value and return.
		return IntTypeAndVal{ Value: 0 }
	case *UnaryExpr:
		switch ast.Kind {
		case TKSizeof:
			// TODO: syms here.
			return IntTypeAndVal{ Value: 0 }
		case TKIncr, TKDecr:
			t := interp.EvalExpr(ast.X)
			if !IsArithmeticTypeAndVal(t) {
				return VoidTypeAndVal{}
			}
			
			switch tnv := t.(type) {
			case IntTypeAndVal:
				return IntTypeAndVal{ Value: Ternary[int32](ast.Kind==TKIncr, tnv.Value + 1, tnv.Value - 1) }
			case CharTypeAndVal:
				return CharTypeAndVal{ Value: Ternary[byte](ast.Kind==TKIncr, tnv.Value + 1, tnv.Value - 1) }
			case FloatTypeAndVal:
				return FloatTypeAndVal{ Value: Ternary[float32](ast.Kind==TKIncr, tnv.Value + 1, tnv.Value - 1) }
			}
		case TKNot:
			t := interp.EvalExpr(ast.X)
			switch tnv := t.(type) {
			case IntTypeAndVal:
				return IntTypeAndVal{ Value: Ternary[int32](tnv.Value==0, 1, 0) }
			case FloatTypeAndVal:
				return IntTypeAndVal{ Value: Ternary[int32](tnv.Value==0.0, 1, 0) }
			default: // error
				interp.MsgSpan.PrepNote(ast.Span(), "here\n")
				fmt.Fprintf(os.Stdout, interp.DoMessage(e, "type error", COLOR_RED, "Non-Int type for NOT expression."))
			}
		case TKCompl:
			t := interp.EvalExpr(ast.X)
			switch tnv := t.(type) {
			case IntTypeAndVal:
				return IntTypeAndVal{ Value: ^tnv.Value }
			default: // error
				interp.MsgSpan.PrepNote(ast.Span(), "here\n")
				fmt.Fprintf(os.Stdout, interp.DoMessage(e, "type error", COLOR_RED, "Non-Int type for Bitwise NOT/Complement expression."))
			}
		case TKSub:
			t := interp.EvalExpr(ast.X)
			switch tnv := t.(type) {
			case IntTypeAndVal:
				return IntTypeAndVal{ Value: -tnv.Value }
			case FloatTypeAndVal:
				return FloatTypeAndVal{ Value: -tnv.Value }
			default: // error
				interp.MsgSpan.PrepNote(ast.Span(), "here\n")
				fmt.Fprintf(os.Stdout, interp.DoMessage(e, "type error", COLOR_RED, "Non-Numeric type for Negative expression."))
			}
		case TKNew:
			// TODO: error on non-objects.
			// ast.X could be a function call expr or type/name index expr.
		}
	case *IndexExpr:
		// TODO: syms here.
		typeOfX := interp.EvalExpr(ast.X)
		if !IsExactType[ArrayTypeAndVal](typeOfX) {
			///interp.MsgSpan.PrepNote(ast.Span(), "")
			fmt.Fprintf(os.Stdout, interp.DoMessage(e, "runtime error", COLOR_RED, "Attempting to index non-Array type."))
			// error
			return VoidTypeAndVal{}
		}
		
		typeOfIdx := interp.EvalExpr(ast.Index)
		if !IsExactType[IntTypeAndVal](typeOfIdx) {
			///interp.MsgSpan.PrepNote(ast.Span(), "here")
			fmt.Fprintf(os.Stdout, interp.DoMessage(e, "runtime error", COLOR_RED, "Attempting to index Array type with non-Int value."))
			// error
			return VoidTypeAndVal{}
		}
		
		int_idx := typeOfIdx.(IntTypeAndVal)
		if int_idx.Value < 0 {
			interp.MsgSpan.PrepNote(ast.Span(), "here\n")
			fmt.Fprintf(os.Stdout, interp.DoMessage(e, "runtime error", COLOR_RED, "Attempting to index Array type with negative index."))
			// invalid access
			return VoidTypeAndVal{}
		}
		
		switch arr := typeOfX.(type) {
		case ArrayTypeAndVal:
			if int(int_idx.Value) >= len(arr.Elems) {
				interp.MsgSpan.PrepNote(ast.Span(), "here\n")
				fmt.Fprintf(os.Stdout, interp.DoMessage(e, "runtime error", COLOR_RED, "Attempting to index Array type with out-of-bounds index."))
				return VoidTypeAndVal{}
			}
			return arr.Elems[int_idx.Value]
		}
	case *ViewAsExpr:
		// check coercion here.
		type_tok := ast.Type.(*TypedExpr)
		targetType, typeOfX := interp.Types[type_tok.TypeName.Lexeme], interp.EvalExpr(ast.X)
		if AreTypeAndValCoercible(targetType, typeOfX) {
			switch targetType.(type) {
			case IntTypeAndVal:
				r, _ := ConvertToInt(typeOfX)
				return r
			case FloatTypeAndVal:
				r, _ := ConvertToFloat(typeOfX)
				return r
			case CharTypeAndVal:
				r, _ := ConvertToChar(typeOfX)
				return r
			}
		} else {
			interp.MsgSpan.PrepNote(ast.Span(), "here\n")
			fmt.Fprintf(os.Stdout, interp.DoMessage(e, "type error", COLOR_RED, "Cannot coerce %s to %s.", GetTypeAndValName(typeOfX), GetTypeAndValName(targetType)))
			// error, trying to convert to invalid type.
		}
	case *BinExpr:
		l, r := interp.EvalExpr(ast.L), interp.EvalExpr(ast.R)
		// if mixing with char type, promote it to int.
		// if 'any' type, autocast to int.
		if IsExactType[CharTypeAndVal](l) {
			int_type, _ := ConvertToInt(l)
			l = int_type
		}
		if IsExactType[CharTypeAndVal](r) {
			int_type, _ := ConvertToInt(r)
			r = int_type
		}
		
		// if mixing with float type, entire expr is float type.
		if IsExactType[FloatTypeAndVal](l) && !IsExactType[FloatTypeAndVal](r) {
			flt_type, _ := ConvertToFloat(r)
			r = flt_type
		} else if !IsExactType[FloatTypeAndVal](l) && IsExactType[FloatTypeAndVal](r) {
			flt_type, _ := ConvertToFloat(l)
			l = flt_type
		}
		
		if IsExactType[VoidTypeAndVal](l) || IsExactType[VoidTypeAndVal](r) {
			return VoidTypeAndVal{}
		}
		
		//var ref_read, ref_write bool
		switch ast.Kind {
		case TKAdd:
			if AreSameType[FloatTypeAndVal](l, r) {
				fL, fR := GetBinaryTypes[FloatTypeAndVal](l, r)
				return FloatTypeAndVal{ Value: fL.Value + fR.Value }
			} else {
				iL, iR := GetBinaryTypes[IntTypeAndVal](l, r)
				return IntTypeAndVal{ Value: iL.Value + iR.Value }
			}
		case TKSub:
			if AreSameType[FloatTypeAndVal](l, r) {
				fL, fR := GetBinaryTypes[FloatTypeAndVal](l, r)
				return FloatTypeAndVal{ Value: fL.Value - fR.Value }
			} else {
				iL, iR := GetBinaryTypes[IntTypeAndVal](l, r)
				return IntTypeAndVal{ Value: iL.Value - iR.Value }
			}
		case TKMul:
			if AreSameType[FloatTypeAndVal](l, r) {
				fL, fR := GetBinaryTypes[FloatTypeAndVal](l, r)
				return FloatTypeAndVal{ Value: fL.Value * fR.Value }
			} else {
				iL, iR := GetBinaryTypes[IntTypeAndVal](l, r)
				return IntTypeAndVal{ Value: iL.Value * iR.Value }
			}
		case TKDiv:
			if AreSameType[FloatTypeAndVal](l, r) {
				fL, fR := GetBinaryTypes[FloatTypeAndVal](l, r)
				return FloatTypeAndVal{ Value: fL.Value / fR.Value }
			} else {
				iL, iR := GetBinaryTypes[IntTypeAndVal](l, r)
				return IntTypeAndVal{ Value: iL.Value / iR.Value }
			}
		case TKMod:
			if !AreSameType[IntTypeAndVal](l, r) {
				// illegal operation for non-int types.
				interp.MsgSpan.PrepNote(ast.Span(), "here\n")
				fmt.Fprintf(os.Stdout, interp.DoMessage(e, "type error", COLOR_RED, "Cannot do Modulo operation with non-Int types."))
				return VoidTypeAndVal{}
			}
			iL, iR := GetBinaryTypes[IntTypeAndVal](l, r)
			return IntTypeAndVal{ Value: iL.Value % iR.Value }
		case TKAnd:
			if !AreSameType[IntTypeAndVal](l, r) {
				// illegal operation for non-int types.
				interp.MsgSpan.PrepNote(ast.Span(), "here\n")
				fmt.Fprintf(os.Stdout, interp.DoMessage(e, "type error", COLOR_RED, "Cannot do bitwise AND operation with non-Int types."))
				return VoidTypeAndVal{}
			}
			iL, iR := GetBinaryTypes[IntTypeAndVal](l, r)
			return IntTypeAndVal{ Value: iL.Value & iR.Value }
		case TKAndNot:
			if !AreSameType[IntTypeAndVal](l, r) {
				// illegal operation for non-int types.
				interp.MsgSpan.PrepNote(ast.Span(), "here\n")
				fmt.Fprintf(os.Stdout, interp.DoMessage(e, "type error", COLOR_RED, "Cannot do bitwise AND-NOT operation with non-Int types."))
				return VoidTypeAndVal{}
			}
			iL, iR := GetBinaryTypes[IntTypeAndVal](l, r)
			return IntTypeAndVal{ Value: iL.Value &^ iR.Value }
		case TKOr:
			if !AreSameType[IntTypeAndVal](l, r) {
				// illegal operation for non-int types.
				interp.MsgSpan.PrepNote(ast.Span(), "here\n")
				fmt.Fprintf(os.Stdout, interp.DoMessage(e, "type error", COLOR_RED, "Cannot do bitwise OR operation with non-Int types."))
				return VoidTypeAndVal{}
			}
			iL, iR := GetBinaryTypes[IntTypeAndVal](l, r)
			return IntTypeAndVal{ Value: iL.Value | iR.Value }
		case TKXor:
			if !AreSameType[IntTypeAndVal](l, r) {
				// illegal operation for non-int types.
				interp.MsgSpan.PrepNote(ast.Span(), "here\n")
				fmt.Fprintf(os.Stdout, interp.DoMessage(e, "type error", COLOR_RED, "Cannot do bitwise XOR operation with non-Int types."))
				return VoidTypeAndVal{}
			}
			iL, iR := GetBinaryTypes[IntTypeAndVal](l, r)
			return IntTypeAndVal{ Value: iL.Value ^ iR.Value }
		case TKShAL:
			if !AreSameType[IntTypeAndVal](l, r) {
				// illegal operation for non-int types.
				interp.MsgSpan.PrepNote(ast.Span(), "here\n")
				fmt.Fprintf(os.Stdout, interp.DoMessage(e, "type error", COLOR_RED, "Cannot do left bit-shift operation with non-Int types."))
				return VoidTypeAndVal{}
			}
			iL, iR := GetBinaryTypes[IntTypeAndVal](l, r)
			if iR.Value >= 32 {
				for iR.Value >= 32 {
					iR.Value %= 32
				}
				// warn about shifting overflow.
				interp.MsgSpan.PrepNote(ast.Span(), "here\n")
				fmt.Fprintf(os.Stdout, interp.DoMessage(e, "type warning", COLOR_MAGENTA, "Left bit-shift overflows int."))
			} else if iR.Value < 0 {
				// warn about shifting with negative numbers.
				interp.MsgSpan.PrepNote(ast.Span(), "here\n")
				fmt.Fprintf(os.Stdout, interp.DoMessage(e, "type warning", COLOR_MAGENTA, "Left bit-shifting with negative numbers."))
			} else if iR.Value==0 {
				interp.MsgSpan.PrepNote(ast.Span(), "here\n")
				fmt.Fprintf(os.Stdout, interp.DoMessage(e, "type warning", COLOR_MAGENTA, "Left bit-shift has no effect."))
				return IntTypeAndVal{ Value: iL.Value }
			}
			return IntTypeAndVal{ Value: iL.Value << iR.Value }
		case TKShAR:
			if !AreSameType[IntTypeAndVal](l, r) {
				// illegal operation for non-int types.
				interp.MsgSpan.PrepNote(ast.Span(), "here\n")
				fmt.Fprintf(os.Stdout, interp.DoMessage(e, "type error", COLOR_RED, "Cannot do right arithmetic bit-shift operation with non-Int types."))
				return VoidTypeAndVal{}
			}
			iL, iR := GetBinaryTypes[IntTypeAndVal](l, r)
			if iR.Value >= 32 || iR.Value < 0 {
				if uint32(iL.Value) & (1 << 31) > 0 {
					u := uint32(1 << 31)
					interp.MsgSpan.PrepNote(ast.Span(), "here\n")
					fmt.Fprintf(os.Stdout, interp.DoMessage(e, "type warning", COLOR_MAGENTA, "Right arithmetic bit-shift underflows keeps sign-bit."))
					return IntTypeAndVal{ Value: *(*int32)((unsafe.Pointer)(&u)) }
				}
				interp.MsgSpan.PrepNote(ast.Span(), "here\n")
				fmt.Fprintf(os.Stdout, interp.DoMessage(e, "type warning", COLOR_MAGENTA, "Right arithmetic bit-shift underflows to 0."))
				return IntTypeAndVal{ Value: 0 }
			} else if iR.Value==0 {
				interp.MsgSpan.PrepNote(ast.Span(), "here\n")
				fmt.Fprintf(os.Stdout, interp.DoMessage(e, "type warning", COLOR_MAGENTA, "Right arithmetic bit-shift has no effect."))
				return IntTypeAndVal{ Value: iL.Value }
			}
			return IntTypeAndVal{ Value: iL.Value >> iR.Value }
		case TKShLR:
			if !AreSameType[IntTypeAndVal](l, r) {
				// illegal operation for non-int types.
				interp.MsgSpan.PrepNote(ast.Span(), "here\n")
				fmt.Fprintf(os.Stdout, interp.DoMessage(e, "type error", COLOR_RED, "Cannot do right logical bit-shift operation with non-Int types."))
				return VoidTypeAndVal{}
			}
			iL, iR := GetBinaryTypes[IntTypeAndVal](l, r)
			if iR.Value >= 32 || iR.Value < 0 {
				interp.MsgSpan.PrepNote(ast.Span(), "here\n")
				fmt.Fprintf(os.Stdout, interp.DoMessage(e, "type warning", COLOR_MAGENTA, "Right logical bit-shift underflows to 0."))
				return IntTypeAndVal{ Value: 0 }
			} else if iR.Value==0 {
				return IntTypeAndVal{ Value: iL.Value }
			}
			return IntTypeAndVal{ Value: iL.Value >> iR.Value }
		case TKNotEq:
			if AreSameType[FloatTypeAndVal](l, r) {
				fL, fR := GetBinaryTypes[FloatTypeAndVal](l, r)
				return IntTypeAndVal{ Value: Ternary[int32](fL.Value != fR.Value, 1, 0) }
			} else {
				iL, iR := GetBinaryTypes[IntTypeAndVal](l, r)
				return IntTypeAndVal{ Value: Ternary[int32](iL.Value != iR.Value, 1, 0) }
			}
		case TKEq:
			if AreSameType[FloatTypeAndVal](l, r) {
				fL, fR := GetBinaryTypes[FloatTypeAndVal](l, r)
				return IntTypeAndVal{ Value: Ternary[int32](fL.Value == fR.Value, 1, 0) }
			} else {
				iL, iR := GetBinaryTypes[IntTypeAndVal](l, r)
				return IntTypeAndVal{ Value: Ternary[int32](iL.Value == iR.Value, 1, 0) }
			}
		case TKAndL:
			if AreSameType[FloatTypeAndVal](l, r) {
				fL, fR := GetBinaryTypes[FloatTypeAndVal](l, r)
				return IntTypeAndVal{ Value: Ternary[int32](fL.Value > 0.0 && fR.Value > 0.0, 1, 0) }
			} else {
				iL, iR := GetBinaryTypes[IntTypeAndVal](l, r)
				return IntTypeAndVal{ Value: Ternary[int32](iL.Value > 0 && iR.Value > 0, 1, 0) }
			}
		case TKOrL:
			if AreSameType[FloatTypeAndVal](l, r) {
				fL, fR := GetBinaryTypes[FloatTypeAndVal](l, r)
				return IntTypeAndVal{ Value: Ternary[int32](fL.Value > 0.0 || fR.Value > 0.0, 1, 0) }
			} else {
				iL, iR := GetBinaryTypes[IntTypeAndVal](l, r)
				return IntTypeAndVal{ Value: Ternary[int32](iL.Value > 0 || iR.Value > 0, 1, 0) }
			}
			
		// assignments.
		// "write" value into lval.
		case TKAssign:
		case TKAddA:
		case TKSubA:
		case TKMulA:
		case TKDivA:
		case TKModA:
		case TKAndA:
		case TKAndNotA:
		case TKOrA:
		case TKXorA:
		case TKShALA:
		case TKShARA:
		case TKShLRA:
		}
	case *ChainExpr:
		// a # b # c => a # b && b # c
		// a # b ==> a = b; b = next; a # b; repeat.
		a := interp.EvalExpr(ast.A)
		int_res := IntTypeAndVal{ Value: 1 }
		for i := range ast.Kinds {
			b := interp.EvalExpr(ast.Bs[i])
			switch ast.Kinds[i] {
			case TKLess:
				if AreSameType[FloatTypeAndVal](a, b) {
					fL, fR := GetBinaryTypes[FloatTypeAndVal](a, b)
					int_res = IntTypeAndVal{ Value: Ternary[int32](int_res.Value > 0 && fL.Value < fR.Value, 1, 0) }
				} else {
					iL, iR := GetBinaryTypes[IntTypeAndVal](a, b)
					int_res = IntTypeAndVal{ Value: Ternary[int32](int_res.Value > 0 && iL.Value < iR.Value, 1, 0) }
				}
			case TKGreater:
				if AreSameType[FloatTypeAndVal](a, b) {
					fL, fR := GetBinaryTypes[FloatTypeAndVal](a, b)
					int_res = IntTypeAndVal{ Value: Ternary[int32](int_res.Value > 0 && fL.Value > fR.Value, 1, 0) }
				} else {
					iL, iR := GetBinaryTypes[IntTypeAndVal](a, b)
					int_res = IntTypeAndVal{ Value: Ternary[int32](int_res.Value > 0 && iL.Value > iR.Value, 1, 0) }
				}
			case TKGreaterE:
				if AreSameType[FloatTypeAndVal](a, b) {
					fL, fR := GetBinaryTypes[FloatTypeAndVal](a, b)
					int_res = IntTypeAndVal{ Value: Ternary[int32](int_res.Value > 0 && fL.Value >= fR.Value, 1, 0) }
				} else {
					iL, iR := GetBinaryTypes[IntTypeAndVal](a, b)
					int_res = IntTypeAndVal{ Value: Ternary[int32](int_res.Value > 0 && iL.Value >= iR.Value, 1, 0) }
				}
			case TKLessE:
				if AreSameType[FloatTypeAndVal](a, b) {
					fL, fR := GetBinaryTypes[FloatTypeAndVal](a, b)
					int_res = IntTypeAndVal{ Value: Ternary[int32](int_res.Value > 0 && fL.Value <= fR.Value, 1, 0) }
				} else {
					iL, iR := GetBinaryTypes[IntTypeAndVal](a, b)
					int_res = IntTypeAndVal{ Value: Ternary[int32](int_res.Value > 0 && iL.Value <= iR.Value, 1, 0) }
				}
			}
			a = b
		}
		return int_res
	case *TernaryExpr:
		if cond := interp.EvalExpr(ast.A); !IsArithmeticTypeAndVal(cond) {
			interp.MsgSpan.PrepNote(ast.Span(), "here\n")
			fmt.Fprintf(os.Stdout, interp.DoMessage(e, "type error", COLOR_RED, "cannot evaluate non-arithmetic type in ternary condition."))
			// error, need an int value here.
			return VoidTypeAndVal{}
		} else if (IsExactType[IntTypeAndVal](cond) && cond.(IntTypeAndVal).Value != 0) || (IsExactType[FloatTypeAndVal](cond) && cond.(FloatTypeAndVal).Value != 0.0) {
			return interp.EvalExpr(ast.B)
		} else {
			return interp.EvalExpr(ast.C)
		}
	case *CommaExpr:
		var ret_typeval TypeAndVal
		for i := range ast.Exprs {
			ret_typeval = interp.EvalExpr(ast.Exprs[i])
		}
		return ret_typeval
	case *FuncLit:
		// TODO: implement function creation here.
	case *BadExpr:
		return VoidTypeAndVal{}
	}
	return VoidTypeAndVal{}
}


type ControlFlow int8
const (
	// continue execution.
	FLOW_EXC = ControlFlow(iota)
	
	// doing loop 'continue'.
	FLOW_CNT
	
	// doing return.
	FLOW_RET
	
	// breaking out of loop.
	FLOW_BRK
)

func (interp Interp) EvalStmt(s Stmt, flow *ControlFlow) TypeAndVal {
	if s==nil {
		return VoidTypeAndVal{}
	}
	
	switch ast := s.(type) {
	case *BlockStmt:
		// make new scope here.
		for i := range ast.Stmts {
			blk_flow := FLOW_EXC
			tnv := interp.EvalStmt(ast.Stmts[i], &blk_flow)
			if blk_flow != FLOW_EXC {
				if blk_flow==FLOW_RET {
					*flow = blk_flow
				}
				return tnv
			}
		}
	case *WhileStmt:
		blk_flow := FLOW_EXC
		counter := 0
		const inf_protect = 999_999
		if ast.Do {
		do_while:
			tnv_body := interp.EvalStmt(ast.Body, &blk_flow)
			switch blk_flow {
			case FLOW_CNT:
				goto do_while
			case FLOW_RET, FLOW_BRK:
				if blk_flow==FLOW_RET {
					*flow = blk_flow
				}
				return tnv_body
			default:
			}
			
			tnv_cond := interp.EvalExpr(ast.Cond)
			if tnv_conv, res := ConvertToInt(tnv_cond); res {
				tnv_cond = tnv_conv
			}
			
			if IsExactType[IntTypeAndVal](tnv_cond) && tnv_cond.(IntTypeAndVal).Value != 0 {
				counter++
				if counter >= inf_protect {
					interp.MsgSpan.PrepNote(ast.Span(), "here\n")
					fmt.Fprintf(os.Stdout, interp.DoMessage(s, "runtime error", COLOR_RED, "infinite loop (counter went over 1M iterations) detected."))
					// throw error
					return VoidTypeAndVal{}
				}
				goto do_while
			}
		} else {
		while:
			tnv_cond := interp.EvalExpr(ast.Cond)
			if tnv_conv, res := ConvertToInt(tnv_cond); res {
				tnv_cond = tnv_conv
			}
			
			if IsExactType[IntTypeAndVal](tnv_cond) && tnv_cond.(IntTypeAndVal).Value != 0 {
				tnv_body := interp.EvalStmt(ast.Body, &blk_flow)
				switch blk_flow {
				case FLOW_CNT:
					goto while
				case FLOW_RET, FLOW_BRK:
					if blk_flow==FLOW_RET {
						*flow = blk_flow
					}
					return tnv_body
				}
				
				counter++
				if counter >= inf_protect {
					interp.MsgSpan.PrepNote(ast.Span(), "here\n")
					fmt.Fprintf(os.Stdout, interp.DoMessage(s, "runtime error", COLOR_RED, "infinite loop (counter went over 1M iterations) detected."))
					// throw error
					return VoidTypeAndVal{}
				}
				goto while
			}
		}
	case *IfStmt:
		tnv_cond := interp.EvalExpr(ast.Cond)
		if tnv_conv, res := ConvertToInt(tnv_cond); res {
			tnv_cond = tnv_conv
		}
		
		blk_flow := FLOW_EXC
		if IsExactType[IntTypeAndVal](tnv_cond) && tnv_cond.(IntTypeAndVal).Value != 0 {
			tnv_then := interp.EvalStmt(ast.Then, &blk_flow)
			if blk_flow==FLOW_RET {
				*flow = blk_flow
				return tnv_then
			}
		} else if ast.Else != nil {
			tnv_then := interp.EvalStmt(ast.Else, &blk_flow)
			if blk_flow==FLOW_RET {
				*flow = blk_flow
				return tnv_then
			}
		}
	case *FlowStmt:
		switch ast.Kind {
		case TKContinue:
			*flow = FLOW_CNT
		case TKBreak:
			*flow = FLOW_BRK
		}
	case *RetStmt:
		*flow = FLOW_RET
		return interp.EvalExpr(ast.X)
	case *ExprStmt:
		interp.EvalExpr(ast.X)
	case *BadStmt:
		return VoidTypeAndVal{}
	}
	return VoidTypeAndVal{}
}