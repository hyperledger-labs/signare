package validators

import (
	"github.com/asaskevich/govalidator"
)

func setPositiveValidator() {
	govalidator.TagMap["positive"] = func(str string) bool {
		value, err := govalidator.ToInt(str)
		if err != nil {
			return false
		}
		return value > 0
	}
}
