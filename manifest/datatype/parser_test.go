package datatype

import (
	"bufio"
	"reflect"
	"strings"
	"testing"
)

type testValue struct {
	tokenType tokenType
	value     string
}

func checkTokens(t *testing.T, input string, expectedResult []testValue) {
	r := tokenReader{bufio.NewReader(strings.NewReader(input))}
	if len(expectedResult) == 0 || expectedResult[len(expectedResult)-1].tokenType != EOF {
		expectedResult = append(expectedResult, testValue{EOF, ""})
	}
	for _, expected := range expectedResult {
		tokenType, value := r.Read()
		if tokenType != expected.tokenType {
			t.Errorf("Expected token type %v got %v", expected.tokenType, tokenType)
		}
		if value != expected.value {
			t.Errorf("Expected token type %v got %v", expected.value, value)
		}
	}
}

func TestTokens(t *testing.T) {
	checkTokens(t, "", []testValue{})
	checkTokens(t, "", []testValue{testValue{EOF, ""}})
	checkTokens(t, ",", []testValue{testValue{COMMA, ","}})
	checkTokens(t, "<", []testValue{testValue{OPEN_BRK, "<"}})
	checkTokens(t, ">", []testValue{testValue{CLOSING_BRK, ">"}})
	checkTokens(t, "test", []testValue{testValue{ALPHANUM, "test"}})
	checkTokens(t, "0", []testValue{testValue{ERROR, "0"}})
	checkTokens(t, "tes_t0", []testValue{testValue{ALPHANUM, "tes_t0"}})
	checkTokens(t, "list<string>", []testValue{
		testValue{ALPHANUM, "list"},
		testValue{OPEN_BRK, "<"},
		testValue{ALPHANUM, "string"},
		testValue{CLOSING_BRK, ">"},
	})
	checkTokens(t, "list< list <string> >", []testValue{
		testValue{ALPHANUM, "list"},
		testValue{OPEN_BRK, "<"},
		testValue{ALPHANUM, "list"},
		testValue{OPEN_BRK, "<"},
		testValue{ALPHANUM, "string"},
		testValue{CLOSING_BRK, ">"},
		testValue{CLOSING_BRK, ">"},
	})
	checkTokens(t, "list<map<string,string>>", []testValue{
		testValue{ALPHANUM, "list"},
		testValue{OPEN_BRK, "<"},
		testValue{ALPHANUM, "map"},
		testValue{OPEN_BRK, "<"},
		testValue{ALPHANUM, "string"},
		testValue{COMMA, ","},
		testValue{ALPHANUM, "string"},
		testValue{CLOSING_BRK, ">"},
		testValue{CLOSING_BRK, ">"},
	})
}

func TestPinTypes(t *testing.T) {
	cases := map[string]DataType{
		"int":                         Int{},
		"string":                      String{},
		"bool":                        Bool{},
		"list<string>":                List{String{}},
		"list<list<string>>":          List{List{String{}}},
		"list< list <string> >":       List{List{String{}}},
		"map<string, string>":         Map{String{}, String{}},
		"record<string foo>":          Record{map[string]DataType{"foo": String{}}},
		"record<string foo, int bar>": Record{map[string]DataType{"foo": String{}, "bar": Int{}}},
	}
	for manifest, expected := range cases {
		parsed, err := Parse(manifest)
		if err != nil {
			t.Error(err)
		}
		if !reflect.DeepEqual(parsed, expected) {
			t.Errorf("\nParsed: %v\nExpect: %v", parsed, expected)
		}
	}
}
