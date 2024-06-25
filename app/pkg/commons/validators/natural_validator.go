package validators

import (
	"github.com/asaskevich/govalidator"
)

func setNaturalValidator() {
	govalidator.TagMap["natural"] = func(str string) bool {
		value, err := govalidator.ToFloat(str)
		if err != nil {
			return false
		}
		if value == 0 {
			return true
		}
		return govalidator.IsNatural(value)
	}
}
