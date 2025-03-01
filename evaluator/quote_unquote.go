package evaluator

import (
	"fmt"

	"github.com/igormichalak/monkey/ast"
	"github.com/igormichalak/monkey/object"
	"github.com/igormichalak/monkey/token"
)

func quote(env *object.Environment, node ast.Node) object.Object {
	node = evalUnquoteCalls(env, node)
	return &object.Quote{Node: node}
}

func evalUnquoteCalls(env *object.Environment, quoted ast.Node) ast.Node {
	return ast.Modify(quoted, func(node ast.Node) ast.Node {
		if !isUnquoteCall(node) {
			return node
		}

		call, ok := node.(*ast.CallExpression)

		if ok && len(call.Arguments) == 1 {
			unquoted := Eval(env, call.Arguments[0])
			return convertObjectToASTNode(unquoted)
		} else {
			return node
		}
	})
}

func isUnquoteCall(node ast.Node) bool {
	callExpression, ok := node.(*ast.CallExpression)
	if !ok {
		return false
	}
	return callExpression.Function.TokenLiteral() == "unquote"
}

func convertObjectToASTNode(obj object.Object) ast.Node {
	switch obj := obj.(type) {
	case *object.Integer:
		t := token.Token{
			Type:    token.INT,
			Literal: fmt.Sprintf("%d", obj.Value),
		}
		return &ast.IntegerLiteral{Token: t, Value: obj.Value}
	case *object.Boolean:
		var t token.Token
		if obj.Value {
			t = token.Token{Type: token.TRUE, Literal: "true"}
		} else {
			t = token.Token{Type: token.FALSE, Literal: "false"}
		}
		return &ast.Boolean{Token: t, Value: obj.Value}
	case *object.Quote:
		return obj.Node
	default:
		return nil
	}
}
