package dgx

import v "github.com/go-playground/validator/v10"

const validatorTag = "valid"

var validator *v.Validate

func init() {
	validator = v.New()
	validator.SetTagName(validatorTag)
}
