package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	"github.com/go-playground/universal-translator"
	"gopkg.in/go-playground/validator.v9"
	en_translations "gopkg.in/go-playground/validator.v9/translations/en"
	zh_translations "gopkg.in/go-playground/validator.v9/translations/zh"
	"my_gateway/public"
	"reflect"
	"regexp"
	"strings"
)

//设置Translation
func TranslationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		//参照：https://github.com/go-playground/validator/blob/v9/_examples/translations/main.go

		//设置支持语言
		en := en.New()
		zh := zh.New()

		//设置国际化翻译器
		uni := ut.New(zh, zh, en)
		val := validator.New()

		//根据参数取翻译器实例
		locale := c.DefaultQuery("locale", "zh")
		trans, _ := uni.GetTranslator(locale)

		//翻译器注册到validator
		switch locale {
		case "en":
			en_translations.RegisterDefaultTranslations(val, trans)
			val.RegisterTagNameFunc(func(fld reflect.StructField) string {
				return fld.Tag.Get("en_comment")
			})
			break
		default:
			zh_translations.RegisterDefaultTranslations(val, trans)
			val.RegisterTagNameFunc(func(fld reflect.StructField) string {
				return fld.Tag.Get("comment")
			})

			//自定义验证方法
			//https://github.com/go-playground/validator/blob/v9/_examples/custom-validation/main.go
			val.RegisterValidation("is_valid_username", func(fl validator.FieldLevel) bool {
				return fl.Field().String() == "admin"
			})
			//自定义Translator
			//https://github.com/go-playground/validator/blob/v9/_examples/translations/main.go
			val.RegisterTranslation("is_valid_username", trans, func(ut ut.Translator) error {
				return ut.Add("is_valid_username", "{0} 填写不正确哦", true)
			}, func(ut ut.Translator, fe validator.FieldError) string {
				t, _ := ut.T("is_valid_username", fe.Field())
				return t
			})

			val.RegisterValidation("is_valid_service_name", func(fl validator.FieldLevel) bool {
				matched, _ := regexp.Match(`^[a-zA-Z0-9_]{6,128}$`, []byte(fl.Field().String()))
				return matched
			})
			val.RegisterTranslation("is_valid_service_name", trans, func(ut ut.Translator) error {
				return ut.Add("is_valid_service_name", "{0} format is wrong.", true)
			}, func(ut ut.Translator, fe validator.FieldError) string {
				t, _ := ut.T("is_valid_service_name", fe.Field())
				return t
			})

			val.RegisterValidation("is_valid_rule", func(fl validator.FieldLevel) bool {
				matched, _ := regexp.Match(`^\S+$`, []byte(fl.Field().String()))
				return matched
			})
			val.RegisterTranslation("is_valid_rule", trans, func(ut ut.Translator) error {
				return ut.Add("is_valid_rule", "{0} Must be a non-empty string!", true)
			}, func(ut ut.Translator, fe validator.FieldError) string {
				t, _ := ut.T("is_valid_rule", fe.Field())
				return t
			})

			val.RegisterValidation("is_valid_url_rewrite", func(fl validator.FieldLevel) bool {
				if fl.Field().String() == "" {
					return true
				}
				for _, rewrites := range strings.Split(fl.Field().String(), "\n") {
					if len(strings.Split(rewrites, " ")) != 2 {
						return false
					}
				}
				return true
			})
			val.RegisterTranslation("is_valid_url_rewrite", trans, func(ut ut.Translator) error {
				return ut.Add("is_valid_url_rewrite", "{0} Must be a non-empty string!", true)
			}, func(ut ut.Translator, fe validator.FieldError) string {
				t, _ := ut.T("is_valid_url_rewrite", fe.Field())
				return t
			})

			val.RegisterValidation("is_valid_header_transform", func(fl validator.FieldLevel) bool {
				if fl.Field().String() == "" {
					return true
				}
				for _, rewrites := range strings.Split(fl.Field().String(), "\n") {
					if len(strings.Split(rewrites, " ")) != 3 {
						return false
					}
				}
				return true
			})
			val.RegisterTranslation("is_valid_header_transform", trans, func(ut ut.Translator) error {
				return ut.Add("is_valid_header_transform", "{0} Must be a non-empty string!", true)
			}, func(ut ut.Translator, fe validator.FieldError) string {
				t, _ := ut.T("is_valid_header_transform", fe.Field())
				return t
			})

			val.RegisterValidation("is_valid_ip_list", func(fl validator.FieldLevel) bool {

				for _, ms := range strings.Split(fl.Field().String(), "\n") {
					matched, _ := regexp.Match(`^\S+\:\d+$`, []byte(ms))
					if !matched {
						return false
					}
				}
				return true
			})
			val.RegisterTranslation("is_valid_ip_list", trans, func(ut ut.Translator) error {
				return ut.Add("is_valid_ip_list", "{0} Must be a non-empty string!", true)
			}, func(ut ut.Translator, fe validator.FieldError) string {
				t, _ := ut.T("is_valid_ip_list", fe.Field())
				return t
			})

			val.RegisterValidation("is_valid_weight_list", func(fl validator.FieldLevel) bool {
				for _, ms := range strings.Split(fl.Field().String(), "\n") {
					matched, _ := regexp.Match(`^\d+$`, []byte(ms))
					if !matched {
						return false
					}
				}
				return true
			})
			val.RegisterTranslation("is_valid_weight_list", trans, func(ut ut.Translator) error {
				return ut.Add("is_valid_weight_list", "{0} Must be a non-empty string!", true)
			}, func(ut ut.Translator, fe validator.FieldError) string {
				t, _ := ut.T("is_valid_weight_list", fe.Field())
				return t
			})

			break
		}
		c.Set(public.TranslatorKey, trans)
		c.Set(public.ValidatorKey, val)
		c.Next()
	}
}
