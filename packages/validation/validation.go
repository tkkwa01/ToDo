package validation

import (
	"ToDo/packages/util"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"reflect"
)

var (
	validate   = validator.New()
	translator ut.Translator
	uni        *ut.UniversalTranslator
)

//func init() {
//	validate.UseActualTagWhenTranslate()
//
//	jp := ja.New()
//	uni = ut.New(jp, jp)
//	translator, _ = uni.GetTranslator("ja")
//
//	_ = jaTranslations.RegisterDefaultTranslations(validate, translator)
//}

func Validate() *validator.Validate {
	return validate
}

func Translator() ut.Translator {
	return translator
}

func Register(tag string, fn validator.Func, translation string, option *RegisterTransOption) {
	_ = validate.RegisterValidation(tag, fn)
	RegisterTrans(tag, translation, option)
}

type RegisterTransOption struct {
	CustomRegisTag  string
	CustomRegisFunc func(ut ut.Translator) (err error)
	CustomTransFunc func(ut ut.Translator, fe validator.FieldError) []string
}

func RegisterTrans(tag string, translation string, option *RegisterTransOption) {
	regisTag := tag
	if option != nil && option.CustomRegisTag != "" {
		regisTag = option.CustomRegisTag
	}

	registrationFunc := func(ut ut.Translator) (err error) {
		if err = ut.Add(regisTag, translation, true); err != nil {
			panic(err)
		}

		if option != nil && option.CustomRegisFunc != nil {
			if err = option.CustomRegisFunc(ut); err != nil {
				panic(err)
			}
		}

		return
	}

	translateFunc := func(ut ut.Translator, fe validator.FieldError) string {
		params := []string{fe.Field()}

		if option != nil && option.CustomTransFunc != nil {
			params = append(params, option.CustomTransFunc(ut, fe)...)
		}

		t, err := ut.T(fe.ActualTag(), params...)
		if err != nil {
			return "入力された値が正しくありません。"
		}
		return t
	}
	_ = validate.RegisterTranslation(tag, translator, registrationFunc, translateFunc)
}

func RegisterFieldTrans(values map[string]string) {
	validate.RegisterTagNameFunc(
		func(fld reflect.StructField) string {
			if value, ok := values[fld.Name]; ok {
				return value
			}
			return util.SnakeCase(fld.Name)
		},
	)
}
