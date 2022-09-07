package SPTools

import (
	"fmt"
	"os"
	"io"
	"io/ioutil"
	"strings"
)

// colorful strings for printing.
const (
	COLOR_RED     = "\x1B[31m"    // used for errors.
	COLOR_GREEN   = "\x1B[32m"
	COLOR_YELLOW  = "\x1B[33m"
	COLOR_BLUE    = "\x1B[34m"
	COLOR_MAGENTA = "\x1B[35m"    // used for warnings.
	COLOR_CYAN    = "\x1B[36m"
	COLOR_WHITE   = "\x1B[37m"
	COLOR_RESET   = "\033[0m"     // used to reset the color.
)


/* Example:
 * {error | warning} [{ErrCode}]: {Msg}
 * --> {file, line, col}
 * Line1 | Code1
 * Span1 | ^^^^^ note1
 * Line2 | Code2
 * Span1 | ----- note2
 * ...
 * LineN | CodeN
 * SpanN | ~~~~~ noteN
 */

/*
 * compute an array of the line start offsets
 */

type MsgSpan struct {
	// TODO: have MsgSpan hold filename?
	spans []Span
	notes []string
	code *[]string
}


func MakeMsgSpan(lines *[]string) MsgSpan {
	return MsgSpan{ code: lines }
}

func (m *MsgSpan) PrepNote(span Span, note_fmt string, args ...any) {
	m.spans = append(m.spans, span)
	m.notes = append(m.notes, fmt.Sprintf(note_fmt, args...))
}

func (m *MsgSpan) PurgeNotes() {
	m.spans = nil
	m.notes = nil
}

func (m *MsgSpan) Report(msgtype, errcode, color, msg_fmt, filename string, line, col *uint16, args ...any) string {
	var sb strings.Builder
	sb.WriteString(color)
	sb.WriteString(msgtype)
	if errcode != "" {
		sb.WriteRune('[')
		sb.WriteString(errcode)
		sb.WriteRune(']')
	}
	sb.WriteString(COLOR_RESET)
	sb.WriteString(": " + fmt.Sprintf(msg_fmt, args...))
	if filename != "" {
		sb.WriteRune('\n')
		sb.WriteString("--> ")
		sb.WriteString(filename)
		if line != nil {
			sb.WriteString(fmt.Sprintf(":%d", *line))
		}
		if col != nil {
			sb.WriteString(fmt.Sprintf(":%d", *col))
		}
	}
	
	if len(m.spans) > 0 {
		sb.WriteRune('\n')
		largest := 0
		for _, l := range m.spans {
			if largest < int(l.LineEnd) {
				largest = int(l.LineEnd)
			}
		}
		big_line_str := fmt.Sprintf("%d", largest)
		largest_span := len(big_line_str) + 1
		span_write := func (sb *strings.Builder, span int, c rune) {
			for i := 0; i < span; i++ {
				sb.WriteRune(c)
			}
		}
		span_write(&sb, largest_span, ' ')
		sb.WriteRune('|')
		sb.WriteRune('\n')
		for i := range m.spans {
			span := m.spans[i]
			for line := span.LineStart; line <= span.LineEnd; line++ {
				line_num_str := fmt.Sprintf("%d", line)
				sb.WriteString(line_num_str)
				line_num_len := len(line_num_str)
				span_write(&sb, largest_span - line_num_len, ' ')
				sb.WriteRune('|')
				code_line := (*m.code)[line - 1]
				sb.WriteString(code_line)
				sb.WriteRune('\n')
			}
			note := m.notes[i]
			span_write(&sb, largest_span, ' ')
			sb.WriteRune('|')
			span_write(&sb, int(span.ColStart), ' ')
			span_write(&sb, int(span.ColEnd - span.ColStart), '^')
			sb.WriteRune(' ')
			sb.WriteString(COLOR_CYAN);
			sb.WriteString(note);
			sb.WriteString(COLOR_RESET)
		}
	}
	return sb.String()
}

func SpewReport(w io.Writer, message string, msg_cnt *uint32) {
	fmt.Fprintf(w, "%s\n", message)
	if msg_cnt != nil {
		*msg_cnt++
	}
}


const (
	// Runs preprocessor.
	LEXFLAG_PREPROCESS     = (1 << iota)
	
	// Self explanatory, strips out all comment tokens.
	LEXFLAG_STRIP_COMMENTS = (1 << iota)
	
	// Keeps newlines.
	LEXFLAG_NEWLINES       = (1 << iota)
	
	// Adds #include <sourcemod> automatically.
	LEXFLAG_SM_INCLUDE     = (1 << iota)
	
	// Enable ALL the above flags.
	LEXFLAG_ALL            = -1
)


func loadFile(filename string) (string, string) {
	if text, read_err := ioutil.ReadFile(filename); read_err==nil {
		code := string(text)
		code = strings.ReplaceAll(code, "\r\n", "\n")
		code = strings.ReplaceAll(code, "\r", "\n")
		code = strings.ReplaceAll(code, "\t", "    ")
		return code, "none"
	} else {
		return "", read_err.Error()
	}
}

// Lexes and preprocesses a file, returning its token array.
func LexFile(filename string, flags int, macros map[string]Macro) (*TokenReader, bool) {
	code, err_str := loadFile(filename)
	if len(code) <= 0 {
		fmt.Fprintf(os.Stdout, "sptools %sIO error%s: **** file error:: '%s'. ****\n", COLOR_RED, COLOR_RESET, err_str)
		return &TokenReader{}, false
	}
	if flags & LEXFLAG_SM_INCLUDE > 0 {
		code = "#include <sourcemod>\n" + code
	}
	return finishLexing(Tokenize(code, filename), flags, macros)
}

func LexCodeString(code string, flags int, macros map[string]Macro) (*TokenReader, bool) {
	return finishLexing(Tokenize(code, ""), flags, macros)
}

func finishLexing(tr *TokenReader, flags int, macros map[string]Macro) (*TokenReader, bool) {
	if flags & LEXFLAG_PREPROCESS > 0 {
		if output, res := Preprocess(tr, flags, macros); res {
			*tr = *output
		} else {
			return nil, false
		}
	}
	tr = ConcatStringLiterals(tr)
	if flags & LEXFLAG_STRIP_COMMENTS > 0 {
		tr = RemoveComments(tr)
	}
	tr = StripSpaceTokens(tr, flags & LEXFLAG_NEWLINES > 0)
	return tr, true
}

func MakeParser(tr *TokenReader) Parser {
	return Parser{ TokenReader: tr }
}


func ParseTokens(tr *TokenReader, old bool) Node {
	parser := MakeParser(tr)
	if old {
		return parser.OldStart()
	} else {
		return parser.Start()
	}
}


func ParseFile(filename string, flags int, macros map[string]Macro, old bool) Node {
	if tr, result := LexFile(filename, flags, macros); !result {
		return nil
	} else {
		return ParseTokens(tr, old)
	}
}

func ParseString(code string, flags int, macros map[string]Macro, old bool) Node {
	output, good := finishLexing(Tokenize(code, ""), flags, macros)
	if good {
		return ParseTokens(output, old)
	}
	return nil
}

func ParseExpression(code string, flags int, macros map[string]Macro, old bool) Expr {
	output, good := finishLexing(Tokenize(code, ""), flags, macros)
	if good {
		parser := MakeParser(output)
		if old {
			return parser.OldMainExpr()
		} else {
			return parser.MainExpr()
		}
	}
	return nil
}

func ParseStatement(code string, flags int, macros map[string]Macro, old bool) Stmt {
	output, good := finishLexing(Tokenize(code, ""), flags, macros)
	if good {
		parser := MakeParser(output)
		if old {
			return nil
		} else {
			return parser.Statement()
		}
	}
	return nil
}

func EvalExpression(code string, flags int, macros map[string]Macro, old bool) TypeAndVal {
	output, good := finishLexing(Tokenize(code, ""), flags, macros)
	if good {
		parser := MakeParser(output)
		expr_node := Ternary[Expr](old, parser.OldMainExpr(), parser.MainExpr())
		interp := MakeInterpreter(parser)
		return interp.EvalExpr(expr_node)
	}
	return VoidTypeAndVal{}
}


func Ternary[T any](cond bool, a, b T) T {
	if cond {
		return a
	}
	return b
}