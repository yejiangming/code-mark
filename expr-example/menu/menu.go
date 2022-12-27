package menu

import (
	"fmt"
	"github.com/antonmedv/expr"
	"github.com/antonmedv/expr/ast"
	"github.com/antonmedv/expr/parser"
	"io"
	"strconv"
	"strings"
)

type MenuEnv struct {
	Ucid string
	Appid string
	SdkVersion string
}

func (c MenuEnv) IsBUser() bool {
	if strings.HasPrefix(c.Ucid, "B") {
		return true
	}
	return false
}


// v1 > v2 1
// v1 = v2 0
// v1 < v2 -1
func sdkVersionCompare(v1, v2 string) int {
	v1Slice := strings.Split(v1, ".")
	if len(v1Slice) != 3 {
		panic("v1 format error")
	}
	v1ISlice := make([]int, 0, 3)
	for _, s := range v1Slice {
		i, err := strconv.Atoi(s)
		if err != nil {
			panic("v1 format error")
		}
		v1ISlice = append(v1ISlice, i)
	}

	v2Slice := strings.Split(v2, ".")
	if len(v2Slice) != 3 {
		panic("v2 format error")
	}
	v2ISlice := make([]int, 0, 3)
	for _, s := range v2Slice {
		i, err := strconv.Atoi(s)
		if err != nil {
			panic("v2 format error")
		}
		v2ISlice = append(v2ISlice, i)
	}

	for i := 0; i < 3; i++ {
		if v1ISlice[i] < v2ISlice[i] {
			return -1
		} else if v1ISlice[i] > v2ISlice[i] {
			return 1
		}
	}

	return 0

}

func (c MenuEnv) SdkVersionGE(targetVersion string) bool {
	return sdkVersionCompare(c.SdkVersion, targetVersion) >= 0
}

func (c MenuEnv) SdkVersionL(targetVersion string) bool {
	return sdkVersionCompare(c.SdkVersion, targetVersion) < 0
}

func EvalExpr(expression string, env MenuEnv) bool{
	// TODO expression vm cache
	vm, err := expr.Compile(expression, expr.Env(MenuEnv{}))
	if err != nil {
		panic(err)
	}
	valueI, err := expr.Run(vm, env)
	if err != nil {
		panic(err)
	}

	if value, ok := valueI.(bool); !ok {
		panic("expr value not bool")
	} else {
		return value
	}

}

func PrintExpression(expression string, writer io.Writer) {
	tree, err := parser.Parse(expression)
	if err != nil {
		panic(err)
	}

	PrintNode(0, tree.Node, writer)
}

func PrintTreeMessage(space int, message interface{}, writer io.Writer) {
	for i := 0; i < space; i++ {
		writer.Write([]byte(" "))
	}
	msg := fmt.Sprintf("%v\n", message)
	writer.Write([]byte(msg))

}

func PrintNode(space int, node interface{}, writer io.Writer) {
	switch node.(type) {
	case *ast.BinaryNode:
		entry := node.(*ast.BinaryNode)
		PrintNode(space+4, entry.Left, writer)
		PrintTreeMessage(space, entry.Operator, writer)
		PrintNode(space+4, entry.Right, writer)
	case *ast.UnaryNode:
		entry := node.(*ast.UnaryNode)
		PrintTreeMessage(space, entry.Operator, writer)
		PrintNode(space+4, entry.Node, writer)

	case *ast.FunctionNode:
		entry := node.(*ast.FunctionNode)
		PrintTreeMessage(space, entry.Name, writer)
		for _, n := range entry.Arguments {
			PrintNode(space+4, n, writer)
		}
	case *ast.IdentifierNode:
		entry := node.(*ast.IdentifierNode)
		PrintTreeMessage(space, entry.Value, writer)
	case *ast.ArrayNode:
		entry := node.(*ast.ArrayNode)
		PrintTreeMessage(space, "array", writer)
		for _, n := range entry.Nodes {
			PrintNode(space+4, n, writer)
		}
	case *ast.StringNode:
		entry := node.(*ast.StringNode)
		PrintTreeMessage(space, entry.Value, writer)
	case *ast.IntegerNode:
		entry := node.(*ast.IntegerNode)
		PrintTreeMessage(space, entry.Value, writer)
	default:
		panic("unknown type")
	}
}
