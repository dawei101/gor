package yuntongxun

import (
	"context"

	"github.com/cloopen/go-sms-sdk/cloopen"

	"github.com/dawei101/gor/rconfig"
	"github.com/dawei101/gor/rlog"
)

var ccfg *cloopen.Config

var config struct {
	Account string `yml:"account"`
	Token   string `yml:"token"`
	AppId   string `yml:"appId"`
}

func init() {
	rconfig.ValueAssignTo("t.yunpian", &config, nil)
	if config.Account == "" || config.Token == "" {
		rlog.Warning(context.Background(), "config is not correct for t.yuntongxun")
	}
	ccfg = cloopen.DefaultConfig().WithAPIAccount(config.Account).WithAPIToken(config.Token)
}

// args 参见：https://learnku.com/laravel/t/15438/rong-lian-cloud-communication-pit-sharing-about-template-data-settings
//
func SendSms(ctx context.Context, mobile string, tpl string, args []string) bool {
	sms := cloopen.NewJsonClient(ccfg).SMS()
	// 下发包体参数
	input := &cloopen.SendRequest{
		// 应用的APPID
		AppId: config.AppId,
		// 手机号码
		To: mobile,
		// 模版ID
		TemplateId: tpl,
		// 模版变量内容 非必填
		Datas: args,
	}
	// 下发
	_, err := sms.Send(input)
	if err != nil {
		rlog.Warning(ctx, "send sms to", mobile, "failed:", err.Error())
		return false
	}
	rlog.Warning(ctx, "send sms to", mobile, "ok")
	return true
}

const voiceUrl = "https://app.cloopen.com:8883"

// TODO
func VoiceCall(ctx context.Context, mobile string, tpl string, args map[string]string) bool {

	_ = "/2013-12-26/Accounts/" + config.Account + "/Calls/VoiceVerify?sig={SigParameter}"
	return false
}
