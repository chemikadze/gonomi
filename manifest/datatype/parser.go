package datatype

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"github.com/chemikadze/gonomi/manifest/parsing"
	"io"
	"strings"
)

type tokenType int

const (
	ALPHANUM = iota
	OPEN_BRK
	CLOSING_BRK
	COMMA
	ERROR
	EOF
)

type tokenReader struct {
	reader *bufio.Reader
}

func isSpace(r rune) bool {
	return r == ' ' || r == '\t' || r == '\n'
}

func isAlphanumStart(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || r == '_'
}

func isAlphanum(r rune) bool {
	return (r >= '0' && r <= '9') || isAlphanumStart(r)
}

func (t *tokenReader) Read() (tokenType, string) {
	for {
		r, _, err := t.reader.ReadRune()
		if err == io.EOF {
			return EOF, ""
		}
		if err != nil {
			return ERROR, ""
		}
		if r == '>' {
			return CLOSING_BRK, ">"
		} else if r == '<' {
			return OPEN_BRK, "<"
		} else if r == ',' {
			return COMMA, ","
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
					return ERROR, ""
				}
				if !isAlphanum(r) {
					t.reader.UnreadRune()
					break
				}
				acc.WriteRune(r)
			}
			return ALPHANUM, acc.String()
		} else {
			return ERROR, string(r)
		}
	}
}

func Parse(repr string) (DataType, error) {
	r := tokenReader{bufio.NewReader(strings.NewReader(repr))}
	return ParseFromTokens(&r)
}

func ParseFromTokens(r *tokenReader) (DataType, error) {
	tokenType, value := r.Read()
	if tokenType == ERROR {
		return nil, errors.New(fmt.Sprintf("Unexpected token: %s", value))
	}
	if tokenType != ALPHANUM {
		return nil, errors.New(fmt.Sprintf("Unexpected token: %s", value))
	}
	switch value {
	case "int":
		return Int{}, nil
	case "bool":
		return Bool{}, nil
	case "string":
		return String{}, nil
	case "list":
		tokenType, _ := r.Read()
		if tokenType != OPEN_BRK {
			return nil, errors.New(fmt.Sprintf("Unexpected token: %s", value))
		}
		typeParameter, err := ParseFromTokens(r)
		if err != nil {
			return nil, err
		}
		tokenType, _ = r.Read()
		if tokenType != CLOSING_BRK {
			return nil, errors.New(fmt.Sprintf("Unexpected token: %s", value))
		}
		return List{typeParameter}, err
	case "map":
		tokenType, _ := r.Read()
		if tokenType != OPEN_BRK {
			return nil, errors.New(fmt.Sprintf("Unexpected token: %s", value))
		}
		keyType, err := ParseFromTokens(r)
		if err != nil {
			return nil, err
		}
		tokenType, _ = r.Read()
		if tokenType != COMMA {
			return nil, errors.New(fmt.Sprintf("Unexpected token: %s", value))
		}
		valueType, err := ParseFromTokens(r)
		if err != nil {
			return nil, err
		}
		tokenType, _ = r.Read()
		if tokenType != CLOSING_BRK {
			return nil, errors.New(fmt.Sprintf("Unexpected token: %s", value))
		}
		return Map{keyType, valueType}, err
	case "record":
		fields := map[string]DataType{}
		tokenType, _ := r.Read()
		if tokenType != OPEN_BRK {
			return nil, errors.New(fmt.Sprintf("Unexpected token: %s", value))
		}
		// read key-value types
	loop:
		for {
			// read key type
			valueType, err := ParseFromTokens(r)
			if err != nil {
				return nil, err
			}
			// read value
			tokenType, value := r.Read()
			if tokenType != ALPHANUM {
				return nil, errors.New(fmt.Sprintf("Unexpected token: %s", value))
			}
			// save field
			fields[value] = valueType
			// next cycle deciding
			tokenType, _ = r.Read()
			switch tokenType {
			case CLOSING_BRK: // finish processing
				break loop
			case COMMA: // to the next record
			default:
				return nil, errors.New(fmt.Sprintf("Unexpected token: %s", value))
			}
		}
		return Record{fields}, nil
	}
	return DataType(nil), parsing.ManifestError{"Unknown type " + value, 0, 0}
}
