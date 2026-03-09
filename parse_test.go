package dispatch

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// Test parsePositionalArgs
func TestParsePositionalArgs(t *testing.T) {
	cmd := &Command{
		Command:      "test",
		ArgsRequired: 2,
		Args: []*Arg{
			{Name: "arg1", Type: String},
			{Name: "arg2", Type: Int},
		},
	}

	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{
			name:    "valid arguments",
			args:    []string{"value1", "123"},
			wantErr: false,
		},
		{
			name:    "missing argument",
			args:    []string{"value1"},
			wantErr: true,
		},
		{
			name:    "too many arguments",
			args:    []string{"value1", "123", "extra"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := &Data{}
			err := parsePositionalArgs(cmd, data, tt.args)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Len(t, data.ParsedArgs, len(cmd.Args))
			}
		})
	}
}

// Test parseComboArgs
func TestParseComboArgs(t *testing.T) {
	cmd := &Command{
		Command: "test",
		ArgumentCombos: [][]int{
			{0, 1}, // Combo 1: arg1 and arg2
		},
		Args: []*Arg{
			{Name: "arg1", Type: String},
			{Name: "arg2", Type: Int},
		},
	}

	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{
			name:    "valid combo",
			args:    []string{"value1", "123"},
			wantErr: false,
		},
		{
			name:    "missing argument in combo",
			args:    []string{"value1"},
			wantErr: true,
		},
		{
			name:    "invalid combo",
			args:    []string{"value1", "invalid"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := &Data{}
			err := parseComboArgs(cmd, data, tt.args)

			t.Logf("received error: %v", err)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Len(t, data.ParsedArgs, len(cmd.Args))
			}
		})
	}
}
