package util

import (
	"fmt"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	translations "github.com/go-playground/validator/v10/translations/zh"
	"log"
	"reflect"
	"regexp"
)

var (
	validate *validator.Validate
	Trans    ut.Translator
)

func RegValid() {
	// 注册自定义的校验器错误消息
	Trans, _ = ut.New(zh.New()).GetTranslator("zh")
	validate = binding.Validator.Engine().(*validator.Validate)
	if err := translations.RegisterDefaultTranslations(validate, Trans); err != nil {
		fmt.Println("RegisterDefaultTranslations failed", err)
		return
	}

	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		if jsonTag := fld.Tag.Get("json"); jsonTag != "" {
			return jsonTag
		}
		if formTag := fld.Tag.Get("form"); formTag != "" {
			return formTag
		}
		return fld.Name
	})

	err := validate.RegisterValidation("phone", validatePhone)
	if err != nil {
		log.Printf("validate register failed: %v\n", err)
	}

	err = validate.RegisterValidation("varname", validateVarName)
	if err != nil {
		log.Printf("validate register failed: %v\n", err)
	}

	setTranslation(validate, Trans)
	setVarNameTranslation(validate, Trans)
}

// 校验手机号
func validatePhone(fl validator.FieldLevel) bool {
	phone := fl.Field().String()
	pattern := "^((13[0-9])|(14[5,7])|(15[0-3,5-9])|(17[0,3,5-8])|(18[0-9])|166|198|199|(147))\\d{8}$"
	regex, err := regexp.Compile(pattern)
	if err != nil {
		return false
	}
	isValid := regex.MatchString(phone)
	return isValid
}

// 中文翻译
func setTranslation(validate *validator.Validate, trans ut.Translator) {
	validate.RegisterTranslation(
		"phone",
		trans,
		func(ut ut.Translator) error {
			return ut.Add("phone", "{0}不是一个有效的手机号码！", true)
		},
		func(ut ut.Translator, fe validator.FieldError) string {
			return fmt.Sprintf("%s不是一个有效的手机号码！", fe.Field())
		},
	)
}

// 校验变量名规范
// 必须英文首字母开否,后面只能出现英文字母、数字、下划线和横行符号
func validateVarName(fl validator.FieldLevel) bool {
	indexName := fl.Field().String()
	pattern := "^[a-zA-Z][a-zA-Z0-9_-]*$"
	regex, err := regexp.Compile(pattern)
	if err != nil {
		return false
	}
	isValid := regex.MatchString(indexName)
	return isValid
}

// 变量名规范中文翻译
func setVarNameTranslation(validate *validator.Validate, trans ut.Translator) {
	validate.RegisterTranslation(
		"varname",
		trans,
		func(ut ut.Translator) error {
			return ut.Add("varname", "{0}不是一个有效的变量名，请按照^[a-zA-Z][a-zA-Z0-9_-]*$规则进行输入", true)
		},
		func(ut ut.Translator, fe validator.FieldError) string {
			return fmt.Sprintf("%s不是一个有效的变量名，请按照^[a-zA-Z][a-zA-Z0-9_-]*$规则进行输入", fe.Field())
		},
	)
}
