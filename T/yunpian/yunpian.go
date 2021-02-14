package yunpian

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/dawei101/gor/rconfig"
	"github.com/dawei101/gor/rlog"

	ypclnt "github.com/yunpian/yunpian-go-sdk/sdk"
)

var config struct {
	ApiKey          string `yml:"apikey"`
	VoiceDisplayNum string `yml:"voiceDisplayNum"`
}

func init() {
	rconfig.ValueAssignTo("t.yunpian", &config, nil)
	if config.ApiKey == "" {
		rlog.Warning(context.Background(), "no apikey config for t.yunpian")
	}
}

func args2string(args map[string]string) string {
	tplvalues := []string{}
	for key, val := range args {
		tplvalues = append(tplvalues, fmt.Sprintf("#%s#=%s", key, url.QueryEscape(val)))
	}
	return strings.Join(tplvalues, "&amp;")
}

func SendSms(ctx context.Context, mobile string, tpl string, args map[string]string) bool {
	params := map[string]string{
		"apikey":    config.ApiKey,
		"mobile":    mobile,
		"tpl_id":    tpl,
		"tpl_value": args2string(args),
	}
	res := ypclnt.NewSms().TplSend(params)
	rlog.Info(ctx, fmt.Sprintf("send sms to: %s, result:%s", mobile, res))
	return res.Code == ypclnt.SUCC
}

func VoiceCallVerify(ctx context.Context, mobile string, code string) bool {
	params := map[string]string{
		"apikey": config.ApiKey,
		"mobile": mobile,
		"code":   code,
	}
	if config.VoiceDisplayNum != "" {
		params["display_num"] = config.VoiceDisplayNum
	}
	res := ypclnt.NewVoice().TplNotify(params)
	rlog.Info(ctx, fmt.Sprintf("send voice call to: %s, result:%s", mobile, res))
	return res.Code == ypclnt.SUCC
}
