package notice

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	lark "github.com/larksuite/oapi-sdk-go/v3"
	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"

	"github.com/aptos-monkey-monitor/env"
	"github.com/aptos-monkey-monitor/pkg/request"
)

const feishuImagePostUrl = "https://open.feishu.cn/open-apis/im/v1/images"
const feishuTenantAccessTokenUrl = "https://open.feishu.cn/open-apis/auth/v3/tenant_access_token/internal"

type FeishuApp struct {
	AppID     string `json:"app_id"`
	AppSecret string `json:"app_secret"`
}

type RichText struct {
	Tag      string   `json:"tag"`
	Text     string   `json:"text,omitempty"`
	Href     string   `json:"href,omitempty"`
	Style    []string `json:"style.omitempty"`
	ImageKey string   `json:"image_key"`
}

type ZhcnBo struct {
	Title   string       `json:"title"`
	Content [][]RichText `json:"content"`
}

type FeishuNotice struct {
	MsgType string `json:"msg_type"`
	Content struct {
		Post struct {
			ZhCn ZhcnBo `json:"zh_cn"`
		} `json:"post"`
	} `json:"content"`
}

type FeishuImageBo struct {
	ImageType string `json:"image_type"`
	Image     []byte `josn:"image"`
}

type TenantAccessTokenBo struct {
	Code              int    `json:"code"`
	Expire            int    `json:"expire"`
	Msg               string `json:"msg"`
	TenantAccessToken string `json:"tenant_access_token"`
}

func Feishu(data FeishuNotice) ([]byte, error) {
	c := request.New(env.FEISHU_ROBOT_URL, http.MethodPost, nil)

	bodyBytes := new(bytes.Buffer)
	json.NewEncoder(bodyBytes).Encode(data)

	return c.Post(bodyBytes)
}

func GetTenantAccessToken() (string, error) {
	c := request.New(feishuTenantAccessTokenUrl, http.MethodGet, nil)

	payload := FeishuApp{
		AppID:     env.FEISHU_APP_ID,
		AppSecret: env.FEISHU_APP_SECRET,
	}

	payloadByte, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	resp, err := c.Post(strings.NewReader(string(payloadByte)))
	if err != nil {
		return "", err
	}

	var data TenantAccessTokenBo
	if err := json.Unmarshal(resp, &data); err != nil {
		return "", err
	}

	if data.Code != 0 {
		return "", fmt.Errorf("get TenantAccessToken faild")
	}

	return data.TenantAccessToken, nil
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func GetImageKey(imageUrl string, tat string) (string, error) {
	filename := filepath.Base(imageUrl)

	// 文件存在删除
	if _, err := os.Stat(filename); err == nil {
		os.Remove(filename)
	}

	// 从url先下载图片
	c := request.New(imageUrl, http.MethodGet, nil)
	resp, err := c.Get(nil)
	if err != nil {
		return "", err
	}

	fs, err := os.Create(filename)
	if err != nil {
		return "", err
	}

	// 保存文件到本地
	_, err = io.Copy(fs, strings.NewReader(string(resp)))
	if err != nil {
		return "", err
	}
	fs.Close()

	// 读取文件
	fh, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer fh.Close()

	client := lark.NewClient(env.FEISHU_APP_ID, env.FEISHU_APP_SECRET)
	req := larkim.NewCreateImageReqBuilder().Body(larkim.NewCreateImageReqBodyBuilder().ImageType("message").Image(fh).Build()).Build()

	newResp, err := client.Im.Image.Create(context.Background(), req, larkcore.WithTenantAccessToken(tat))
	if err != nil {
		return "", err
	}

	if !newResp.Success() {
		return "", fmt.Errorf("faild")
	}

	return *newResp.Data.ImageKey, nil
}
