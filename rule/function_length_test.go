package rule_test

import (
	"errors"
	"testing"

	"github.com/mgechev/revive/lint"
	"github.com/mgechev/revive/rule"
)

func TestFunctionLength_Configure(t *testing.T) {
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
			name:      "int64 arguments, as a TOML configuration yields them",
			arguments: lint.Arguments{int64(10), int64(20)},
		},
		{
			name:      "int arguments, as a YAML configuration yields them",
			arguments: lint.Arguments{10, 20},
		},
		{
			name:      "float64 arguments, as a JSON configuration yields them",
			arguments: lint.Arguments{float64(10), float64(20)},
		},
		{
			name:      "non-integral float max statements",
			arguments: lint.Arguments{10.5, 20},
			wantErr:   errors.New(`invalid configuration value for max statements in "function-length" rule; need an integer but got float64`),
		},
		{
			name:      "non-integral float max lines",
			arguments: lint.Arguments{10, 20.5},
			wantErr:   errors.New(`invalid configuration value for max lines in "function-length" rule; need an integer but got float64`),
		},
		{
			name:      "string max statements",
			arguments: lint.Arguments{"10", 20},
			wantErr:   errors.New(`invalid configuration value for max statements in "function-length" rule; need an integer but got string`),
		},
		{
			name:      "string max lines",
			arguments: lint.Arguments{10, "20"},
			wantErr:   errors.New(`invalid configuration value for max lines in "function-length" rule; need an integer but got string`),
		},
		{
			name:      "negative max statements",
			arguments: lint.Arguments{-1, 20},
			wantErr:   errors.New(`the configuration value for max statements in "function-length" rule cannot be negative, got -1`),
		},
		{
			name:      "negative max lines",
			arguments: lint.Arguments{10, -1},
			wantErr:   errors.New(`the configuration value for max lines in "function-length" rule cannot be negative, got -1`),
		},
		{
			name:      "too few arguments",
			arguments: lint.Arguments{10},
			wantErr:   errors.New(`invalid configuration for "function-length" rule, expected 2 arguments but got 1`),
		},
		{
			name:      "too many arguments",
			arguments: lint.Arguments{10, 20, 30},
			wantErr:   errors.New(`invalid configuration for "function-length" rule, expected 2 arguments but got 3`),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var r rule.FunctionLength

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
