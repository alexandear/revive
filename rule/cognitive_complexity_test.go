package rule_test

import (
	"errors"
	"testing"

	"github.com/mgechev/revive/lint"
	"github.com/mgechev/revive/rule"
)

func TestCognitiveComplexityRule_Configure(t *testing.T) {
	tests := []struct {
		name      string
		arguments lint.Arguments
		wantErr   error
	}{
		{
			name:      "no arguments",
			arguments: lint.Arguments{},
		},
		{
			name:      "int64 argument, as a TOML configuration yields it",
			arguments: lint.Arguments{int64(10)},
		},
		{
			name:      "int argument, as a YAML configuration yields it",
			arguments: lint.Arguments{10},
		},
		{
			name:      "float64 argument, as a JSON configuration yields it",
			arguments: lint.Arguments{float64(10)},
		},
		{
			name:      "non-integral float argument",
			arguments: lint.Arguments{10.5},
			wantErr:   errors.New("invalid argument type for cognitive-complexity, expected an integer, got float64"),
		},
		{
			name:      "string argument",
			arguments: lint.Arguments{"10"},
			wantErr:   errors.New("invalid argument type for cognitive-complexity, expected an integer, got string"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var r rule.CognitiveComplexityRule

			err := r.Configure(tt.arguments)

			if tt.wantErr == nil {
				if err != nil {
					t.Errorf("Configure() unexpected non-nil error %q", err)
				}
				return
			}
			if err == nil || err.Error() != tt.wantErr.Error() {
				t.Errorf("Configure() unexpected error: got %q, want %q", err, tt.wantErr)
			}
		})
	}
}
