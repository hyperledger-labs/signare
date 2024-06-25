package validators

import (
	"github.com/hyperledger-labs/signare/app/pkg/entities/address"

	"github.com/asaskevich/govalidator"
)

func setAddressValidator() {
	govalidator.CustomTypeTagMap.Set("address", func(i any, o any) bool {
		addr, ok := i.(address.Address)
		if ok {
			return !addr.IsEmpty()
		}
		addresses, ok := i.([]address.Address)
		if !ok {
			return false
		}
		for _, a := range addresses {
			if a.IsEmpty() {
				return false
			}
		}
		return true
	})
}
