package SPTools

import (
	"os"
	"fmt"
	///"time"
)


type Type interface {
	aType()
}

type BaseType uint8
const (
	TYPE_VOID = BaseType(iota)
	TYPE_HANDLE
	TYPE_BOOL
	TYPE_INT
	TYPE_CHAR
	TYPE_FLOAT
	TYPE_ANY
)
func (BaseType) aType() {}

type RefType struct {
	Base BaseType
}
func (RefType) aType() {}

type ArrayType struct {
	ElemType Type
	Len int
	Dynamic bool
}
func (ArrayType) aType() {}

type FuncType struct {
	Params map[*Name]Type
	RetType Type
}
func (FuncType) aType() {}


func IsExactType[T any](a any) bool {
	_, is_type := a.(T)
	return is_type
}

func AreSameType[T any](a, b any) bool {
	_, a_type := a.(T)
	_, b_type := b.(T)
	return a_type && b_type
}

// for binary operations.
// please make sure they're both the same type by using 'AreSameType'.
func GetBinaryTypes[T any](a, b any) (T, T) {
	x := a.(T)
	y := b.(T)
	return x, y
}

func IsBaseType(a Type) bool {
	switch a.(type) {
	case BaseType:
		return true
	default:
		return false
	}
}

func IsBaseTypeOfType(a Type, bt ...BaseType) bool {
	switch t := a.(type) {
	case BaseType:
		for _, val := range bt {
			if t==val {
				return true
			}
		}
	}
	return false
}

func IsArithmeticType(a Type) bool {
	switch t := a.(type) {
	case BaseType:
		return t >= TYPE_INT
	default:
		return false
	}
}

func GetTypeName(a Type) string {
	switch t := a.(type) {
	case BaseType:
		switch t {
		case TYPE_VOID:
			return "'void' type"
		case TYPE_ANY:
			return "'any' type"
		case TYPE_BOOL:
			return "'bool' type"
		case TYPE_CHAR:
			return "'char' type"
		case TYPE_FLOAT:
			return "'float' type"
		case TYPE_INT:
			return "'int' type"
		case TYPE_HANDLE:
			return "'Handle' type"
		}
	case RefType:
		return "Reference of " + GetTypeName(t.Base)
	case ArrayType:
		return "Array Type of element " + GetTypeName(t.ElemType)
	case FuncType:
		return "Function Type"
	}
	return "Unknown Type"
}


/*
 * 1. Two expressions are convertible when their reduced forms are the same. E.g 2 + 2 is convertible to 4
 * 2. Two expressions are coercible when you can safely cast one to the other. E.g 22 : int32 might be coercible to 22 : int64
 */
func AreTypesCoercible(a, b Type) bool {
	if IsArithmeticType(a) && IsArithmeticType(b) {
		return true
	} else if IsBaseType(a) && IsBaseType(b) {
		// probably bool type. Allow it to coerce.
		tA, tB := GetBinaryTypes[BaseType](a, b)
		return tA <= tB || tA >= tB
	} else if IsExactType[ArrayType](a) && IsExactType[ArrayType](b) {
		if tA, tB := GetBinaryTypes[ArrayType](a, b); tA.Len != tB.Len || !AreTypesCoercible(tA.ElemType, tB.ElemType) {
			// trying to coerce array with different size or arent same type.
			return false
		}
		return true
	}
	return false
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


type TypeChecker struct {
	MsgSpan
	Syms, Types map[string]Type
}

func MakeTypeChecker(p Parser) TypeChecker {
	var tc = TypeChecker{ MsgSpan: p.TokenReader.MsgSpan, Types: make(map[string]Type) }
	tc.Types["int"] = TYPE_INT
	tc.Types["any"] = TYPE_ANY
	tc.Types["bool"] = TYPE_BOOL
	tc.Types["char"] = TYPE_CHAR
	tc.Types["float"] = TYPE_FLOAT
	return tc
}

func (c TypeChecker) DoMessage(n Node, msgtype, color, msg string, args ...any) string {
	t := n.Tok()
	report := c.MsgSpan.Report(msgtype, "", color, msg, *t.Path, &t.Span.LineStart, &t.Span.ColStart, args...)
	c.MsgSpan.PurgeNotes()
	return report
}


func (c *TypeChecker) CheckExpr(e Expr) {
	if e==nil {
		return
	}
	
	switch ast := e.(type) {
	case *BasicLit:
		switch ast.Kind {
		case IntLit, BoolLit, CharLit:
			ast.tag = TYPE_INT
		case FloatLit:
			ast.tag = TYPE_FLOAT
		case StringLit:
			ast.tag = ArrayType{ ElemType: TYPE_CHAR, Len: len(ast.node.Tok().Lexeme) }
		}
	case *BracketExpr:
		for i := range ast.Exprs {
			c.CheckExpr(ast.Exprs[i])
		}
		ast.tag = ArrayType{ ElemType: nil, Len: len(ast.Exprs) }
	case *NullExpr:
		ast.tag = TYPE_HANDLE
	case *ThisExpr: // return type of 'this'.
		ast.tag = TYPE_INT
	case *Name:
		ast.tag = c.Types[ast.Value]
	case *UnaryExpr:
		switch ast.Kind {
		case TKIncr, TKDecr:
			// TODO: syms/lvalue here. Check if lvalue is float or int type.
			ast.tag = TYPE_INT
		case TKNot, TKCompl, TKSub, TKSizeof:
			c.CheckExpr(ast.X)
			if t := ast.X.Tag(); IsBaseTypeOfType(t, TYPE_INT, TYPE_CHAR, TYPE_ANY, TYPE_BOOL) {
				ast.tag = t
			} else {
				// error
				c.MsgSpan.PrepNote(ast.Span(), "here\n")
				var expr_type string
				switch ast.Kind {
				case TKNot:
					expr_type = "Logical NOT"
				case TKCompl:
					expr_type = "Bitwise NOT/Complement"
				case TKSub:
					expr_type = "Negation"
				}
				fmt.Fprintf(os.Stdout, c.DoMessage(e, "type error", COLOR_RED, "Non-Int type for " + expr_type + " expression."))
			}
		case TKNew:
			// TODO: error on non-objects.
			// ast.X could be a function call expr or type/name index expr.
		}
	case *IndexExpr:
		// TODO: syms here.
		c.CheckExpr(ast.X)
		c.CheckExpr(ast.Index)
		typeOfX := ast.X.Tag()
		if !IsExactType[ArrayType](typeOfX) {
			///interp.MsgSpan.PrepNote(ast.Span(), "")
			fmt.Fprintf(os.Stdout, c.DoMessage(e, "runtime error", COLOR_RED, "Attempting to index non-Array type."))
			// error
		}
		
		typeOfIdx := ast.Index.Tag()
		if !IsBaseTypeOfType(typeOfIdx, TYPE_INT, TYPE_CHAR, TYPE_BOOL) {
			///interp.MsgSpan.PrepNote(ast.Span(), "here")
			fmt.Fprintf(os.Stdout, c.DoMessage(e, "runtime error", COLOR_RED, "Attempting to index Array type with non-Int value."))
			// error
		}
		ast.tag = typeOfX.(ArrayType).ElemType
	case *ViewAsExpr:
		// check coercion here.
		type_tok := ast.Type.(*TypedExpr)
		targetType := c.Types[type_tok.TypeName.Lexeme]
		c.CheckExpr(ast.X)
		typeOfX := ast.X.Tag()
		if AreTypesCoercible(targetType, typeOfX) {
			ast.tag = targetType
		} else {
			c.MsgSpan.PrepNote(ast.Span(), "here\n")
			fmt.Fprintf(os.Stdout, c.DoMessage(e, "type error", COLOR_RED, "Cannot coerce %s to %s.", GetTypeName(typeOfX), GetTypeName(targetType)))
			// error, trying to convert to invalid type.
		}
	case *BinExpr:
		c.CheckExpr(ast.L)
		c.CheckExpr(ast.R)
		l, r := ast.L.Tag(), ast.R.Tag()
		var binary_type Type

		if ast.Kind >= TKAdd && ast.Kind <= TKOrL {
			// if mixing with char type, promote it to int.
			// if 'any' type, autocast to int.
			if IsBaseTypeOfType(l, TYPE_CHAR) || IsBaseTypeOfType(r, TYPE_CHAR) {
				// sign extend to int.
				binary_type = TYPE_INT
			}
			
			// if mixing with float type, entire expr is float type.
			if IsBaseTypeOfType(l, TYPE_FLOAT) || IsBaseTypeOfType(r, TYPE_FLOAT) {
				binary_type = TYPE_FLOAT
			}
		}
		
		switch ast.Kind {
		case TKMod, TKAnd, TKAndNot, TKOr, TKXor, TKShAL, TKShAR, TKShLR:
			if !IsBaseTypeOfType(l, TYPE_INT, TYPE_CHAR, TYPE_ANY) || !IsBaseTypeOfType(r, TYPE_INT, TYPE_CHAR, TYPE_ANY) {
				// illegal operation for non-int types.
				c.MsgSpan.PrepNote(ast.Span(), "here\n")
				fmt.Fprintf(os.Stdout, c.DoMessage(e, "type error", COLOR_RED, "Cannot do int operation with non-int types."))
				return
			}
			fallthrough
		case TKAdd, TKSub, TKMul, TKDiv, TKNotEq, TKEq, TKAndL, TKOrL:
			ast.tag = binary_type
		
		// assignments.
		// "write" value into lvalue.
		// make sure lvalue is of compatible typing or coercable.
		case TKAssign, TKAddA, TKSubA, TKMulA, TKDivA, TKModA, TKAndA, TKAndNotA, TKOrA, TKXorA, TKShALA, TKShARA, TKShLRA:
			// we can't promote L to higher type here.
			if !AreTypesCoercible(l, r) {
				// illegal operation for non-int types.
				c.MsgSpan.PrepNote(ast.Span(), "here\n")
				fmt.Fprintf(os.Stdout, c.DoMessage(e, "type error", COLOR_RED, "Cannot coerce type %s to %s in assignment.", GetTypeName(l), GetTypeName(r)))
				return
			}
			ast.tag = l
		}
	case *ChainExpr:
		// a # b # c => a # b && b # c
		// a # b ==> a = b; b = next; a # b; repeat.
		c.CheckExpr(ast.A)
		a := ast.A.Tag()
		var binary_type Type
		if IsBaseTypeOfType(a, TYPE_CHAR) {
			binary_type = TYPE_INT
		}
		
		if IsBaseTypeOfType(a, TYPE_FLOAT) {
			binary_type = TYPE_FLOAT
		}
		
		for i := range ast.Kinds {
			c.CheckExpr(ast.Bs[i])
			b := ast.Bs[i].Tag()
			switch ast.Kinds[i] {
			case TKLess, TKGreater, TKGreaterE, TKLessE:
				ast.tag = TYPE_BOOL
			}
			a = b
		}
		
		if ast.tag==nil {
			ast.tag = binary_type
		}
	case *TernaryExpr:
		c.CheckExpr(ast.A)
		c.CheckExpr(ast.B)
		c.CheckExpr(ast.C)
		if !AreTypesCoercible(ast.A.Tag(), ast.B.Tag()) {
			// error here.
		}
	case *CommaExpr:
		for i := range ast.Exprs {
			c.CheckExpr(ast.Exprs[i])
			ast.tag = ast.Exprs[i].Tag()
		}
	case *FuncLit:
		// TODO: implement function creation here.
	}
}

func (c *TypeChecker) CheckStmt(s Stmt) {
	if s==nil {
		return
	}
	
	switch ast := s.(type) {
	case *BlockStmt:
		// new scope here.
		for i := range ast.Stmts {
			c.CheckStmt(ast.Stmts[i])
		}
	case *WhileStmt:
	case *IfStmt:
	case *FlowStmt:
	case *RetStmt:
		if ast.X != nil {
			c.CheckExpr(ast.X)
		}
	case *ExprStmt:
		c.CheckExpr(ast.X)
	case *BadStmt:
	}
}