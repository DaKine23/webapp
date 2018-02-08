package oauth2

import (
	"reflect"
	"testing"
)

func TestParsePredefinedTokens(t *testing.T) {
	probes := []struct {
		expr   string
		result map[string]string
	}{
		{
			expr: "user-1=token-1,user-2=token-2",
			result: map[string]string{
				"user-1": "token-1",
				"user-2": "token-2",
			},
		},
		{
			expr: "user-1=token-1,user-2token-2",
			result: map[string]string{
				"user-1": "token-1",
			},
		},
		{
			expr:   "user-1token-1,user-2token-2",
			result: map[string]string{},
		},
	}

	for _, probe := range probes {
		if result := parsePredefinedTokens(probe.expr); !reflect.DeepEqual(result, probe.result) {
			t.Fatalf("Expected: %s, Got: %s", probe.result, result)
		}
	}
}
