package validator

import (
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

type Validator struct {
	// trans     ut.Translator
	validator *validator.Validate
}

func NewValidator() *Validator {
	validatorCustom := &Validator{}

	// en := en.New()
	// uni := ut.New(en, en)
	// trans, _ := uni.GetTranslator("en")

	v := validator.New()
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		var name string

		name = strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]

		if name == "" {
			name = strings.SplitN(fld.Tag.Get("query"), ",", 2)[0]
		}

		if name == "" {
			name = strings.SplitN(fld.Tag.Get("form"), ",", 2)[0]
		}

		if name == "" {
			name = strings.SplitN(fld.Tag.Get("params"), ",", 2)[0]
		}

		if name == "" {
			name = strings.SplitN(fld.Tag.Get("prop"), ",", 2)[0]
		}

		if name == "-" {
			return ""
		}

		return name
	})

	// en_translations.RegisterDefaultTranslations(v, trans)
	// if err := v.RegisterValidation("email_blacklist", isEmailBlacklist); err != nil {
	// 	log.Fatal().Err(err).Msg("Error while registering email_blacklist validator")
	// }
	// if err := v.RegisterValidation("strong_password", isStrongPassword); err != nil {
	// 	log.Fatal().Err(err).Msg("Error while registering strong_password validator")
	// }
	// if err := v.RegisterValidation("unique_in_slice", isUniqueInSlice); err != nil {
	// 	log.Fatal().Err(err).Msg("Error while registering unique validator")
	// }

	validatorCustom.validator = v
	// validatorCustom.trans = trans

	return validatorCustom
}

func (v *Validator) Validate(i any) error {
	return v.validator.Struct(i)
}
