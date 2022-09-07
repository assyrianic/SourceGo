package SPTools


func (parser *Parser) OldStart() Node {
	if parser.TokenReader.Len() <= 0 {
		return nil
	}
	return parser.DoOldBlock()
}


// OldBlockStmt = '{' *OldStatement '}' .
func (parser *Parser) DoOldBlock() Stmt {
	block := new(BlockStmt)
	parser.want(TKLCurl, "{")
	copyPosToNode(&block.node, parser.GetToken(-1))
	for t := parser.GetToken(0); t.Kind != TKRCurl && t.Kind != TKEoF; t = parser.GetToken(0) {
		///time.Sleep(100 * time.Millisecond)
		///fmt.Printf("current tok: %v\n", t)
		n := parser.OldStatement()
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
 * OldStatement = OldIfStmt | OldWhileStmt | OldForStmt | OldSwitchStmt | OldBlockStmt |
 *             OldRetStmt | OldAssertStmt | OldDeclStmt | OldExprStmt .
 */
func (parser *Parser) OldStatement() Stmt {
	///defer fmt.Printf("parser.OldStatement()\n")
	if tIsEoF := parser.GetToken(0); tIsEoF.Kind==TKEoF {
		bad := new(BadStmt)
		copyPosToNode(&bad.node, tIsEoF)
		return bad
	}
	switch t := parser.GetToken(0); t.Kind {
	/*case TKConst, TKPublic, TKForward, TKNative, TKStock, TKDecl, TKStatic, TKVar, TKNew:
		// parse declaration.
		vardecl := new(DeclStmt)
		copyPosToNode(&vardecl.node, t)
		vardecl.D = parser.OldDoVarOrFuncDecl(false)
		if !parser.got(TKSemi) {
			return parser.noSemi()
		}
		return vardecl
		*/
	case TKLCurl:
		return parser.DoOldBlock()
	case TKIf:
		return parser.DoOldIf()
	case TKDo, TKWhile:
		return parser.OldWhile(t)
	case TKFor:
		return parser.DoOldFor()
	case TKSwitch:
		return parser.OldSwitch()
	case TKReturn:
		// RetStmt = 'return' [ Expr ] ';' .
		ret := new(RetStmt)
		copyPosToNode(&ret.node, t)
		parser.Advance(1)
		if parser.GetToken(0).Kind != TKSemi {
			ret.X = parser.OldMainExpr()
			if !parser.got(TKSemi) {
				return parser.noSemi()
			}
			return ret
		} else {
			parser.Advance(1)
			return ret
		}
	case TKAssert:
		// AssertStmt = 'assert' Expr ';' .
		asrt := new(AssertStmt)
		copyPosToNode(&asrt.node, t)
		parser.Advance(1)
		asrt.X = parser.OldMainExpr()
		if !parser.got(TKSemi) {
			return parser.noSemi()
		}
		return asrt
	case TKIdent:
		// vast majority of the times,
		// an expression starts with an identifier.
		exp := new(ExprStmt)
		copyPosToNode(&exp.node, t)
		exp.X = parser.OldMainExpr()
		if !parser.got(TKSemi) {
			return parser.noSemi()
		}
		return exp
	case TKIncr, TKDecr:
		exp := new(ExprStmt)
		copyPosToNode(&exp.node, t)
		exp.X = parser.OldMainExpr()
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
		parser.syntaxErr("unknown statement: %q", t.String())
		bad := new(BadStmt)
		copyPosToNode(&bad.node, parser.GetToken(0))
		parser.Advance(1)
		return bad
	}
}

// DoStmt = 'do' OldStatement 'while' '(' Expr ')' ';' .
// WhileStmt = 'while' '(' Expr ')' OldStatement .
func (parser *Parser) OldWhile(t Token) Stmt {
	///defer fmt.Printf("parser.OldWhile()\n")
	while := new(WhileStmt)
	copyPosToNode(&while.node, t)
	parser.Advance(1)
	if t.Kind==TKDo {
		// do-while
		while.Do = true
		while.Body = parser.OldStatement()
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
		while.Body = parser.OldStatement()
	}
	return while
}

// IfStmt = 'if' '(' Expr ')' OldStatement [ 'else' OldStatement ] .
func (parser *Parser) DoOldIf() Stmt {
	///defer fmt.Printf("parser.DoOldIf()\n")
	parser.want(TKIf, "if")
	ifstmt := new(IfStmt)
	copyPosToNode(&ifstmt.node, parser.GetToken(-1))
	parser.want(TKLParen, "(")
	ifstmt.Cond = parser.MainExpr()
	parser.want(TKRParen, ")")
	ifstmt.Then = parser.OldStatement()
	if parser.GetToken(0).Kind==TKElse {
		parser.Advance(1)
		ifstmt.Else = parser.OldStatement()
		// personally, I'd prefer to fix the not-needing-{ thing but whatever.
		/*
		switch t := parser.GetToken(0); t.Kind {
		case TKIf:
			ifstmt.Else = parser.DoOldIf()
		case TKLCurl:
			ifstmt.Else = parser.DoBlock()
		default:
			parser.syntaxErr("ill-formed else block, missing 'if' or { curl.")
			bad := new(BadStmt)
			copyPosToNode(&bad.node, t)
			ifstmt.Else = bad
		}
		*/
	}
	return ifstmt
}

// OldForStmt = 'for' '(' [ OldDecl | OldExpr ] ';' [ OldExpr ] ';' [ OldExpr ] ')' OldStatement .
func (parser *Parser) DoOldFor() Stmt {
	parser.want(TKFor, "for")
	forstmt := new(ForStmt)
	copyPosToNode(&forstmt.node, parser.GetToken(-1))
	parser.want(TKLParen, "(")
	if parser.GetToken(0).Kind != TKSemi {
		if t := parser.GetToken(0); t.Kind==TKNew || t.IsStorageClass() {
			//forstmt.Init = parser.DoOldVarOrFuncDecl(false)
		} else {
			forstmt.Init = parser.OldMainExpr()
		}
	}
	parser.want(TKSemi, ";")
	if parser.GetToken(0).Kind != TKSemi {
		forstmt.Cond = parser.OldMainExpr()
	}
	parser.want(TKSemi, ";")
	if parser.GetToken(0).Kind != TKRParen {
		forstmt.Post = parser.OldMainExpr()
	}
	parser.want(TKRParen, ")")
	forstmt.Body = parser.OldStatement()
	return forstmt
}

// OldSwitchStmt = 'switch' '(' OldExpr ')' '{' *OldCaseClause '}' .
// OldCaseClause = 'case' OldExprList ':' OldStatement | 'default' ':' OldStatement .
func (parser *Parser) OldSwitch() Stmt {
	parser.want(TKSwitch, "switch")
	swtch := new(SwitchStmt)
	copyPosToNode(&swtch.node, parser.GetToken(-1))
	parser.want(TKLParen, "(")
	swtch.Cond = parser.OldAssignExpr()
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
			case_expr := parser.OldMainExpr()
			if _, is_bad := case_expr.(*BadExpr); is_bad {
				parser.syntaxErr("bad case expr.")
				bad_case = true
				goto errd_case
			} else {
				_case.Case = case_expr
			}
			parser.want(TKColon, ":")
			_case.Body = parser.OldStatement()
			swtch.Cases = append(swtch.Cases, _case)
		case TKDefault:
			parser.Advance(1)
			parser.want(TKColon, ":")
			swtch.Default = parser.OldStatement()
		default:
			parser.syntaxErr("bad switch label: %+v.", parser.GetToken(0))
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
/***********************************************************************************/
//                                End Statements
/***********************************************************************************/

/***********************************************************************************/
//                                Start Expressions
/***********************************************************************************/
// OldExpr = OldAssignExpr *( ',' OldAssignExpr ) .
func (parser *Parser) OldMainExpr() Expr {
	a := parser.OldAssignExpr()
	if parser.GetToken(0).Kind==TKComma {
		c := new(CommaExpr)
		copyPosToNode(&c.node, parser.GetToken(0))
		c.Exprs = append(c.Exprs, a)
		for t := parser.GetToken(0); t.Kind != TKEoF && t.Kind==TKComma; t = parser.GetToken(0) {
			parser.Advance(1)
			c.Exprs = append(c.Exprs, parser.OldAssignExpr())
		}
		a = c
	}
	return a
}

// OldAssignExpr = OldSubMainExpr *( '['+' | '-' | '*' | '/' | '%' | '&' | '|' | '^' | '<<' | '>>' | '>>>' ] =' OldSubMainExpr ) .
func (parser *Parser) OldAssignExpr() Expr {
	a := parser.OldSubMainExpr()
	for t := parser.GetToken(0); t.Kind >= TKAssign && t.Kind <= TKShLRA; t = parser.GetToken(0) {
		parser.Advance(1)
		assign_expr := new(BinExpr)
		copyPosToNode(&assign_expr.node, t)
		assign_expr.L = a
		assign_expr.Kind = t.Kind
		assign_expr.R = parser.OldSubMainExpr()
		a = assign_expr
	}
	return a
}

// OldSubMainExpr = OldLogicalOrExpr [ TernaryExpr ] .
func (parser *Parser) OldSubMainExpr() Expr {
	a := parser.OldLogicalOrExpr()
	if parser.GetToken(0).Kind==TKQMark {
		// ternary
		a = parser.DoOldTernary(a)
	}
	return a
}

// OldTernaryExpr = '?' OldSubMainExpr ':' Expr .
func (parser *Parser) DoOldTernary(a Expr) Expr {
	tk := parser.GetToken(0)
	t := new(TernaryExpr)
	copyPosToNode(&t.node, tk)
	t.A = a
	parser.Advance(1) // advance past question mark.
	t.B = parser.OldSubMainExpr()
	parser.want(TKColon, ":")
	t.C = parser.OldMainExpr()
	return t
}

// OldLogicalOrExpr = OldLogicalAndExpr *( '||' OldLogicalAndExpr ) .
func (parser *Parser) OldLogicalOrExpr() Expr {
	if tIsEoF := parser.GetToken(0); tIsEoF.Kind==TKEoF {
		bad := new(BadExpr)
		copyPosToNode(&bad.node, tIsEoF)
		return bad
	}
	
	e := parser.OldLogicalAndExpr()
	for t := parser.GetToken(0); t.Kind==TKOrL; t = parser.GetToken(0) {
		b := new(BinExpr)
		copyPosToNode(&b.node, t)
		b.L = e
		b.Kind = t.Kind
		parser.Advance(1)
		b.R = parser.OldLogicalAndExpr()
		e = b
	}
	return e
}

// OldLogicalAndExpr = OldEqualExpr *( '&&' OldEqualExpr ) .
func (parser *Parser) OldLogicalAndExpr() Expr {
	if tIsEoF := parser.GetToken(0); tIsEoF.Kind==TKEoF {
		bad := new(BadExpr)
		copyPosToNode(&bad.node, tIsEoF)
		return bad
	}
	
	e := parser.OldEqualExpr()
	for t := parser.GetToken(0); t.Kind==TKAndL; t = parser.GetToken(0) {
		b := new(BinExpr)
		copyPosToNode(&b.node, t)
		b.L = e
		b.Kind = t.Kind
		parser.Advance(1)
		b.R = parser.OldEqualExpr()
		e = b
	}
	return e
}

// OldEqualExpr = OldRelExpr *( ( '==' | '!=' ) OldRelExpr ) .
func (parser *Parser) OldEqualExpr() Expr {
	if tIsEoF := parser.GetToken(0); tIsEoF.Kind==TKEoF {
		bad := new(BadExpr)
		copyPosToNode(&bad.node, tIsEoF)
		return bad
	}
	
	e := parser.OldRelExpr()
	for t := parser.GetToken(0); t.Kind==TKEq || t.Kind==TKNotEq; t = parser.GetToken(0) {
		b := new(BinExpr)
		copyPosToNode(&b.node, t)
		b.L = e
		b.Kind = t.Kind
		parser.Advance(1)
		b.R = parser.OldRelExpr()
		e = b
	}
	return e
}

// OldRelExpr = OldBitOrExpr *( ( '<[=]' | '>[=]' ) OldBitOrExpr ) .
func (parser *Parser) OldRelExpr() Expr {
	///defer fmt.Printf("parser.OldRelExpr()\n")
	if tIsEoF := parser.GetToken(0); tIsEoF.Kind==TKEoF {
		bad := new(BadExpr)
		copyPosToNode(&bad.node, tIsEoF)
		return bad
	}
	
	e := parser.OldBitOrExpr()
	for t := parser.GetToken(0); t.Kind>=TKLess && t.Kind<=TKLessE; t = parser.GetToken(0) {
		b := new(BinExpr)
		copyPosToNode(&b.node, t)
		b.L = e
		b.Kind = t.Kind
		parser.Advance(1)
		b.R = parser.OldBitOrExpr()
		e = b
	}
	return e
}

// OldBitOrExpr = OldBitXorExpr *( '|' OldBitXorExpr ) .
func (parser *Parser) OldBitOrExpr() Expr {
	///defer fmt.Printf("parser.OldBitOrExpr()\n")
	if tIsEoF := parser.GetToken(0); tIsEoF.Kind==TKEoF {
		bad := new(BadExpr)
		copyPosToNode(&bad.node, tIsEoF)
		return bad
	}
	
	e := parser.OldBitXorExpr()
	for t := parser.GetToken(0); t.Kind==TKOr; t = parser.GetToken(0) {
		b := new(BinExpr)
		copyPosToNode(&b.node, t)
		b.L = e
		b.Kind = t.Kind
		parser.Advance(1)
		b.R = parser.OldBitXorExpr()
		e = b
	}
	return e
}

// OldBitXorExpr = OldBitAndExpr *( '^' OldBitAndExpr ) .
func (parser *Parser) OldBitXorExpr() Expr {
	///defer fmt.Printf("parser.OldBitXorExpr()\n")
	if tIsEoF := parser.GetToken(0); tIsEoF.Kind==TKEoF {
		bad := new(BadExpr)
		copyPosToNode(&bad.node, tIsEoF)
		return bad
	}
	
	e := parser.OldBitAndExpr()
	for t := parser.GetToken(0); t.Kind==TKXor; t = parser.GetToken(0) {
		b := new(BinExpr)
		copyPosToNode(&b.node, t)
		b.L = e
		b.Kind = t.Kind
		parser.Advance(1)
		b.R = parser.OldBitAndExpr()
		e = b
	}
	return e
}

// OldBitAndExpr = OldShiftExpr *( '&' OldShiftExpr ) .
func (parser *Parser) OldBitAndExpr() Expr {
	///defer fmt.Printf("parser.OldBitAndExpr()\n")
	if tIsEoF := parser.GetToken(0); tIsEoF.Kind==TKEoF {
		bad := new(BadExpr)
		copyPosToNode(&bad.node, tIsEoF)
		return bad
	}
	
	e := parser.OldShiftExpr()
	for t := parser.GetToken(0); t.Kind==TKAnd; t = parser.GetToken(0) {
		b := new(BinExpr)
		copyPosToNode(&b.node, t)
		b.L = e
		b.Kind = t.Kind
		parser.Advance(1)
		b.R = parser.OldShiftExpr()
		e = b
	}
	return e
}

// OldShiftExpr = OldAddExpr *( ( '<<' | '>>' | '>>>' ) OldAddExpr ) .
func (parser *Parser) OldShiftExpr() Expr {
	///defer fmt.Printf("parser.OldShiftExpr()\n")
	if tIsEoF := parser.GetToken(0); tIsEoF.Kind==TKEoF {
		bad := new(BadExpr)
		copyPosToNode(&bad.node, tIsEoF)
		return bad
	}
	
	e := parser.OldAddExpr()
	for t := parser.GetToken(0); t.Kind>=TKShAL && t.Kind<=TKShLR; t = parser.GetToken(0) {
		b := new(BinExpr)
		copyPosToNode(&b.node, t)
		b.L = e
		b.Kind = t.Kind
		parser.Advance(1)
		b.R = parser.OldAddExpr()
		e = b
	}
	return e
}

// OldAddExpr = OldMulExpr *( ( '+' | '-' ) OldMulExpr ) .
func (parser *Parser) OldAddExpr() Expr {
	///defer fmt.Printf("parser.OldAddExpr()\n")
	if tIsEoF := parser.GetToken(0); tIsEoF.Kind==TKEoF {
		bad := new(BadExpr)
		copyPosToNode(&bad.node, tIsEoF)
		return bad
	}
	
	e := parser.OldMulExpr()
	for t := parser.GetToken(0); t.Kind==TKAdd || t.Kind==TKSub; t = parser.GetToken(0) {
		b := new(BinExpr)
		copyPosToNode(&b.node, t)
		b.L = e
		b.Kind = t.Kind
		parser.Advance(1)
		b.R = parser.OldMulExpr()
		e = b
	}
	return e
}

// OldMulExpr = OldPrefixExpr *( ( '*' | '/' | '%' ) OldPrefixExpr ) .
func (parser *Parser) OldMulExpr() Expr {
	///defer fmt.Printf("parser.OldMulExpr()\n")
	if tIsEoF := parser.GetToken(0); tIsEoF.Kind==TKEoF {
		bad := new(BadExpr)
		copyPosToNode(&bad.node, tIsEoF)
		return bad
	}
	
	e := parser.OldPrefixExpr()
	for t := parser.GetToken(0); t.Kind==TKMul || t.Kind==TKDiv || t.Kind==TKMod; t = parser.GetToken(0) {
		b := new(BinExpr)
		copyPosToNode(&b.node, t)
		b.L = e
		b.Kind = t.Kind
		parser.Advance(1)
		b.R = parser.OldPrefixExpr()
		e = b
	}
	return e
}

// OldPrefixExpr = *( '!' | '~' | '-' | '++' | '--' | 'sizeof' | 'defined' ) OldPostfixExpr .
func (parser *Parser) OldPrefixExpr() Expr {
	// certain patterns are allowed to recursively run Prefix.
	switch t := parser.GetToken(0); t.Kind {
	case TKIncr, TKDecr, TKNot, TKCompl, TKSub, TKSizeof, TKDefined:
		n := new(UnaryExpr)
		parser.Advance(1)
		copyPosToNode(&n.node, t)
		n.X = parser.OldPrefixExpr()
		n.Kind = t.Kind
		return n
	default:
		return parser.OldPostfixExpr()
	}
}

// OldTypeExpr = ident | 'Float' | 'String' | 'bool' .
func (parser *Parser) OldTypeExpr() Expr {
	if tIsEoF := parser.GetToken(0); tIsEoF.Kind==TKEoF {
		bad := new(BadExpr)
		copyPosToNode(&bad.node, tIsEoF)
		return bad
	}
	
	ret_expr := Expr(nil)
	if t := parser.GetToken(0); t.IsType() || t.Kind==TKIdent {
		texp := new(TypedExpr)
		copyPosToNode(&texp.node, t)
		texp.TypeName = t
		texp.TypeName.Lexeme = func() string {
			switch t.Lexeme {
			case "String":
				return "char"
			case "Float":
				return "float"
			case "_":
				return "int"
			default:
				return t.Lexeme
			}
		}()
		ret_expr = texp
		parser.Advance(1)
	} else {
		parser.syntaxErr("missing type expression.")
		bad := new(BadExpr)
		copyPosToNode(&bad.node, t)
		ret_expr = bad
	}
	return ret_expr
}


// OldNamedArgExpr = '.' OldAssignExpr .
// OldExprList = START OldListedExpr *( SEP OldListedExpr ) END .
// OldListedExpr = NamedArgExpr | OldAssignExpr .
func (parser *Parser) OldExprList(end, sep TokenKind, sep_at_end bool) []Expr {
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
				parser.syntaxErr("expected identifier for named arg.")
			}
			named_arg.X = parser.OldAssignExpr()
			exprs = append(exprs, named_arg)
		} else {
			exprs = append(exprs, parser.OldAssignExpr())
		}
		
		if sep_at_end && parser.GetToken(0).Kind==sep {
			parser.Advance(1)
		}
	}
	return exprs
}

// OldPostfixExpr = OldPrimaryExpr *( ':' OldPrefixExpr | '[' Expr ']' | '(' [ OldExprList ] ')' | '++' | '--' ) .
func (parser *Parser) OldPostfixExpr() Expr {
	n := parser.OldPrimaryExpr()
	for t := parser.GetToken(0); t.Kind==TKColon || t.Kind==TKLBrack || t.Kind==TKLParen || t.Kind==TKIncr || t.Kind==TKDecr; t = parser.GetToken(0) {
		parser.Advance(1)
		switch t.Kind {
		case TKColon:
			// retagging.
			if _, is_type_expr := n.(*TypedExpr); !is_type_expr {
				break
			}
			view_as := new(ViewAsExpr)
			copyPosToNode(&view_as.node, t)
			view_as.Type = n
			view_as.X = parser.OldPrefixExpr()
			n = view_as
		case TKLBrack:
			arr := new(IndexExpr)
			copyPosToNode(&arr.node, t)
			arr.X = n
			if parser.GetToken(0).Kind != TKRBrack {
				arr.Index = parser.OldMainExpr()
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
			call.ArgList = parser.OldExprList(TKRParen, TKComma, false)
			parser.want(TKRParen, ")")
			n = call
		}
	}
	return n
}

// BoolLit = 'true' | 'false' .
// BasicLit = int_lit | rune_lit | string_lit .
// OldPrimary = BasicLit | identifier | 'operator' op | BoolLit | '...' | '(' Expr ')' | '{' OldExprList '}' .
func (parser *Parser) OldPrimaryExpr() Expr {
	ret_expr := Expr(nil)
	if tIsEoF := parser.GetToken(0); tIsEoF.Kind==TKEoF {
		bad := new(BadExpr)
		copyPosToNode(&bad.node, tIsEoF)
		return bad
	}
	
	if t := parser.GetToken(0); t.IsType() || t.Lexeme=="String" || t.Lexeme=="Float" || ( (t.Lexeme=="_" || t.Kind==TKIdent) && parser.GetToken(1).Kind==TKColon) {
		return parser.OldTypeExpr()
	}
	
	switch prim := parser.GetToken(0); prim.Kind {
	case TKEllipses:
		ell := new(EllipsesExpr)
		copyPosToNode(&ell.node, prim)
		ret_expr = ell
	case TKLParen:
		parser.Advance(1)
		ret_expr = parser.OldMainExpr()
		if t := parser.GetToken(0); t.Kind != TKRParen {
			parser.syntaxErr("missing ending ')' right paren for nested expression")
			bad := new(BadExpr)
			copyPosToNode(&bad.node, parser.GetToken(0))
			ret_expr = bad
		}
	case TKLCurl:
		brktexpr := new(BracketExpr)
		copyPosToNode(&brktexpr.node, prim)
		parser.Advance(1)
		brktexpr.Exprs = parser.OldExprList(TKRCurl, TKComma, true)
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
	case TKTrue, TKFalse:
		boolean := new(BasicLit)
		boolean.Value = prim.Lexeme
		boolean.Kind = BoolLit
		copyPosToNode(&boolean.node, prim)
		ret_expr = boolean
	default:
		parser.syntaxErr("bad primary expression '%s'", prim.Lexeme)
		bad := new(BadExpr)
		copyPosToNode(&bad.node, prim)
		ret_expr = bad
	}
	parser.Advance(1)
	return ret_expr
}

/* NOTES.
 * rename `String` to `char` & `Float` to `float`.
 * declarations: ( storageclass | 'decl' | 'new' | '&' ) [ Type ':' ] ident [ +('[' [ OldExpr ] ']') ] [ initializer ]
 */