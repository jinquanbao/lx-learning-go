package main

import (
	"fmt"
)

type AiAnalyseStatisticsVO struct {
	TenantId   int    `json:"id"`         //租户Id
	TenantName string `json:"tenantName"` //租户名称
	WaitNum    int    `json:"waitNum"`    //等待数量
	WaitTime   int    `json:"waitTime"`   //等待持续时长
}

func main() {
	templateStr := "租户【{{.TenantName}}】AI分析等待量:{{.WaitNum}},持续时间:{{.WaitTime}}秒;"
	templateStr2 := "租户ID【{{.TenantId}}】;"

	list := make([]AiAnalyseStatisticsVO, 0)
	vo := AiAnalyseStatisticsVO{
		TenantId:   1,
		TenantName: "租户1",
		WaitNum:    301,
		WaitTime:   10801,
	}
	list = append(list, vo)
	vo2 := AiAnalyseStatisticsVO{
		TenantId:   1,
		TenantName: "租户2",
		WaitNum:    302,
		WaitTime:   10802,
	}
	list = append(list, vo2)

	for _, v := range list {
		res, err := ParseContent(templateStr, v)
		if err != nil {
			panic(err)
		}
		fmt.Println(res)
		res, err = ParseContent(templateStr2, v)
		if err != nil {
			panic(err)
		}
		fmt.Println(res)
	}
}
