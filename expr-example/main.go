package main

import (
	"bytes"
	"expr-example/menu"
	"fmt"
)

func main() {
	expression := `((not IsBUser()) and Appid in ["App1"] and SdkVersionGE("8.8.8")) or
(IsBUser() and Appid not in ["App2"] and SdkVersionL("8.8.8"))`

	writer := bytes.NewBuffer(make([]byte, 0, 0))
	menu.PrintExpression(expression, writer)
	fmt.Println(writer.String())

	env := menu.MenuEnv{
		SdkVersion: "7.9.8",
		Ucid: "Buser1",
		Appid: "App1",
	}

	res := menu.EvalExpr(expression, env)
	fmt.Println(res)
}
