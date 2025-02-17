package question

import (
	"errors"
	"testing"
)

func TestNewQuestion(t *testing.T) {
	type args struct {
		name      string
		a         int32
		b         int32
		operation string
	}
	tests := []struct {
		name    string
		args    args
		want    float64
		wantErr error
	}{
		// TODO: Add test cases.
		{name: "1+1=2", args: args{a: 1, b: 1, operation: "+"}, want: 2, wantErr: nil},
		{name: "1^1=invalid operation", args: args{a: 1, b: 1, operation: "^"}, want: -1, wantErr: InvalidOperation},
		{name: "1/0=invalid input", args: args{a: 1, b: 0, operation: "/"}, want: -1, wantErr: InvalidInputDivisionBy0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewQuestion("", tt.args.a, tt.args.b, tt.args.operation)
			if err != nil && !errors.Is(err, tt.wantErr) {
				t.Errorf("NewQuestion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != nil && tt.want != got.answer {
				t.Errorf("NewQuestion().answer got = %v, want %v", got, tt.want)
			}
		})
	}
}
