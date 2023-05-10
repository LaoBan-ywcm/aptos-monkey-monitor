package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/robfig/cron/v3"

	"github.com/aptos-monkey-monitor/env"
	"github.com/aptos-monkey-monitor/pkg/notice"
	"github.com/aptos-monkey-monitor/pkg/snype"
)

var LP = 0.0
var TAT = ""

func main() {

	if err := env.Init(); err != nil {
		fmt.Println(err.Error())
		return
	}

	c := cron.New(cron.WithChain())
	c.AddFunc("@every 3s", getMonkeyLowPrice)
	c.AddFunc("@every 30m", updateTAT)
	c.Start()

	done := make(chan bool)
	<-done
}

func getMonkeyLowPrice() {
	monkeyData, err := snype.GetMonkeyData()
	if err != nil {
		fmt.Println(err)
	}
	timeStr := time.Now().Format("2006-01-02 15:04:05")

	// 地板价变化
	if LP != monkeyData.Price {
		LP = monkeyData.Price
		noticeApp(monkeyData)
	}

	// 第一次运行
	if LP == 0 {
		LP = monkeyData.Price
	}

	fmt.Printf("【%s】【%s】 %f apt    %f usd\n", timeStr, monkeyData.Name, monkeyData.Price, monkeyData.UsdPrice)
}

func updateTAT() {
	tat, err := notice.GetTenantAccessToken()
	if err != nil {
		return
	}

	TAT = tat

	timeStr := time.Now().Format("2006-01-02 15:04:05")
	fmt.Printf("【%s】update tenant_access_token: %s\n", timeStr, tat)
}

func noticeApp(monkeyData *snype.MonkeyLowPriceBo) {
	noticeContent := []notice.RichText{
		{
			Tag:   "text",
			Text:  fmt.Sprintf("【NFT】 %s \n", monkeyData.Name),
			Style: []string{"bold"},
		},
		{
			Tag:  "text",
			Text: fmt.Sprintf("【APT】 %.2f \n", monkeyData.Price),
		},
		{
			Tag:  "text",
			Text: fmt.Sprintf("【USD】 %.2f\n", monkeyData.UsdPrice),
		},
	}

	// 发送图片
	imageKey, err := notice.GetImageKey(monkeyData.ImageUrl, TAT)
	if err == nil {
		noticeContent = append([]notice.RichText{
			{
				Tag:      "img",
				ImageKey: imageKey,
			},
		}, noticeContent...)
	}

	// 添加list链接
	if monkeyData.ListUrl != "" {
		noticeContent = append(noticeContent, notice.RichText{
			Tag:  "a",
			Href: monkeyData.ListUrl,
			Text: fmt.Sprintf("购买链接\n"),
		})
	}

	noticeData := notice.FeishuNotice{}
	noticeData.MsgType = "post"
	noticeData.Content.Post.ZhCn.Title = "【Aptos Monkeys】地板价变化通知"
	noticeData.Content.Post.ZhCn.Content = make([][]notice.RichText, 1)
	noticeData.Content.Post.ZhCn.Content[0] = noticeContent

	_, err = notice.Feishu(noticeData)
	if err != nil {
		fmt.Println(err)
	}
}

func prettyPrint(i interface{}) string {
	s, _ := json.MarshalIndent(i, "", "\t")
	return string(s)
}
