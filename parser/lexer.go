package parser

import (
	"bufio"
	"go/token"
	"io"
	"unicode"

	"github.com/cznic/golex/lex"
	"github.com/z7zmey/php-parser/comment"
	t "github.com/z7zmey/php-parser/token"
)

// Allocate Character classes anywhere in [0x80, 0xFF].
const (
	classUnicodeLeter = iota + 0x80
	classUnicodeDigit
	classOther
)

type lexer struct {
	*lex.Lexer
	stateStack    []int
	phpDocComment string
	comments      []comment.Comment
}

func rune2Class(r rune) int {
	if r >= 0 && r < 0x80 { // Keep ASCII as it is.
		return int(r)
	}
	if unicode.IsLetter(r) {
		return classUnicodeLeter
	}
	if unicode.IsDigit(r) {
		return classUnicodeDigit
	}
	// return classOther
	return -1
}

func newLexer(src io.Reader, fName string) *lexer {
	file := token.NewFileSet().AddFile(fName, -1, 1<<31-1)
	lx, err := lex.New(file, bufio.NewReader(src), lex.RuneClass(rune2Class))
	if err != nil {
		panic(err)
	}
	return &lexer{lx, []int{0}, "", nil}
}

func (l *lexer) ungetN(n int) []byte {
	l.Unget(l.Lookahead())

	chars := l.Token()

	for i := 1; i <= n; i++ {
		char := chars[len(chars)-i]
		l.Unget(char)
	}

	buf := l.TokenBytes(nil)
	buf = buf[:len(buf)-n]

	return buf
}

func (l *lexer) pushState(state int) {
	l.stateStack = append(l.stateStack, state)
}

func (l *lexer) popState() {
	len := len(l.stateStack)
	if len <= 1 {
		return
	}

	l.stateStack = l.stateStack[:len-1]
}

func (l *lexer) begin(state int) {
	len := len(l.stateStack)
	l.stateStack = l.stateStack[:len-1]
	l.stateStack = append(l.stateStack, state)
}

func (l *lexer) getCurrentState() int {
	return l.stateStack[len(l.stateStack)-1]
}

func (l *lexer) newToken(tokenBytes []byte) t.Token {
	tokenBytesEnd := len(tokenBytes) - 1

	startLine := l.File.Line(l.First.Pos())
	lastChar := l.Token()[tokenBytesEnd]
	endLine := l.File.Line(lastChar.Pos())

	startPos := int(l.First.Pos())
	endPos := startPos + tokenBytesEnd

	return t.NewToken(tokenBytes, startLine, endLine, startPos, endPos).SetComments(l.comments)
}

func (l *lexer) addComment(c comment.Comment) {
	l.comments = append(l.comments, c)
}
