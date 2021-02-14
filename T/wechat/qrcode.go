package wechat

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"time"

	"roo.bo/rlib"
)

type RGBColor struct {
	R int `json: "r"`
	G int `json: "g"`
	B int `json: "b"`
}

type qrCodeForm struct {
	Page      string `json:"page"`
	Width     int    `json:"width"`
	Scene     string `json:"scene"`
	AutoColor bool   `json:"auto_color"`
	// line_color 微信文档搞的不明不白
	// https://developers.weixin.qq.com/miniprogram/dev/api-backend/open-api/qr-code/wxacode.getUnlimited.html
	//LineColor *RGBColor `json:"line_color"`
	//
	IsHyaline bool `json:"is_hyaline"`
}

func NewQRCodeForm() *qrCodeForm {
	return &qrCodeForm{
		Page:  "/",
		Width: 430,
		Scene: "",
	}
}

var (
	WechatBaseUrl = "https://api.weixin.qq.com"
)

func CreateQRCode(ctx context.Context, accessToken string, form *qrCodeForm) ([]byte, error) {
	httpoption := &rlib.HTTPClientOption{
		CloseLog: true,
		Timeout:  time.Second * 3,
	}

	client := rlib.NewRooboHTTPClient(ctx, "vx_qr_code", httpoption, nil)
	urlstr := WechatBaseUrl + "/wxa/getwxacodeunlimit?access_token=" + accessToken

	postBody, _ := json.Marshal(form)
	rlib.Info(ctx, urlstr, string(postBody))
	resp, err := client.Post(urlstr, postBody, map[string]string{"Content-Type": "application/json"})
	if err != nil {
		return nil, err
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	type WXResp struct {
		Errcode int    `json:"errcode"`
		Errmsg  string `json:"errmsg"`
	}
	wxresp := &WXResp{}
	byts := []byte(data)
	if err = json.Unmarshal(byts, wxresp); err == nil && wxresp.Errcode != 0 {
		rlib.Info(ctx, "get response error from wechat:", wxresp)
		return nil, errors.New(wxresp.Errmsg)
	}
	return byts, nil
}
