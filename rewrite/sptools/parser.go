package SPTools

import (
	"os"
	"fmt"
	///"time"
)


const ERR_LIMIT = 5

type Parser struct {
	*TokenReader
	Errs []string
}

func (parser *Parser) GetToken(offset int) Token {
	return parser.TokenReader.Get(offset, TOKFLAG_IGNORE_ALL)
}

func (parser *Parser) Advance(i int) Token {
	parser.TokenReader.Advance(i)
	return parser.GetToken(0)
}

func (parser *Parser) HasTokenKindSeq(kinds ...TokenKind) bool {
	return parser.TokenReader.HasTokenKindSeq(TOKFLAG_IGNORE_ALL, kinds...)
}

func (parser *Parser) ReportErrs() bool {
	if len(parser.Errs)==0 {
		t := parser.GetToken(0)
		report := parser.MsgSpan.Report("success", "", COLOR_GREEN, "successfully parsed.", *t.Path, nil, nil)
		SpewReport(os.Stdout, report, nil)
		return true
	} else {
		for _, err := range parser.Errs {
			fmt.Fprintf(os.Stdout, "%s\n", err)
		}
		return false
	}
}

// REMEMBER, this auto-increments the token index if it matches.
// so DO NOT increment the token index after using this.
func (parser *Parser) got(tk TokenKind) bool {
	if token := parser.GetToken(0); token.Kind==tk {
		parser.Advance(1)
		return true
	}
	return false
}

func (parser *Parser) syntaxErr(msg string, args ...any) {
	token := parser.GetToken(-1)
	report := parser.MsgSpan.Report("syntax error", "", COLOR_RED, msg, *token.Path, &token.Span.LineStart, &token.Span.ColStart, args...)
	parser.MsgSpan.PurgeNotes()
	if len(parser.Errs) <= ERR_LIMIT {
		parser.Errs = append(parser.Errs, report)
		if len(parser.Errs) > ERR_LIMIT {
			parser.Errs = append(parser.Errs, parser.MsgSpan.Report("system", "", COLOR_RED, "too many errors.", *token.Path, nil, nil))
		}
	}
}

// TODO: allow newlines in place of semicolons.
func (parser *Parser) want(tk TokenKind, lexeme string) bool {
	if !parser.got(tk) {
		t := parser.GetToken(0)
		parser.MsgSpan.PrepNote(t.Span, "")
		parser.syntaxErr("expecting '%s' but got '%s'", lexeme, parser.GetToken(0).Lexeme)
		// continue on and try to parse the remainder
		parser.Advance(1)
		return false
	}
	return true
}


func (parser *Parser) Start() Node {
	if parser.TokenReader.Len() <= 0 {
		report := parser.MsgSpan.Report("parsing error", "", COLOR_RED, "Token buffer is EMPTY!", "", nil, nil)
		SpewReport(os.Stdout, report, nil)
		parser.MsgSpan.PurgeNotes()
		return nil
	}
	return parser.TopDecl()
}

// Plugin = +TopDecl .
// TopDecl = FuncDecl | TypeDecl | VarDecl | StaticAssertion .
func (parser *Parser) TopDecl() Node {
	///defer fmt.Printf("parser.TopDecl()\n")
	plugin := new(Plugin)
	for t := parser.GetToken(0); t.Kind != TKEoF; t = parser.GetToken(0) {
		///time.Sleep(100 * time.Millisecond)
		///fmt.Printf("TopDecl :: current tok: %v\n", t)
		if t.IsStorageClass() || t.IsType() || t.Kind==TKIdent && parser.GetToken(1).Kind==TKIdent {
			///fmt.Printf("TopDecl :: func or var decl: %v\n", t)
			v_or_f_decl := parser.DoVarOrFuncDecl(false)
			if vdecl, is_var_decl := v_or_f_decl.(*VarDecl); is_var_decl {
				if parser.GetToken(-1).Kind==TKRCurl {
					if parser.GetToken(0).Kind==TKSemi {
						parser.Advance(1)
					}
				} else if !parser.got(TKSemi) {
					parser.syntaxErr("missing ';' semicolon for global variable:")
					for i := range vdecl.Names {
						PrintNode(vdecl.Names[i], 1, os.Stdout)
					}
					bad := new(BadDecl)
					copyPosToNode(&bad.node, t)
					plugin.Decls = append(plugin.Decls, bad)
					goto err_exit
				}
			}
			plugin.Decls = append(plugin.Decls, v_or_f_decl)
		} else if t.Kind==TKStaticAssert {
			stasrt := new(StaticAssert)
			copyPosToNode(&stasrt.node, t)
			parser.Advance(1)
			parser.want(TKLParen, "(")
			stasrt.A = parser.MainExpr()
			if parser.GetToken(0).Kind==TKComma {
				parser.Advance(1)
				stasrt.B = parser.MainExpr()
			}
			parser.want(TKRParen, ")")
			if !parser.got(TKSemi) {
				parser.noSemi()
				goto err_exit
			}
			plugin.Decls = append(plugin.Decls, stasrt)
		} else {
			///fmt.Printf("TopDecl :: type decl: %v\n", t)
			type_decl := new(TypeDecl)
			copyPosToNode(&type_decl.node, t)
			switch t.Kind {
			case TKMethodMap:
				type_decl.Type = parser.DoMethodMap()
			case TKTypedef:
				type_decl.Type = parser.DoTypedef()
			case TKTypeset:
				type_decl.Type = parser.DoTypeSet()
			case TKEnum:
				type_decl.Type = parser.DoEnumSpec()
			case TKStruct:
				type_decl.Type = parser.DoStruct(false)
			case TKUsing:
				type_decl.Type = parser.DoUsingSpec()
			default:
				parser.MsgSpan.PrepNote(t.Span, "")
				parser.syntaxErr("bad declaration: %q", t.String())
				bad := new(BadDecl)
				copyPosToNode(&bad.node, t)
				plugin.Decls = append(plugin.Decls, bad)
				goto err_exit
			}
			plugin.Decls = append(plugin.Decls, type_decl)
		}
	}
err_exit:
	parser.ReportErrs()
	return plugin
}


// VarDecl  = VarOrFuncSpec VarDeclarator .
// FuncDecl = VarOrFuncSpec FuncDeclarator .
func (parser *Parser) DoVarOrFuncDecl(param bool) Decl {
	///defer fmt.Printf("parser.DoVarOrFuncDecl()\n")
	saved_token := parser.GetToken(0)
	class_flags, spec_type := parser.VarOrFuncSpec()
	ident := parser.PrimaryExpr() // get NAME only.
	if t := parser.GetToken(0); t.Kind==TKLParen {
		fdecl := new(FuncDecl)
		copyPosToNode(&fdecl.node, saved_token)
		fdecl.RetType = spec_type
		fdecl.ClassFlags = class_flags
		fdecl.Ident = ident
		parser.DoFuncDeclarator(fdecl)
		return fdecl
	} else {
		vdecl := new(VarDecl)
		copyPosToNode(&vdecl.node, saved_token)
		vdecl.Type = spec_type
		vdecl.ClassFlags = class_flags
		vdecl.Names = append(vdecl.Names, ident)
		parser.DoVarDeclarator(vdecl, param)
		return vdecl
	}
}

// VarDeclarator = Ident [ IndexExpr ] [ Initializer ] *( ',' VarDeclarator ) .
// Initializer = '=' SubMainExpr | '{' Expr [ ',' ( '...' | *Expr ) ] '}' .
func (parser *Parser) DoVarDeclarator(vdecl *VarDecl, param bool) {
	///defer fmt.Printf("parser.DoVarDeclarator()\n")
	// This is structured as if it's a do-while loop.
	for {
		///time.Sleep(100 * time.Millisecond)
		///fmt.Printf("parser.DoVarDeclarator() - '%+v'\n", parser.GetToken(0).ToString())
		if parser.GetToken(0).Kind==TKLBrack {
			var dims []Expr
			for t := parser.GetToken(0); t.Kind != TKEoF && t.Kind==TKLBrack; t = parser.GetToken(0) {
				parser.want(TKLBrack, "[")
				if parser.GetToken(0).Kind != TKRBrack {
					dims = append(dims, parser.SubMainExpr())
				}
				parser.want(TKRBrack, "]")
			}
			vdecl.Dims = append(vdecl.Dims, dims)
		} else {
			vdecl.Dims = append(vdecl.Dims, nil)
		}
		
		if parser.GetToken(0).Kind==TKAssign {
			parser.Advance(1)
			if parser.GetToken(0).Kind==TKLCurl {
				vdecl.Inits = append(vdecl.Inits, parser.PrimaryExpr())
			} else {
				vdecl.Inits = append(vdecl.Inits, parser.SubMainExpr())
			}
		} else {
			vdecl.Inits = append(vdecl.Inits, nil)
		}
		
		if ending := parser.GetToken(0); param || ending.Kind==TKEoF || ending.Kind==TKSemi {
			break
		} else if ending.Kind==TKComma {
			parser.Advance(1)
		} else {
			break
		}
		
		ident := parser.PrimaryExpr()
		vdecl.Names = append(vdecl.Names, ident)
	}
	///fmt.Printf("Leaving parser.DoVarDeclarator(): '%+v'\n", parser.GetToken(0).ToString())
}

// ParamList = '(' *VarDecl ')' .
func (parser *Parser) DoParamList() []Decl {
	///defer fmt.Printf("parser.DoParamList()\n")
	var params []Decl
	parser.want(TKLParen, "(")
	for t := parser.GetToken(0); t.Kind != TKEoF && t.Kind != TKRParen; t = parser.GetToken(0) {
		///time.Sleep(100 * time.Millisecond)
		if len(params) > 0 {
			parser.want(TKComma, ",")
		}
		
		var_decl := parser.DoVarOrFuncDecl(true)
		if _, is_var_decl := var_decl.(*VarDecl); !is_var_decl {
			loc := var_decl.Span()
			parser.MsgSpan.PrepNote(loc, "this declaration here.")
			parser.syntaxErr("bad declaration: expected var declaration.")
			bad_decl := new(BadDecl)
			copyPosToNode(&bad_decl.node, t)
			params = append(params, bad_decl)
		} else {
			params = append(params, var_decl)
		}
	}
	parser.want(TKRParen, ")")
	return params
}

// FuncSpec = Ident ParamList ( Initializer | Block | ';' ) .
func (parser *Parser) DoFuncDeclarator(fdecl *FuncDecl) {
	///defer fmt.Printf("parser.DoFuncDeclarator()\n")
	fdecl.Params = parser.DoParamList()
	switch t := parser.GetToken(0); t.Kind {
	case TKSemi:
		parser.want(TKSemi, ";")
		fdecl.Body = nil
	case TKLCurl:
		fdecl.Body = parser.DoBlock()
	case TKAssign:
		parser.Advance(1)
		fdecl.Body = parser.MainExpr()
		if parser.GetToken(0).Kind==TKSemi {
			parser.Advance(1)
		}
	default:
		name := fdecl.Ident.Tok()
		parser.MsgSpan.PrepNote(name.Span, "this function here.")
		parser.MsgSpan.PrepNote(t.Span, "needs placement here.")
		parser.syntaxErr("bad declaration: expected body, assignment, or semicolon for function.")
		fdecl.Body = new(BadStmt)
		copyPosToNode(&fdecl.node, t)
	}
}


// StorageClass = 'native' | 'forward' | 'const' | 'static' | 'stock' | 'public' | 'private' | 'protected' | 'readonly' | 'sealed' | 'virtual' .
func (parser *Parser) StorageClass() StorageClassFlags {
	///defer fmt.Printf("parser.StorageClass()\n")
	flags := StorageClassFlags(0)
	for parser.GetToken(0).IsStorageClass() {
		flags |= storageClassFromToken(parser.GetToken(0))
		parser.Advance(1)
	}
	return flags
}

// AbstractDecl = Type [ *'[]' | '&' ] .
func (parser *Parser) AbstractDecl() Spec {
	///defer fmt.Printf("parser.AbstractDecl()\n")
	tspec := new(TypeSpec)
	copyPosToNode(&tspec.node, parser.GetToken(0))
	// next get type name.
	tspec.Type = parser.TypeExpr(false)
	
	// check pre-identifier array dims or ampersand reference.
	switch t := parser.GetToken(0); t.Kind {
	case TKLBrack:
		for parser.GetToken(0).Kind==TKLBrack {
			tspec.Dims++
			parser.want(TKLBrack, "[")
			parser.want(TKRBrack, "]")
			tspec.IsRef = true
		}
	case TKAnd:
		tspec.IsRef = true
		parser.Advance(1)
	}
	return tspec
}

// VarOrFuncSpec = *StorageClass AbstractDecl .
func (parser *Parser) VarOrFuncSpec() (StorageClassFlags, Spec) {
	///defer fmt.Printf("parser.VarOrFuncSpec()\n")
	return parser.StorageClass(), parser.AbstractDecl()
}

// SignatureSpec = 'function' AbstractDecl ParamsList .
func (parser *Parser) DoFuncSignature() Spec {
	///defer fmt.Printf("parser.DoFuncSignature()\n")
	sig := new(SignatureSpec)
	copyPosToNode(&sig.node, parser.GetToken(0))
	parser.want(TKFunction, "function")
	sig.Type = parser.AbstractDecl()
	sig.Params = parser.DoParamList()
	return sig
}

// EnumSpec = 'enum' [ ident [ ':' ] '(' operator PrimaryExpr ')' ] '{' +EnumEntry '}' [ ';' ] .
// EnumEntry = Ident [ '=' Expr ] .
func (parser *Parser) DoEnumSpec() Spec {
	///defer fmt.Printf("parser.DoEnumSpec()\n")
	enum := new(EnumSpec)
	copyPosToNode(&enum.node, parser.GetToken(0))
	parser.want(TKEnum, "enum")
	if t := parser.GetToken(0); t.Kind==TKStruct {
		return parser.DoStruct(true)
	}
	
	if parser.GetToken(0).Kind==TKIdent {
		enum.Ident = parser.PrimaryExpr()
	}
	
	if parser.GetToken(0).Kind==TKColon {
		parser.want(TKColon, ":")
	}
	
	if t := parser.GetToken(0); t.Kind==TKLParen {
		parser.Advance(1)
		t = parser.GetToken(0)
		if !t.IsOperator() {
			parser.MsgSpan.PrepNote(enum.Span(), "this enum here.\n")
			parser.MsgSpan.PrepNote(t.Span, "placement here.")
			parser.syntaxErr("expected math operator for enum auto-incrementer.")
			bad := new(BadSpec)
			copyPosToNode(&bad.node, parser.GetToken(0))
			return bad
		} else {
			enum.StepOp = t.Kind
		}
		parser.Advance(1)
		enum.Step = parser.SubMainExpr()
		parser.want(TKRParen, ")")
	}
	
	parser.want(TKLCurl, "{")
	for {
		///time.Sleep(100 * time.Millisecond)
		if parser.GetToken(0).Kind==TKRCurl {
			break
		}
		
		///fmt.Printf("DoEnumSpec :: current tok: %v\n", parser.GetToken(0))
		enum.Names = append(enum.Names, parser.PrimaryExpr())
		if parser.GetToken(0).Kind==TKAssign {
			parser.Advance(1)
			enum.Values = append(enum.Values, parser.SubMainExpr())
		} else {
			enum.Values = append(enum.Values, nil)
		}
		
		if parser.GetToken(0).Kind==TKComma {
			parser.Advance(1)
		} else {
			break
		}
	}
	parser.want(TKRCurl, "}")
	if parser.GetToken(0).Kind==TKSemi {
		parser.Advance(1)
	}
	return enum
}

// StructSpec = 'struct' Ident '{' *Field '}' [ ';' ] .
// Field = VarDecl ';' | FuncDecl .
func (parser *Parser) DoStruct(is_enum bool) Spec {
	///defer fmt.Printf("parser.DoStruct()\n")
	struc := new(StructSpec)
	copyPosToNode(&struc.node, parser.GetToken(0))
	parser.want(TKStruct, "struct")
	struc.IsEnum = is_enum
	struc.Ident = parser.PrimaryExpr()
	parser.want(TKLCurl, "{")
	for t := parser.GetToken(0); t.Kind != TKEoF && t.Kind != TKRCurl; t = parser.GetToken(0) {
		///time.Sleep(100 * time.Millisecond)
		///fmt.Printf("DoStruct :: current tok: %v\n", t)
		v_or_f_decl := parser.DoVarOrFuncDecl(false)
		switch ast := v_or_f_decl.(type) {
		case *VarDecl:
			struc.Fields = append(struc.Fields, ast)
			parser.want(TKSemi, ";")
		case *FuncDecl:
			struc.Methods = append(struc.Methods, ast)
		default:
			name := struc.Ident.Tok()
			parser.MsgSpan.PrepNote(name.Span, "this struct here.\n")
			a := v_or_f_decl.Tok()
			parser.MsgSpan.PrepNote(a.Span, "illegal construct here.")
			if is_enum {
				parser.syntaxErr("bad field/method in enum struct")
			} else {
				parser.syntaxErr("bad field/method in struct")
			}
			bad := new(BadSpec)
			copyPosToNode(&bad.node, parser.GetToken(0))
			return bad
		}
	}
	parser.want(TKRCurl, "}")
	if parser.GetToken(0).Kind==TKSemi {
		parser.Advance(1)
	}
	return struc
}

// UsingSpec = 'using' Expr ';' .
func (parser *Parser) DoUsingSpec() Spec {
	///defer fmt.Printf("parser.DoUsingSpec()\n")
	using := new(UsingSpec)
	copyPosToNode(&using.node, parser.GetToken(0))
	parser.want(TKUsing, "using")
	using.Namespace = parser.SubMainExpr()
	if !parser.got(TKSemi) {
		end := parser.GetToken(-1)
		parser.MsgSpan.PrepNote(end.Span, "missing ';' here.")
		parser.syntaxErr("missing ending ';' semicolon for 'using' specification.")
		bad := new(BadSpec)
		copyPosToNode(&bad.node, parser.GetToken(-1))
		return bad
	}
	return using
}

// TypeSetSpec = 'typeset' Ident '{' *( SignatureSpec ';' ) '}' [ ';' ] .
func (parser *Parser) DoTypeSet() Spec {
	///defer fmt.Printf("parser.DoTypeSet()\n")
	typeset := new(TypeSetSpec)
	copyPosToNode(&typeset.node, parser.GetToken(0))
	parser.want(TKTypeset, "typeset")
	typeset.Ident = parser.PrimaryExpr()
	parser.want(TKLCurl, "{")
	for t := parser.GetToken(0); t.Kind != TKEoF && t.Kind != TKRCurl; t = parser.GetToken(0) {
		///time.Sleep(100 * time.Millisecond)
		///fmt.Printf("DoTypeSet :: current tok: %v\n", t)
		signature := parser.DoFuncSignature()
		typeset.Signatures = append(typeset.Signatures, signature)
		parser.want(TKSemi, ";")
	}
	parser.want(TKRCurl, "}")
	if parser.GetToken(0).Kind==TKSemi {
		parser.Advance(1)
	}
	return typeset
}

// TypeDefSpec = 'typedef' Ident '=' SignatureSpec ';' .
func (parser *Parser) DoTypedef() Spec {
	///defer fmt.Printf("parser.DoTypedef()\n")
	typedef := new(TypeDefSpec)
	copyPosToNode(&typedef.node, parser.GetToken(0))
	parser.want(TKTypedef, "typedef")
	typedef.Ident = parser.PrimaryExpr()
	parser.want(TKAssign, "=")
	typedef.Sig = parser.DoFuncSignature()
	if !parser.got(TKSemi) {
		name := typedef.Ident.Tok()
		parser.MsgSpan.PrepNote(name.Span, "for this typedef here.\n")
		end := parser.GetToken(-1)
		parser.MsgSpan.PrepNote(end.Span, "missing ';' here.")
		parser.syntaxErr("missing ending ';' semicolon for 'typedef' specification.")
		bad := new(BadSpec)
		copyPosToNode(&bad.node, end)
		return bad
	}
	return typedef
}


// MethodMapSpec = 'methodmap' Ident [ '__nullable__' ] [ '<' TypeExpr ] '{' [ MethodCtor ] *MethodMapEntry '}' [ ';' ] .
// MethodCtor = 'public' Ident ParamList ( Block | ';' ) .
// MethodMapEntry = MethodMapProp | FuncDecl .
// MethodMapProp = 'property' TypeExpr Ident '{' PropGetter [ PropSetter ] | PropSetter '}' .
// PropGetter = 'public' [ 'native' ] 'get' '(' ')' ( Block | ';' ) .
// PropSetter = 'public' [ 'native' ] 'set' ParamList ( Block | ';' ) .
func (parser *Parser) DoMethodMap() Spec {
	///defer fmt.Printf("parser.DoMethodMap()\n")
	methodmap := new(MethodMapSpec)
	copyPosToNode(&methodmap.node, parser.GetToken(0))
	parser.want(TKMethodMap, "methodmap")
	methodmap.Ident = parser.PrimaryExpr()
	if parser.GetToken(0).Kind==TKNullable {
		parser.Advance(1)
		methodmap.Nullable = true
	}
	
	if parser.GetToken(0).Kind==TKLess {
		parser.Advance(1)
		methodmap.Parent = parser.TypeExpr(false)
	}
	
	parser.want(TKLCurl, "{")
	for t := parser.GetToken(0); t.Kind != TKEoF && (t.Kind==TKPublic || t.Kind==TKProperty); t = parser.GetToken(0) {
		switch t.Kind {
		case TKPublic:
			// gotta use lookahead for this...
			// if after the 'public' or 'native' keyword is an identifier
			// and after identifier is a left parenthesis, it's likely the constructor.
			//t1, t2, t3 := parser.GetToken(1), parser.GetToken(2), parser.GetToken(3)
			//if (t1.Kind==TKIdent && t2.Kind==TKLParen) || (t1.Kind==TKNative && t2.Kind==TKIdent && t3.Kind==TKLParen) {
			if parser.HasTokenKindSeq(TKPublic, TKIdent, TKLParen) || parser.HasTokenKindSeq(TKPublic, TKNative, TKIdent, TKLParen) {
				// kinda have to make a *FuncDecl from scratch here...
				ctor_decl := new(FuncDecl)
				//ctor_decl.RetType = methodmap
				copyPosToNode(&ctor_decl.node, t)
				// eats up the 'public' and 'native' keyword if it's there.
				ctor_decl.ClassFlags = parser.StorageClass()
				ctor_decl.Ident = parser.PrimaryExpr()
				parser.DoFuncDeclarator(ctor_decl)
				method := new(MethodMapMethodSpec)
				copyPosToNode(&method.node, parser.GetToken(0))
				method.Impl = ctor_decl
				method.IsCtor = true
				methodmap.Methods = append(methodmap.Methods, method)
			} else {
				method := new(MethodMapMethodSpec)
				copyPosToNode(&method.node, t)
				method.Impl = parser.DoVarOrFuncDecl(false)
				methodmap.Methods = append(methodmap.Methods, method)
			}
		case TKProperty:
			parser.Advance(1)
			prop := new(MethodMapPropSpec)
			prop.Type = parser.TypeExpr(false)
			prop.Ident = parser.PrimaryExpr()
			parser.want(TKLCurl, "{")
			if spec_ret := parser.DoMethodMapProperty(prop); spec_ret != nil {
				methodmap.Props = append(methodmap.Props, spec_ret)
				goto methodmap_loop_exit
			}
			if parser.GetToken(0).Kind==TKPublic {
				if spec_ret := parser.DoMethodMapProperty(prop); spec_ret != nil {
					methodmap.Props = append(methodmap.Props, spec_ret)
					goto methodmap_loop_exit
				}
			}
			parser.want(TKRCurl, "}")
			methodmap.Props = append(methodmap.Props, prop)
		}
	}
methodmap_loop_exit:
	parser.want(TKRCurl, "}")
	if parser.GetToken(0).Kind==TKSemi {
		parser.Advance(1)
	}
	return methodmap
}


func (parser *Parser) DoMethodMapProperty(prop *MethodMapPropSpec) *BadSpec {
	storage_cls := parser.StorageClass()
	if g := parser.GetToken(0); g.Lexeme=="get" {
		prop.GetterClass = storage_cls
		parser.Advance(1)
		parser.want(TKLParen, "(")
		parser.want(TKRParen, ")")
		if end := parser.GetToken(0); end.Kind==TKLCurl {
			prop.GetterBlock = parser.DoBlock()
		} else if end.Kind==TKSemi {
			parser.Advance(1)
		} else {
			name := prop.Ident.Tok()
			parser.MsgSpan.PrepNote(name.Span, "in property here.\n")
			parser.MsgSpan.PrepNote(g.Span, "offending 'get' starts here.\n")
			parser.MsgSpan.PrepNote(end.Span, "end of 'get' property here.")
			parser.syntaxErr("expected ending } or ; for get implementation on methodmap property.")
			bad := new(BadSpec)
			copyPosToNode(&bad.node, parser.GetToken(-1))
			return bad
		}
	} else if s := parser.GetToken(0); s.Lexeme=="set" {
		prop.SetterClass = storage_cls
		parser.Advance(1)
		prop.SetterParams = parser.DoParamList()
		if end := parser.GetToken(0); end.Kind==TKLCurl {
			prop.SetterBlock = parser.DoBlock()
		} else if end.Kind==TKSemi {
			parser.Advance(1)
		} else {
			name := prop.Ident.Tok()
			parser.MsgSpan.PrepNote(name.Span, "in property here.\n")
			parser.MsgSpan.PrepNote(s.Span, "offending 'set' starts here.\n")
			parser.MsgSpan.PrepNote(end.Span, "end of 'set' property here.")
			parser.syntaxErr("expected ending } or ; for set implementation on methodmap property.")
			bad := new(BadSpec)
			copyPosToNode(&bad.node, parser.GetToken(-1))
			return bad
		}
	}
	return nil
}


func (parser *Parser) noSemi() Stmt {
	t := parser.GetToken(-1)
	parser.MsgSpan.PrepNote(t.Span, "missing semicolon here")
	parser.syntaxErr("missing ';' semicolon, got %q.", t.String())
	bad := new(BadStmt)
	copyPosToNode(&bad.node, t)
	return bad
}

// BlockStmt = '{' *Statement '}' .
func (parser *Parser) DoBlock() Stmt {
	///defer fmt.Printf("parser.DoBlock()\n")
	block := new(BlockStmt)
	///fmt.Printf("starting tok: %v\n", parser.GetToken(0))
	parser.want(TKLCurl, "{")
	copyPosToNode(&block.node, parser.GetToken(-1))
	for t := parser.GetToken(0); t.Kind != TKRCurl && t.Kind != TKEoF; t = parser.GetToken(0) {
		///time.Sleep(100 * time.Millisecond)
		///fmt.Printf("current tok: %v\n", t)
		n := parser.Statement()
		if n==nil {
			///fmt.Printf("n == nil\n")
			continue
		}
		block.Stmts = append(block.Stmts, n)
	}
	parser.want(TKRCurl, "}")
	return block
}

/*
 * Statement = IfStmt | WhileStmt | ForStmt | SwitchStmt | BlockStmt |
 *             RetStmt | AssertStmt | StaticAssertStmt | DeclStmt | DeleteStmt | ExprStmt .
 */
func (parser *Parser) Statement() Stmt {
	///defer fmt.Printf("parser.Statement()\n")
	if tIsEoF := parser.GetToken(0); tIsEoF.Kind==TKEoF {
		bad := new(BadStmt)
		copyPosToNode(&bad.node, tIsEoF)
		return bad
	}
	switch t := parser.GetToken(0); t.Kind {
	case TKConst, TKPublic, TKPrivate, TKProtected, TKForward, TKNative, TKReadOnly, TKSealed, TKVirtual, TKStock:
		fallthrough
	case TKInt, TKInt8, TKInt16, TKInt32, TKInt64, TKIntN:
		fallthrough
	case TKUInt8, TKUInt16, TKUInt32, TKUInt64, TKChar, TKDouble, TKVoid, TKObject, TKDecl, TKStatic, TKVar:
		// parse declaration.
		vardecl := new(DeclStmt)
		copyPosToNode(&vardecl.node, t)
		vardecl.D = parser.DoVarOrFuncDecl(false)
		if !parser.got(TKSemi) {
			return parser.noSemi()
		}
		return vardecl
	case TKLCurl:
		return parser.DoBlock()
	case TKIf:
		return parser.DoIf()
	case TKDo, TKWhile:
		return parser.While(t)
	case TKFor:
		return parser.DoFor()
	///case TKForEach:
		// ForEachStmt = 'foreach' '(' DeclStmt 'in' Expr ')' Statement.
	case TKSwitch:
		return parser.Switch()
	case TKReturn:
		// RetStmt = 'return' [ Expr ] ';' .
		ret := new(RetStmt)
		copyPosToNode(&ret.node, t)
		parser.Advance(1)
		if parser.GetToken(0).Kind != TKSemi {
			ret.X = parser.MainExpr()
			if !parser.got(TKSemi) {
				return parser.noSemi()
			}
			return ret
		} else {
			parser.Advance(1)
			return ret
		}
	case TKStaticAssert:
		// StaticAssertStmt = 'static_assert' '(' Expr [ ',' Expr ] ')' ';' .
		stasrt := new(StaticAssert)
		copyPosToNode(&stasrt.node, t)
		parser.Advance(1)
		parser.want(TKLParen, "(")
		stasrt.A = parser.MainExpr()
		if parser.GetToken(0).Kind==TKComma {
			parser.Advance(1)
			stasrt.B = parser.MainExpr()
		}
		parser.want(TKRParen, ")")
		if !parser.got(TKSemi) {
			return parser.noSemi()
		}
		
		stasrtstmt := new(StaticAssertStmt)
		copyPosToNode(&stasrtstmt.node, t)
		stasrtstmt.A = stasrt
		return stasrtstmt
	case TKAssert:
		// AssertStmt = 'assert' Expr ';' .
		asrt := new(AssertStmt)
		copyPosToNode(&asrt.node, t)
		parser.Advance(1)
		asrt.X = parser.MainExpr()
		if !parser.got(TKSemi) {
			return parser.noSemi()
		}
		return asrt
	case TKDelete:
		// DeleteStmt = 'delete' Expr ';' .
		del := new(DeleteStmt)
		copyPosToNode(&del.node, t)
		parser.Advance(1)
		del.X = parser.MainExpr()
		if !parser.got(TKSemi) {
			return parser.noSemi()
		}
		return del
	case TKIdent, TKThis:
		// vast majority of the times,
		// an expression starts with an identifier.
		if t2 := parser.GetToken(1); t2.Kind==TKIdent || parser.HasTokenKindSeq(TKIdent, TKLBrack, TKRBrack) {
			// possible var decl with custom type.
			vardecl := new(DeclStmt)
			copyPosToNode(&vardecl.node, t)
			vardecl.D = parser.DoVarOrFuncDecl(false)
			if !parser.got(TKSemi) {
				return parser.noSemi()
			}
			return vardecl
		} else {
			exp := new(ExprStmt)
			copyPosToNode(&exp.node, t)
			exp.X = parser.MainExpr()
			if !parser.got(TKSemi) {
				return parser.noSemi()
			}
			return exp
		}
	case TKIncr, TKDecr, TKViewAs:
		exp := new(ExprStmt)
		copyPosToNode(&exp.node, t)
		exp.X = parser.MainExpr()
		if !parser.got(TKSemi) {
			return parser.noSemi()
		}
		return exp
	case TKBreak, TKContinue:
		flow := new(FlowStmt)
		copyPosToNode(&flow.node, t)
		flow.Kind = t.Kind
		parser.Advance(1)
		if !parser.got(TKSemi) {
			return parser.noSemi()
		}
		return flow
	case TKCase, TKDefault, TKSemi:
		// lone semicolon is considered an error.
		parser.Advance(1)
		fallthrough
	default:
		parser.MsgSpan.PrepNote(t.Span, "statement starts here.")
		parser.syntaxErr("unknown statement: %q", t.String())
		bad := new(BadStmt)
		copyPosToNode(&bad.node, parser.GetToken(0))
		parser.Advance(1)
		return bad
	}
}

// DoStmt = 'do' Statement 'while' '(' Expr ')' ';' .
// WhileStmt = 'while' '(' Expr ')' Statement .
func (parser *Parser) While(t Token) Stmt {
	///defer fmt.Printf("parser.While()\n")
	while := new(WhileStmt)
	copyPosToNode(&while.node, t)
	parser.Advance(1)
	if t.Kind==TKDo {
		// do-while
		while.Do = true
		while.Body = parser.Statement()
		parser.want(TKWhile, "while")
		parser.want(TKLParen, "(")
		while.Cond = parser.MainExpr()
		parser.want(TKRParen, ")")
		if !parser.got(TKSemi) {
			return parser.noSemi()
		}
	} else {
		// while
		parser.want(TKLParen, "(")
		while.Cond = parser.MainExpr()
		parser.want(TKRParen, ")")
		while.Body = parser.Statement()
	}
	return while
}

// IfStmt = 'if' '(' Expr ')' Statement [ 'else' Statement ] .
func (parser *Parser) DoIf() Stmt {
	///defer fmt.Printf("parser.DoIf()\n")
	parser.want(TKIf, "if")
	ifstmt := new(IfStmt)
	copyPosToNode(&ifstmt.node, parser.GetToken(-1))
	parser.want(TKLParen, "(")
	ifstmt.Cond = parser.MainExpr()
	parser.want(TKRParen, ")")
	ifstmt.Then = parser.Statement()
	if parser.GetToken(0).Kind==TKElse {
		parser.Advance(1)
		ifstmt.Else = parser.Statement()
		// personally, I'd prefer to fix the not-needing-{ thing but whatever.
		/*
		switch t := parser.GetToken(0); t.Kind {
		case TKIf:
			ifstmt.Else = parser.DoIf()
		case TKLCurl:
			ifstmt.Else = parser.DoBlock()
		default:
			parser.MsgSpan.PrepNote(Span, "reason here.")
			parser.syntaxErr("ill-formed else block, missing 'if' or { curl.")
			bad := new(BadStmt)
			copyPosToNode(&bad.node, t)
			ifstmt.Else = bad
		}
		*/
	}
	return ifstmt
}

// ForStmt = 'for' '(' [ Decl | Expr ] ';' [ Expr ] ';' [ Expr ] ')' Statement .
func (parser *Parser) DoFor() Stmt {
	///defer fmt.Printf("parser.DoFor()\n")
	parser.want(TKFor, "for")
	forstmt := new(ForStmt)
	copyPosToNode(&forstmt.node, parser.GetToken(-1))
	parser.want(TKLParen, "(")
	if parser.GetToken(0).Kind != TKSemi {
		if t := parser.GetToken(0); t.IsType() || t.IsStorageClass() || (t.Kind==TKIdent && parser.GetToken(1).Kind==TKIdent) {
			forstmt.Init = parser.DoVarOrFuncDecl(false)
		} else {
			forstmt.Init = parser.MainExpr()
		}
	}
	parser.want(TKSemi, ";")
	if parser.GetToken(0).Kind != TKSemi {
		forstmt.Cond = parser.MainExpr()
	}
	parser.want(TKSemi, ";")
	if parser.GetToken(0).Kind != TKRParen {
		forstmt.Post = parser.MainExpr()
	}
	parser.want(TKRParen, ")")
	forstmt.Body = parser.Statement()
	return forstmt
}

// SwitchStmt = 'switch' '(' Expr ')' '{' *CaseClause '}' .
// CaseClause = 'case' ExprList ':' Statement | 'default' ':' Statement .
func (parser *Parser) Switch() Stmt {
	///defer fmt.Printf("parser.Switch()\n")
	parser.want(TKSwitch, "switch")
	swtch := new(SwitchStmt)
	copyPosToNode(&swtch.node, parser.GetToken(-1))
	parser.want(TKLParen, "(")
	swtch.Cond = parser.AssignExpr()
	parser.want(TKRParen, ")")
	parser.want(TKLCurl, "{")
	bad_case := false
	for t := parser.GetToken(0); t.Kind != TKEoF && t.Kind != TKRCurl; t = parser.GetToken(0) {
		switch t.Kind {
		case TKCase:
			// next do case expressions:
			_case := new(CaseStmt)
			copyPosToNode(&_case.node, parser.GetToken(0))
			parser.Advance(1)
			case_expr := parser.MainExpr()
			if _, is_bad := case_expr.(*BadExpr); is_bad {
				n := case_expr.Tok()
				parser.MsgSpan.PrepNote(n.Span, "offending case expression(s).")
				parser.syntaxErr("bad case expr.")
				bad_case = true
				goto errd_case
			} else {
				_case.Case = case_expr
			}
			parser.want(TKColon, ":")
			_case.Body = parser.Statement()
			swtch.Cases = append(swtch.Cases, _case)
		case TKDefault:
			parser.Advance(1)
			parser.want(TKColon, ":")
			swtch.Default = parser.Statement()
		default:
			parser.MsgSpan.PrepNote(t.Span, "illegal switch case.")
			parser.syntaxErr("bad switch control label: %v.", t)
			bad := new(BadStmt)
			copyPosToNode(&bad.node, parser.GetToken(0))
			return bad
		}
	}
	parser.want(TKRCurl, "}")
errd_case:
	if bad_case {
		bad := new(BadStmt)
		copyPosToNode(&bad.node, parser.GetToken(0))
		return bad
	}
	return swtch
}


// Expr = AssignExpr *( ',' AssignExpr ) .
func (parser *Parser) MainExpr() Expr {
	///defer fmt.Printf("parser.MainExpr()\n")
	a := parser.AssignExpr()
	if parser.GetToken(0).Kind==TKComma {
		c := new(CommaExpr)
		copyPosToNode(&c.node, parser.GetToken(0))
		c.Exprs = append(c.Exprs, a)
		for t := parser.GetToken(0); t.Kind != TKEoF && t.Kind==TKComma; t = parser.GetToken(0) {
			///time.Sleep(100 * time.Millisecond)
			parser.Advance(1)
			c.Exprs = append(c.Exprs, parser.AssignExpr())
		}
		a = c
	}
	return a
}

// AssignExpr = SubMainExpr *( '['+' | '-' | '*' | '/' | '%' | '&' | '&~' | '|' | '^' | '<<' | '>>' | '>>>' ] =' SubMainExpr ) .
func (parser *Parser) AssignExpr() Expr {
	///defer fmt.Printf("parser.AssignExpr()\n")
	a := parser.SubMainExpr()
	for t := parser.GetToken(0); t.Kind >= TKAssign && t.Kind <= TKShLRA; t = parser.GetToken(0) {
		parser.Advance(1)
		assign_expr := new(BinExpr)
		copyPosToNode(&assign_expr.node, t)
		assign_expr.L = a
		assign_expr.Kind = t.Kind
		assign_expr.R = parser.SubMainExpr()
		a = assign_expr
	}
	return a
}

// SubMainExpr = LogicalOrExpr [ TernaryExpr ] .
func (parser *Parser) SubMainExpr() Expr {
	///defer fmt.Printf("parser.SubMainExpr()\n")
	a := parser.LogicalOrExpr()
	if parser.GetToken(0).Kind==TKQMark {
		// ternary
		a = parser.DoTernary(a)
	}
	return a
}

// TernaryExpr = '?' SubMainExpr ':' Expr .
func (parser *Parser) DoTernary(a Expr) Expr {
	///defer fmt.Printf("parser.DoTernary()\n")
	tk := parser.GetToken(0)
	t := new(TernaryExpr)
	copyPosToNode(&t.node, tk)
	t.A = a
	parser.Advance(1) // advance past question mark.
	t.B = parser.SubMainExpr()
	parser.want(TKColon, ":")
	t.C = parser.MainExpr()
	return t
}

// LogicalOrExpr = LogicalAndExpr *( '||' LogicalAndExpr ) .
func (parser *Parser) LogicalOrExpr() Expr {
	///defer fmt.Printf("parser.LogicalOrExpr()\n")
	if tIsEoF := parser.GetToken(0); tIsEoF.Kind==TKEoF {
		bad := new(BadExpr)
		copyPosToNode(&bad.node, tIsEoF)
		return bad
	}
	
	e := parser.LogicalAndExpr()
	for t := parser.GetToken(0); t.Kind==TKOrL; t = parser.GetToken(0) {
		b := new(BinExpr)
		copyPosToNode(&b.node, t)
		b.L = e
		b.Kind = t.Kind
		parser.Advance(1)
		b.R = parser.LogicalAndExpr()
		e = b
	}
	return e
}

// LogicalAndExpr = EqualExpr *( '&&' EqualExpr ) .
func (parser *Parser) LogicalAndExpr() Expr {
	///defer fmt.Printf("parser.LogicalAndExpr()\n")
	if tIsEoF := parser.GetToken(0); tIsEoF.Kind==TKEoF {
		bad := new(BadExpr)
		copyPosToNode(&bad.node, tIsEoF)
		return bad
	}
	
	e := parser.EqualExpr()
	for t := parser.GetToken(0); t.Kind==TKAndL; t = parser.GetToken(0) {
		b := new(BinExpr)
		copyPosToNode(&b.node, t)
		b.L = e
		b.Kind = t.Kind
		parser.Advance(1)
		b.R = parser.EqualExpr()
		e = b
	}
	return e
}

// EqualExpr = RelExpr *( ( '==' | '!=' ) RelExpr ) .
func (parser *Parser) EqualExpr() Expr {
	///defer fmt.Printf("parser.EqualExpr()\n")
	if tIsEoF := parser.GetToken(0); tIsEoF.Kind==TKEoF {
		bad := new(BadExpr)
		copyPosToNode(&bad.node, tIsEoF)
		return bad
	}
	
	e := parser.RelExpr()
	for t := parser.GetToken(0); t.Kind==TKEq || t.Kind==TKNotEq; t = parser.GetToken(0) {
		b := new(BinExpr)
		copyPosToNode(&b.node, t)
		b.L = e
		b.Kind = t.Kind
		parser.Advance(1)
		b.R = parser.RelExpr()
		e = b
	}
	return e
}

// RelExpr = BitOrExpr *( ( '<[=]' | '>[=]' ) BitOrExpr ) .
func (parser *Parser) RelExpr() Expr {
	///defer fmt.Printf("parser.RelExpr()\n")
	if tIsEoF := parser.GetToken(0); tIsEoF.Kind==TKEoF {
		bad := new(BadExpr)
		copyPosToNode(&bad.node, tIsEoF)
		return bad
	}
	
	e := parser.BitOrExpr()
	if t := parser.GetToken(0); t.Kind>=TKLess && t.Kind<=TKLessE {
		chain := new(ChainExpr)
		copyPosToNode(&chain.node, t)
		chain.A = e
		for n := t; n.Kind>=TKLess && n.Kind<=TKLessE; n = parser.GetToken(0) {
			chain.Kinds = append(chain.Kinds, n.Kind)
			parser.Advance(1)
			chain.Bs = append(chain.Bs, parser.BitOrExpr())
		}
		e = chain
	}
	return e
}

// BitOrExpr = BitXorExpr *( '|' BitXorExpr ) .
func (parser *Parser) BitOrExpr() Expr {
	///defer fmt.Printf("parser.BitOrExpr()\n")
	if tIsEoF := parser.GetToken(0); tIsEoF.Kind==TKEoF {
		bad := new(BadExpr)
		copyPosToNode(&bad.node, tIsEoF)
		return bad
	}
	
	e := parser.BitXorExpr()
	for t := parser.GetToken(0); t.Kind==TKOr; t = parser.GetToken(0) {
		b := new(BinExpr)
		copyPosToNode(&b.node, t)
		b.L = e
		b.Kind = t.Kind
		parser.Advance(1)
		b.R = parser.BitXorExpr()
		e = b
	}
	return e
}

// BitXorExpr = BitAndExpr *( '^' BitAndExpr ) .
func (parser *Parser) BitXorExpr() Expr {
	///defer fmt.Printf("parser.BitXorExpr()\n")
	if tIsEoF := parser.GetToken(0); tIsEoF.Kind==TKEoF {
		bad := new(BadExpr)
		copyPosToNode(&bad.node, tIsEoF)
		return bad
	}
	
	e := parser.BitAndExpr()
	for t := parser.GetToken(0); t.Kind==TKXor; t = parser.GetToken(0) {
		b := new(BinExpr)
		copyPosToNode(&b.node, t)
		b.L = e
		b.Kind = t.Kind
		parser.Advance(1)
		b.R = parser.BitAndExpr()
		e = b
	}
	return e
}

// BitAndExpr = ShiftExpr *( ('&' | '&~') ShiftExpr ) .
func (parser *Parser) BitAndExpr() Expr {
	///defer fmt.Printf("parser.BitAndExpr()\n")
	if tIsEoF := parser.GetToken(0); tIsEoF.Kind==TKEoF {
		bad := new(BadExpr)
		copyPosToNode(&bad.node, tIsEoF)
		return bad
	}
	
	e := parser.ShiftExpr()
	for t := parser.GetToken(0); t.Kind==TKAnd || t.Kind==TKAndNot; t = parser.GetToken(0) {
		b := new(BinExpr)
		copyPosToNode(&b.node, t)
		b.L = e
		b.Kind = t.Kind
		parser.Advance(1)
		b.R = parser.ShiftExpr()
		e = b
	}
	return e
}

// ShiftExpr = AddExpr *( ( '<<' | '>>' | '>>>' ) AddExpr ) .
func (parser *Parser) ShiftExpr() Expr {
	///defer fmt.Printf("parser.ShiftExpr()\n")
	if tIsEoF := parser.GetToken(0); tIsEoF.Kind==TKEoF {
		bad := new(BadExpr)
		copyPosToNode(&bad.node, tIsEoF)
		return bad
	}
	
	e := parser.AddExpr()
	for t := parser.GetToken(0); t.Kind>=TKShAL && t.Kind<=TKShLR; t = parser.GetToken(0) {
		b := new(BinExpr)
		copyPosToNode(&b.node, t)
		b.L = e
		b.Kind = t.Kind
		parser.Advance(1)
		b.R = parser.AddExpr()
		e = b
	}
	return e
}

// AddExpr = MulExpr *( ( '+' | '-' ) MulExpr ) .
func (parser *Parser) AddExpr() Expr {
	///defer fmt.Printf("parser.AddExpr()\n")
	if tIsEoF := parser.GetToken(0); tIsEoF.Kind==TKEoF {
		bad := new(BadExpr)
		copyPosToNode(&bad.node, tIsEoF)
		return bad
	}
	
	e := parser.MulExpr()
	for t := parser.GetToken(0); t.Kind==TKAdd || t.Kind==TKSub; t = parser.GetToken(0) {
		b := new(BinExpr)
		copyPosToNode(&b.node, t)
		b.L = e
		b.Kind = t.Kind
		parser.Advance(1)
		b.R = parser.MulExpr()
		e = b
	}
	return e
}

// MulExpr = PrefixExpr *( ( '*' | '/' | '%' ) PrefixExpr ) .
func (parser *Parser) MulExpr() Expr {
	///defer fmt.Printf("parser.MulExpr()\n")
	if tIsEoF := parser.GetToken(0); tIsEoF.Kind==TKEoF {
		bad := new(BadExpr)
		copyPosToNode(&bad.node, tIsEoF)
		return bad
	}
	
	e := parser.PrefixExpr()
	for t := parser.GetToken(0); t.Kind==TKMul || t.Kind==TKDiv || t.Kind==TKMod; t = parser.GetToken(0) {
		b := new(BinExpr)
		copyPosToNode(&b.node, t)
		b.L = e
		b.Kind = t.Kind
		parser.Advance(1)
		b.R = parser.PrefixExpr()
		e = b
	}
	return e
}

// PrefixExpr = *( '!' | '~' | '-' | '++' | '--' | 'sizeof' | 'new' ) PostfixExpr .
func (parser *Parser) PrefixExpr() Expr {
	///defer fmt.Printf("parser.PrefixExpr()\n")
	// certain patterns are allowed to recursively run Prefix.
	switch t := parser.GetToken(0); t.Kind {
	case TKIncr, TKDecr, TKNot, TKCompl, TKSub, TKSizeof, TKNew:
		n := new(UnaryExpr)
		parser.Advance(1)
		copyPosToNode(&n.node, t)
		n.X = parser.PrefixExpr()
		n.Kind = t.Kind
		return n
	default:
		return parser.PostfixExpr()
	}
}

// TypeExpr = '<' ( ident | '[u]int[8|16|32|64|n]' | 'float' | 'char' | 'bool' ) '>' .
func (parser *Parser) TypeExpr(need_carots bool) Expr {
	///defer fmt.Printf("parser.TypeExpr()\n")
	if tIsEoF := parser.GetToken(0); tIsEoF.Kind==TKEoF {
		bad := new(BadExpr)
		copyPosToNode(&bad.node, tIsEoF)
		return bad
	}
	
	ret_expr := Expr(nil)
	if need_carots {
		parser.want(TKLess, "<")
	}
	
	if t := parser.GetToken(0); t.IsType() || t.Kind==TKIdent {
		texp := new(TypedExpr)
		copyPosToNode(&texp.node, t)
		texp.TypeName = t
		ret_expr = texp
		parser.Advance(1)
	} else {
		parser.MsgSpan.PrepNote(t.Span, "expected type here.")
		parser.syntaxErr("missing required type expression.")
		bad := new(BadExpr)
		copyPosToNode(&bad.node, t)
		ret_expr = bad
	}
	if need_carots {
		parser.want(TKGreater, ">")
	}
	return ret_expr
}

// ViewAsExpr = TypeExpr '(' MainExpr ')' .
func (parser *Parser) ViewAsExpr() Expr {
	///defer fmt.Printf("parser.ViewAsExpr()\n")
	view_as := new(ViewAsExpr)
	parser.want(TKViewAs, "view_as")
	copyPosToNode(&view_as.node, parser.GetToken(-1))
	view_as.Type = parser.TypeExpr(true)
	parser.want(TKLParen, "(")
	view_as.X = parser.MainExpr()
	parser.want(TKRParen, ")")
	return view_as
}

// NamedArgExpr = '.' AssignExpr .
// ExprList = START ListedExpr *( SEP ListedExpr ) END .
// ListedExpr = NamedArgExpr | AssignExpr .
func (parser *Parser) ExprList(end, sep TokenKind, sep_at_end bool) []Expr {
	///defer fmt.Printf("parser.ExprList()\n")
	var exprs []Expr
	for t := parser.GetToken(0); t.Kind != TKEoF && t.Kind != end; t = parser.GetToken(0) {
		if !sep_at_end && len(exprs) > 0 {
			parser.want(sep, TokenToStr[sep])
		}
		// SP allows setting your params by name.
		// '.param_name = expression'
		// '.' Name '=' Expr
		if parser.GetToken(0).Kind==TKDot {
			parser.Advance(1)
			named_arg := new(NamedArg)
			copyPosToNode(&named_arg.node, parser.GetToken(-1))
			if iden := parser.GetToken(0); iden.Kind != TKIdent {
				parser.MsgSpan.PrepNote(t.Span, "")
				parser.syntaxErr("expected identifier for named arg.")
			}
			named_arg.X = parser.AssignExpr()
			exprs = append(exprs, named_arg)
		} else {
			exprs = append(exprs, parser.AssignExpr())
		}
		
		if sep_at_end && parser.GetToken(0).Kind==sep {
			parser.Advance(1)
		}
	}
	return exprs
}


// PostfixExpr = Primary *( '.' identifier | '[' Expr ']' | '(' [ ExprList ] ')' | '::' identifier | '++' | '--' ) .
func (parser *Parser) PostfixExpr() Expr {
	///defer fmt.Printf("parser.PostfixExpr()\n")
	n := Expr(nil)
	if t := parser.GetToken(0); t.Kind==TKViewAs {
		n = parser.ViewAsExpr()
	} else {
		n = parser.PrimaryExpr()
	}
	
	for t := parser.GetToken(0); t.Kind==TKDot || t.Kind==TKLBrack || t.Kind==TKLParen || t.Kind==TK2Colons || t.Kind==TKIncr || t.Kind==TKDecr; t = parser.GetToken(0) {
		parser.Advance(1)
		switch t.Kind {
		case TKDot:
			field := new(FieldExpr)
			copyPosToNode(&field.node, t)
			field.X = n
			field.Sel = parser.PrimaryExpr()
			n = field
		case TK2Colons:
			namespc := new(NameSpaceExpr)
			copyPosToNode(&namespc.node, t)
			namespc.N = n
			namespc.Id = parser.PrimaryExpr()
			n = namespc
		case TKLBrack:
			arr := new(IndexExpr)
			copyPosToNode(&arr.node, t)
			arr.X = n
			if parser.GetToken(0).Kind != TKRBrack {
				arr.Index = parser.MainExpr()
			}
			parser.want(TKRBrack, "]")
			n = arr
		case TKIncr, TKDecr:
			incr := new(UnaryExpr)
			copyPosToNode(&incr.node, t)
			incr.X = n
			incr.Kind = t.Kind
			incr.Post = true
			n = incr
		case TKLParen:
			call := new(CallExpr)
			copyPosToNode(&call.node, t)
			call.Func = n
			call.ArgList = parser.ExprList(TKRParen, TKComma, false)
			parser.want(TKRParen, ")")
			n = call
		}
	}
	return n
}

// BoolLit = 'true' | 'false' .
// BasicLit = int_lit | rune_lit | string_lit .
// BracketExpr = '{' ExprList '}' .
// Primary = BasicLit | identifier | 'operator' op | BoolLit | 'this' | 'null' | '...' | '(' Expr ')' | BracketExpr .
func (parser *Parser) PrimaryExpr() Expr {
	///defer fmt.Printf("parser.PrimaryExpr()\n")
	ret_expr := Expr(nil)
	if tIsEoF := parser.GetToken(0); tIsEoF.Kind==TKEoF {
		bad := new(BadExpr)
		copyPosToNode(&bad.node, tIsEoF)
		return bad
	}
	
	if parser.GetToken(0).IsType() {
		return parser.TypeExpr(false)
	}
	switch prim := parser.GetToken(0); prim.Kind {
	case TKEllipses:
		ell := new(EllipsesExpr)
		copyPosToNode(&ell.node, prim)
		ret_expr = ell
	case TKLParen:
		parser.Advance(1)
		ret_expr = parser.MainExpr()
		if t := parser.GetToken(0); t.Kind != TKRParen {
			parser.MsgSpan.PrepNote(prim.Span, "starting '(' parenthesis here.")
			parser.syntaxErr("missing ending ')' right paren for nested expression")
			bad := new(BadExpr)
			copyPosToNode(&bad.node, parser.GetToken(0))
			ret_expr = bad
		}
	case TKLCurl:
		brktexpr := new(BracketExpr)
		copyPosToNode(&brktexpr.node, prim)
		parser.Advance(1)
		brktexpr.Exprs = parser.ExprList(TKRCurl, TKComma, true)
		ret_expr = brktexpr
	case TKOperator: // operator# like operator% or operator+
		operator := prim.Lexeme
		operator += parser.GetToken(1).Lexeme
		parser.Advance(1)
		iden := new(Name)
		iden.Value = operator
		copyPosToNode(&iden.node, prim)
		ret_expr = iden
	case TKIdent:
		iden := new(Name)
		iden.Value = prim.Lexeme
		copyPosToNode(&iden.node, prim)
		ret_expr = iden
	case TKIntLit:
		num := new(BasicLit)
		num.Value = prim.Lexeme
		num.Kind = IntLit
		copyPosToNode(&num.node, prim)
		ret_expr = num
	case TKFloatLit:
		num := new(BasicLit)
		num.Value = prim.Lexeme
		num.Kind = FloatLit
		copyPosToNode(&num.node, prim)
		ret_expr = num
	case TKStrLit:
		str := new(BasicLit)
		str.Value = prim.Lexeme
		str.Kind = StringLit
		copyPosToNode(&str.node, prim)
		ret_expr = str
	case TKCharLit:
		str := new(BasicLit)
		str.Value = prim.Lexeme
		str.Kind = CharLit
		copyPosToNode(&str.node, prim)
		ret_expr = str
	case TKThis:
		this := new(ThisExpr)
		copyPosToNode(&this.node, prim)
		ret_expr = this
	case TKTrue, TKFalse:
		boolean := new(BasicLit)
		boolean.Value = prim.Lexeme
		boolean.Kind = BoolLit
		copyPosToNode(&boolean.node, prim)
		ret_expr = boolean
	case TKNull:
		null := new(NullExpr)
		copyPosToNode(&null.node, prim)
		ret_expr = null
	case TKFunction:
		func_lit := new(FuncLit)
		copyPosToNode(&func_lit.node, prim)
		func_lit.Sig  = parser.DoFuncSignature()
		func_lit.Body = parser.DoBlock()
		return func_lit
	default:
		parser.MsgSpan.PrepNote(prim.Span, "offending token.")
		parser.syntaxErr("bad primary expression '%s'", prim.Lexeme)
		bad := new(BadExpr)
		copyPosToNode(&bad.node, prim)
		ret_expr = bad
	}
	parser.Advance(1)
	return ret_expr
}