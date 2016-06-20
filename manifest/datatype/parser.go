package datatype

import (
	"errors"
	"fmt"
	"github.com/chemikadze/gonomi/manifest/parsing"
)

func Parse(repr string) (DataType, error) {
	r := NewTokenReader(repr)
	return parseFromTokens(&r, []TokenType{})
}

func parseFromTokens(r *TokenReader, stopTokens []TokenType) (DataType, error) {
	tokenType, value := r.Read()
	if tokenType == TOKEN_ERROR {
		return nil, errors.New(fmt.Sprintf("Unexpected token: %s", value))
	}
	for _, stopToken := range stopTokens {
		if tokenType == stopToken {
			return nil, nil
		}
	}
	if tokenType != TOKEN_ALPHANUM {
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
		if tokenType != TOKEN_OPEN_BRK {
			return nil, errors.New(fmt.Sprintf("Unexpected token: %s", value))
		}
		typeParameter, err := parseFromTokens(r, []TokenType{TOKEN_CLOSING_BRK})
		if err != nil {
			return nil, err
		}
		tokenType, _ = r.Read()
		if tokenType != TOKEN_CLOSING_BRK {
			return nil, errors.New(fmt.Sprintf("Unexpected token: %s", value))
		}
		return List{typeParameter}, err
	case "map":
		tokenType, _ := r.Read()
		if tokenType != TOKEN_OPEN_BRK {
			return nil, errors.New(fmt.Sprintf("Unexpected token: %s", value))
		}
		keyType, err := parseFromTokens(r, []TokenType{TOKEN_CLOSING_BRK})
		if err != nil {
			return nil, err
		}
		tokenType, _ = r.Read()
		if tokenType != TOKEN_COMMA {
			return nil, errors.New(fmt.Sprintf("Unexpected token: %s", value))
		}
		valueType, err := parseFromTokens(r, []TokenType{TOKEN_CLOSING_BRK})
		if err != nil {
			return nil, err
		}
		tokenType, _ = r.Read()
		if tokenType != TOKEN_CLOSING_BRK {
			return nil, errors.New(fmt.Sprintf("Unexpected token: %s", value))
		}
		return Map{keyType, valueType}, err
	case "record":
		tokenType, _ := r.Read()
		if tokenType != TOKEN_OPEN_BRK {
			return nil, errors.New(fmt.Sprintf("Unexpected token: %s", value))
		}
		fields, err := ParseRecordBodyFromTokens(r, []TokenType{TOKEN_CLOSING_BRK})
		if err != nil {
			return nil, err
		}
		return Record{fields}, nil
	}
	return DataType(nil), parsing.ManifestError{"Unknown type " + value, 0, 0}
}

func ParseRecordBodyFromTokens(r *TokenReader, stopTokens []TokenType) (map[string]DataType, error) {
	fields := map[string]DataType{}
loop:
	for {
		// read key type
		valueType, err := parseFromTokens(r, stopTokens)
		if err != nil {
			return nil, err
		}
		// empty record detected
		if valueType == nil {
			return nil, nil
		}
		// read value
		tokenType, value := r.Read()
		if tokenType != TOKEN_ALPHANUM {
			return nil, errors.New(fmt.Sprintf("Unexpected token: %s", value))
		}
		// save field
		fields[value] = valueType
		// next cycle deciding
		tokenType, _ = r.Read()
		switch tokenType {
		case TOKEN_CLOSING_BRK:
			break loop
		case TOKEN_COMMA: // to the next record
		default:
			for _, stopToken := range stopTokens {
				if tokenType == stopToken {
					break loop
				}
			}
			return nil, errors.New(fmt.Sprintf("Unexpected token: %s", value))
		}
	}
	return fields, nil
}
