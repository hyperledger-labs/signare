package pkcs11hsm

import (
	"github.com/hyperledger-labs/signare/app/pkg/signaturemanager"

	"github.com/miekg/pkcs11"
)

type Curve string

const (
	CurveDefault   Curve = "secp256k1"
	CurveSecp256k1 Curve = "secp256k1"
)

var pkcsErrTranslator = map[pkcs11.Error]*signaturemanager.Error{
	pkcs11.CKR_SLOT_ID_INVALID:              signaturemanager.NewInvalidSlotError(),
	pkcs11.CKR_PIN_INCORRECT:                signaturemanager.NewPinIncorrectError(),
	pkcs11.CKR_CRYPTOKI_ALREADY_INITIALIZED: signaturemanager.NewAlreadyInitializedError(),
}

// PKCS11HSMConnectionDetails configuration to connect to a specific slot in a softHSM instance.
type PKCS11HSMConnectionDetails struct {
	// Configuration configuration to connect to a softHSM instance.
	Configuration PKCS11HSMConfiguration
}

// PKCS11HSMConfiguration configuration to connect to a softHSM instance.
type PKCS11HSMConfiguration struct {
	Curve Curve
}
