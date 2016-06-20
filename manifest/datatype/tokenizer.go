package datatype

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"
)

type TokenType int

const (
	TOKEN_ALPHANUM = iota
	TOKEN_OPEN_BRK
	TOKEN_CLOSING_BRK
	TOKEN_COMMA
	TOKEN_ERROR
	TOKEN_EOF
	TOKEN_OPEN_BRS
	TOKEN_CLOSING_BRS
	TOKEN_NONE
)

type TokenReader struct {
	reader *bufio.Reader
}

func NewTokenReader(input string) TokenReader {
	return TokenReader{bufio.NewReader(strings.NewReader(input))}
}

func isSpace(r rune) bool {
	return r == ' ' || r == '\t' || r == '\n'
}

func isAlphanumStart(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || r == '_' || r == '-'
}

func isAlphanum(r rune) bool {
	return (r >= '0' && r <= '9') || isAlphanumStart(r)
}

func (t *TokenReader) Read() (TokenType, string) {
	for {
		r, _, err := t.reader.ReadRune()
		if err == io.EOF {
			return TOKEN_EOF, ""
		}
		if err != nil {
			return TOKEN_ERROR, ""
		}
		if r == '>' {
			return TOKEN_CLOSING_BRK, ">"
		} else if r == '<' {
			return TOKEN_OPEN_BRK, "<"
		} else if r == ',' {
			return TOKEN_COMMA, ","
		} else if r == '(' {
			return TOKEN_OPEN_BRS, "("
		} else if r == ')' {
			return TOKEN_CLOSING_BRS, ")"
		} else if isSpace(r) {
			// skip
		} else if isAlphanumStart(r) {
			var acc bytes.Buffer
			acc.WriteRune(r)
			for {
				r, _, err := t.reader.ReadRune()
				if err == io.EOF {
					break
				}
				if err != nil {
					return TOKEN_ERROR, ""
				}
				if !isAlphanum(r) {
					t.reader.UnreadRune()
					break
				}
				acc.WriteRune(r)
			}
			return TOKEN_ALPHANUM, acc.String()
		} else {
			return TOKEN_ERROR, string(r)
		}
	}
}

func (t *TokenReader) ReadAssertToken(next ...TokenType) error {
	ttype, token := t.Read()
	for _, candidate := range next {
		if ttype == candidate {
			return nil
		}
	}
	return errors.New(fmt.Sprintf("Unexpected token: %s", token))
}
