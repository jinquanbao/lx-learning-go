package main

import (
	"fmt"
	"testing"
)

func TestParseContent(t *testing.T) {
	templateStr := "租户【{{.TenantName}}】AI分析等待量:{{.WaitNum}},持续时间:{{.WaitTime}}秒;"
	param := make(map[string]interface{})
	param["TenantName"] = "租户1"
	param["WaitNum"] = 301
	param["WaitTime"] = 10801

	res, err := ParseContent(templateStr, param)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(res)
	t.Log(res)
}
