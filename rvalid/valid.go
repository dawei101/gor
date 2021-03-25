package rvalid

import (
	"errors"
	"fmt"
	"regexp"
	"sync"

	"github.com/dawei101/gor/rhttp"
	valid "gopkg.in/go-playground/validator.v9"
)

var _valid *valid.Validate
var _v_mutex sync.Mutex

func init() {
	_valid = valid.New()
	regDefaultValidations()
}

var _validationMap = map[string]valid.Func{
	"sn":        re_match("^[1-9][0-9a-zA-Z]{15}$"),
	"appid":     re_match("^[a-zA-Z][0-9a-zA-Z]{15}$"),
	"cn_mobile": re_match("^1[3-9][0-9]{9}$"),
	"is_userid": re_match("^[a-zA-Z\\d]{2}:[\\da-zA-Z]+$"),
}

func re_match(reg string) func(valid.FieldLevel) bool {
	re := regexp.MustCompile(reg)
	return func(fl valid.FieldLevel) bool {
		return re.MatchString(fl.Field().String())
	}
}

//
//运行时动态加入验证规则
//
func ValidatorReg(tag string, f valid.Func) error {
	_, ok := _validationMap[tag]
	if ok {
		return errors.New("validation already existed!")
	}
	_v_mutex.Lock()
	defer _v_mutex.Unlock()
	_validationMap[tag] = f
	_valid.RegisterValidation(tag, f)
	return nil
}

func regDefaultValidations() {
	for k, v := range _validationMap {
		_valid.RegisterValidation(k, v)
	}
}

//
// 验证实例字段
//use for example:
//		import
//		user = &struct {
//			Appid string `validate:"required,appid"`
//			SN    string `validate:"required,sn"`
//		}{
//			SN:    "102102000100008F",
//			Appid: "TI1ZDE1NzJhZTEwO",
//		}
//
//		if err := rvalid.ValidField(user); err != nil {
//			fmt.Println(err)
//			rvalid.NewErrResp(-422, err.Error(), "")
//		}
//
func ValidField(v interface{}) error {
	err := _valid.Struct(v)
	if nil == err {
		return nil
	}
	verrs := err.(valid.ValidationErrors)
	errmsg := fmt.Sprintf("%s is not correct, notice on the '%s' tag", verrs[0].Field(), verrs[0].Tag())

	msg := ""
	for _, e := range verrs {
		msg += fmt.Sprintf("%+v", e) + ";"
	}
	return rhttp.NewRespErr(422, errmsg, "")
}
