package rules

import (
	"go/ast"
	"reflect"
	"testing"
)

func Test_complexityVisitor_Visit(t *testing.T) {
	type fields struct {
		complexity int
	}
	type args struct {
		n ast.Node
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   ast.Visitor
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &complexityVisitor{
				complexity: tt.fields.complexity,
			}
			if got := v.Visit(tt.args.n); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("complexityVisitor.Visit() = %v, want %v", got, tt.want)
			}
		})
	}
}
