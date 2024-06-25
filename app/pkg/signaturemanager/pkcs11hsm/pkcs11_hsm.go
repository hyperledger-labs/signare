// Package pkcs11hsm adds support for PKCS11 compatible HSMs by satisfying the DigitalSignatureManager interface.
package pkcs11hsm

import (
	"context"
	"encoding/asn1"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"fmt"
	"sort"
	"strconv"

	curves "github.com/btcsuite/btcd/btcec/v2"
	"github.com/miekg/pkcs11"

	"github.com/hyperledger-labs/signare/app/pkg/commons/logger"
	"github.com/hyperledger-labs/signare/app/pkg/commons/time"
	"github.com/hyperledger-labs/signare/app/pkg/entities/address"
	signererrors "github.com/hyperledger-labs/signare/app/pkg/internal/errors"
	"github.com/hyperledger-labs/signare/app/pkg/signaturemanager"
)

const (
	standard = "PKCS#11"
)

// PKCS11HSMSignatureManager implements the DigitalSignatureManager interface.
type PKCS11HSMSignatureManager struct {
	pkcsContext       *pkcs11.Ctx
	connectionDetails PKCS11HSMConnectionDetails
}

// PKCS11HSMSignatureManagerOptions defines options to create a new instance of PKCS11HSMSignatureManager.
type PKCS11HSMSignatureManagerOptions struct {
	PkcsContext *pkcs11.Ctx
}

var _ signaturemanager.DigitalSignatureManager = (*PKCS11HSMSignatureManager)(nil)

// ProvidePKCS11HSMSignatureManager creates a new instance of PKCS11HSMSignatureManager using the provided options, returning an error if it fails.
func ProvidePKCS11HSMSignatureManager(options PKCS11HSMSignatureManagerOptions) (*PKCS11HSMSignatureManager, error) {
	if options.PkcsContext == nil {
		return nil, signaturemanager.NewInvalidArgumentError().WithMessage("nil pkcs11 context")
	}

	return &PKCS11HSMSignatureManager{
		pkcsContext: options.PkcsContext,
		connectionDetails: PKCS11HSMConnectionDetails{
			Configuration: PKCS11HSMConfiguration{
				Curve: CurveDefault,
			},
		},
	}, nil
}

func (s *PKCS11HSMSignatureManager) GenerateKey(_ context.Context, input signaturemanager.GenerateKeyInput) (*signaturemanager.GenerateKeyOutput, error) {
	tracer := input.Tracer
	slot, err := strconv.ParseUint(input.Slot, 10, 32)
	if err != nil {
		return nil, signaturemanager.NewInvalidSlotError().WithMessage(fmt.Sprintf("invalid slot: '%s'", input.Slot))
	}
	tracer.AddProperty("slot", slot)
	tracer.AddProperty("standard", standard)
	session, err := s.pkcsContext.OpenSession(uint(slot), pkcs11.CKF_SERIAL_SESSION|pkcs11.CKF_RW_SESSION)
	if err != nil {
		return nil, toSignatureManagerErr(err, fmt.Sprintf("could not open PKCS11 session. Error: %v", err))
	}
	defer s.closeSession(tracer, session)

	err = s.pkcsContext.Login(session, pkcs11.CKU_USER, input.Pin)
	if err != nil {
		return nil, toSignatureManagerErr(err, fmt.Sprintf("failed to login with the PKCS11 session. Error: %v", err))
	}
	defer s.logOut(tracer, session)

	ecParams := s.getEllipticCurveParameters()

	timestamp := generateTimestampId()
	lb := base64.StdEncoding.EncodeToString(timestamp)
	publicKeyTemplate := []*pkcs11.Attribute{
		pkcs11.NewAttribute(pkcs11.CKA_KEY_TYPE, pkcs11.CKK_EC),
		pkcs11.NewAttribute(pkcs11.CKA_CLASS, pkcs11.CKO_PUBLIC_KEY),
		pkcs11.NewAttribute(pkcs11.CKA_TOKEN, true),
		pkcs11.NewAttribute(pkcs11.CKA_VERIFY, true),
		pkcs11.NewAttribute(pkcs11.CKA_EC_PARAMS, ecParams),
		pkcs11.NewAttribute(pkcs11.CKA_PRIVATE, false),
		pkcs11.NewAttribute(pkcs11.CKA_LABEL, lb),
		pkcs11.NewAttribute(pkcs11.CKA_ID, timestamp),
	}
	timestamp = generateTimestampId()
	lb = base64.StdEncoding.EncodeToString(timestamp)
	privateKeyTemplate := []*pkcs11.Attribute{
		pkcs11.NewAttribute(pkcs11.CKA_KEY_TYPE, pkcs11.CKK_EC),
		pkcs11.NewAttribute(pkcs11.CKA_CLASS, pkcs11.CKO_PRIVATE_KEY),
		pkcs11.NewAttribute(pkcs11.CKA_TOKEN, true),
		pkcs11.NewAttribute(pkcs11.CKA_SIGN, true),
		pkcs11.NewAttribute(pkcs11.CKA_DERIVE, false),
		pkcs11.NewAttribute(pkcs11.CKA_PRIVATE, true),
		pkcs11.NewAttribute(pkcs11.CKA_LABEL, lb),
		pkcs11.NewAttribute(pkcs11.CKA_ID, timestamp),
	}

	tracer.Debug("generating key pair")
	publicKeyHandle, privateKeyHandle, err := s.pkcsContext.GenerateKeyPair(session,
		[]*pkcs11.Mechanism{pkcs11.NewMechanism(pkcs11.CKM_EC_KEY_PAIR_GEN, nil)},
		publicKeyTemplate,
		privateKeyTemplate)
	if err != nil {
		return nil, signaturemanager.NewKeyGenerationError().WithMessage(fmt.Sprintf("error generating key: %v", err))
	}

	tracer.Debug("getting address from public key")
	addr, err := s.getAddress(session, publicKeyHandle)
	if err != nil {
		return nil, err
	}
	publicKeyLabel := calculatePublicKeyLabel(*addr)
	tracer.Debugf("setting public key label for address '%s'", addr.String())
	err = s.setLabel(session, publicKeyHandle, publicKeyLabel)
	if err != nil {
		tracer.Warn(fmt.Sprintf("failed to set the public key label for publicKeyHandle '%d'. Error: %v", publicKeyHandle, err))
	}
	privateKeyLabel := calculatePrivateKeyLabel(*addr)
	tracer.Debugf("setting private key label for address '%s'", addr.String())
	err = s.setLabel(session, privateKeyHandle, privateKeyLabel)
	if err != nil {
		tracer.Warn(fmt.Sprintf("failed to set the private key label for privateKeyHandle '%d'. Error: %v", privateKeyHandle, err))
	}
	return &signaturemanager.GenerateKeyOutput{
		Address: *addr,
	}, nil
}

func (s *PKCS11HSMSignatureManager) RemoveKey(_ context.Context, input signaturemanager.RemoveKeyInput) (*signaturemanager.RemoveKeyOutput, error) {
	tracer := input.Tracer
	tracer.AddProperty("address", input.Address.String())
	slot, err := strconv.ParseUint(input.Slot, 10, 32)
	if err != nil {
		return nil, signaturemanager.NewInvalidSlotError().WithMessage(fmt.Sprintf("invalid slot: '%s'", input.Slot))
	}
	tracer.AddProperty("slot", slot)
	tracer.AddProperty("standard", standard)
	session, err := s.pkcsContext.OpenSession(uint(slot), pkcs11.CKF_SERIAL_SESSION|pkcs11.CKF_RW_SESSION)
	if err != nil {
		return nil, toSignatureManagerErr(err, fmt.Sprintf("could not open PKCS11 session. Error: %v", err))
	}
	defer s.closeSession(tracer, session)

	err = s.pkcsContext.Login(session, pkcs11.CKU_USER, input.Pin)
	if err != nil {
		return nil, toSignatureManagerErr(err, fmt.Sprintf("could not login PKCS11 session. Error: %v", err))
	}
	defer s.logOut(tracer, session)

	// Private key
	tracer.Debug("removing private key")
	privateKeyLabel := calculatePrivateKeyLabel(input.Address)
	templatePrivate := []*pkcs11.Attribute{
		pkcs11.NewAttribute(pkcs11.CKA_CLASS, pkcs11.CKO_PRIVATE_KEY),
		pkcs11.NewAttribute(pkcs11.CKA_LABEL, privateKeyLabel),
	}
	privateKeyObjects, err := s.findObjects(session, templatePrivate)
	if err != nil {
		return nil, signererrors.InternalFromErr(err).WithMessage(fmt.Sprintf("error finding PKCS11 attributes: %v", err))
	}

	// Public key
	tracer.Trace("removing public key")
	publicKeyLabel := calculatePublicKeyLabel(input.Address)
	templatePublic := []*pkcs11.Attribute{
		pkcs11.NewAttribute(pkcs11.CKA_CLASS, pkcs11.CKO_PUBLIC_KEY),
		pkcs11.NewAttribute(pkcs11.CKA_LABEL, publicKeyLabel),
	}
	publicKeyObjects, err := s.findObjects(session, templatePublic)
	if err != nil {
		return nil, signererrors.InternalFromErr(err).WithMessage(fmt.Sprintf("error finding PKCS11 attributes: %v", err))
	}

	if len(privateKeyObjects) == 0 && len(publicKeyObjects) == 0 {
		return nil, signaturemanager.NewNotFoundError().WithMessage(fmt.Sprintf("key pair not found for address '%s'", input.Address))
	}

	for _, object := range privateKeyObjects {
		err = s.pkcsContext.DestroyObject(session, object)
		if err != nil {
			return nil, toSignatureManagerErr(err, "call to PKCS11 destroy object function failed for the private key")
		}
	}

	for _, object := range publicKeyObjects {
		err = s.pkcsContext.DestroyObject(session, object)
		if err != nil {
			return nil, toSignatureManagerErr(err, "call to PKCS11 destroy object function failed for the public key")
		}
	}
	return &signaturemanager.RemoveKeyOutput{}, nil
}

func (s *PKCS11HSMSignatureManager) ListKeys(_ context.Context, input signaturemanager.ListKeysInput) (*signaturemanager.ListKeysOutput, error) {
	tracer := input.Tracer
	slot, err := strconv.ParseUint(input.Slot, 10, 32)
	if err != nil {
		return nil, signaturemanager.NewInvalidSlotError()
	}

	tracer.AddProperty("slot", slot)
	tracer.AddProperty("standard", standard)
	session, err := s.pkcsContext.OpenSession(uint(slot), pkcs11.CKF_SERIAL_SESSION)
	if err != nil {
		return nil, toSignatureManagerErr(err, "error opening PKCS11 session")
	}
	defer s.closeSession(tracer, session)

	err = s.pkcsContext.Login(session, pkcs11.CKU_USER, input.Pin)
	if err != nil {
		return nil, toSignatureManagerErr(err, "error logging in with the PKCS11 session")
	}
	defer s.logOut(tracer, session)

	pubKeyTemplate := []*pkcs11.Attribute{
		pkcs11.NewAttribute(pkcs11.CKA_CLASS, pkcs11.CKO_PUBLIC_KEY),
	}

	// First find the list of all public objects.
	objects, err := s.findObjects(session, pubKeyTemplate)
	if err != nil {
		return nil, err
	}

	addresses := make([]addressTime, 0)

	// For each of the public objects, determine the address and check that it matches the label
	for _, o := range objects {
		label, getLabelErr := s.getLabel(session, o)
		if getLabelErr != nil {
			continue
		}
		addr, _ := s.getAddress(session, o)
		toCompare := calculatePublicKeyLabel(*addr)
		if toCompare != *label {
			continue
		}
		t, getTimestampErr := s.getTimestamp(session, o)
		if getTimestampErr != nil {
			continue
		}
		addresses = append(addresses, addressTime{*addr, t})
	}
	sort.Sort(ByTime(addresses))
	result := ByTime(addresses).Addresses()
	return &signaturemanager.ListKeysOutput{
		Items: result,
	}, nil
}

func (s *PKCS11HSMSignatureManager) Sign(ctx context.Context, input signaturemanager.SignInput) (*signaturemanager.SignOutput, error) {
	tracer := input.Tracer
	slot, err := strconv.ParseUint(input.Slot, 10, 32)
	if err != nil {
		return nil, signaturemanager.NewInvalidSlotError().WithMessage(fmt.Sprintf("invalid slot: '%s'", input.Slot))
	}

	tracer.AddProperty("slot", slot)
	tracer.AddProperty("standard", standard)
	tracer.Debug("signing transaction")

	sig, err := s.sign(ctx, tracer, uint(slot), input.Pin, input.Data[:], input.From)
	if err != nil {

		return nil, err
	}
	return &signaturemanager.SignOutput{
		Signature: sig,
	}, nil
}

func (s *PKCS11HSMSignatureManager) Close(_ context.Context, _ signaturemanager.CloseInput) (*signaturemanager.CloseOutput, error) {
	err := s.pkcsContext.Finalize()
	if err != nil {
		return nil, toSignatureManagerErr(err, "error calling PKCS11 finalize")
	}

	return &signaturemanager.CloseOutput{}, nil
}

func (s *PKCS11HSMSignatureManager) Open(_ context.Context, _ signaturemanager.OpenInput) (*signaturemanager.OpenOutput, error) {
	err := s.pkcsContext.Initialize()
	if err != nil {
		return nil, toSignatureManagerErr(err, "error calling PKCS11 finalize")
	}

	return &signaturemanager.OpenOutput{}, nil
}

func (s *PKCS11HSMSignatureManager) IsAlive(_ context.Context, input signaturemanager.IsAliveInput) (*signaturemanager.IsAliveOutput, error) {
	tracer := input.Tracer
	slot, err := strconv.ParseUint(input.Slot, 10, 32)
	if err != nil {
		return nil, signaturemanager.NewInvalidSlotError().WithMessage(fmt.Sprintf("invalid slot: '%s'", input.Slot))
	}

	tracer.AddProperty("slot", slot)
	tracer.AddProperty("standard", standard)
	session, openSessionErr := s.pkcsContext.OpenSession(uint(slot), pkcs11.CKF_SERIAL_SESSION|pkcs11.CKF_RW_SESSION)
	if openSessionErr != nil {
		return nil, toSignatureManagerErr(openSessionErr, "error opening PKCS11 session")
	}
	defer s.closeSession(tracer, session)

	err = s.pkcsContext.Login(session, pkcs11.CKU_USER, input.Pin)
	if err != nil {
		return nil, toSignatureManagerErr(err, "error logging in with the PKCS11 session")
	}
	defer s.logOut(tracer, session)

	return &signaturemanager.IsAliveOutput{
		IsAlive: true,
	}, nil
}

func (s *PKCS11HSMSignatureManager) sign(_ context.Context, tracer logger.Tracer, slot uint, pin string, payloadToSign []byte, address address.Address) ([]byte, error) {
	tracer.AddProperty("slot", slot)
	tracer.AddProperty("address", address.String())
	tracer.AddProperty("standard", standard)
	session, err := s.pkcsContext.OpenSession(slot, pkcs11.CKF_SERIAL_SESSION|pkcs11.CKF_RW_SESSION)
	if err != nil {
		return nil, toSignatureManagerErr(err, "error opening PKCS11 session")
	}
	defer s.closeSession(tracer, session)

	err = s.pkcsContext.Login(session, pkcs11.CKU_USER, pin)
	if err != nil {
		return nil, toSignatureManagerErr(err, "error logging in with the PKCS11 session")
	}
	defer s.logOut(tracer, session)

	tracer.Debug("retrieving private key")
	privateKeyLabel := calculatePrivateKeyLabel(address)
	template := []*pkcs11.Attribute{
		pkcs11.NewAttribute(pkcs11.CKA_CLASS, pkcs11.CKO_PRIVATE_KEY),
		pkcs11.NewAttribute(pkcs11.CKA_LABEL, privateKeyLabel),
	}
	private, err := s.findObject(session, template)
	if err != nil {
		if signaturemanager.IsNotFoundError(err) {
			return nil, signaturemanager.NewNotFoundError().WithMessage(fmt.Sprintf("private key not found for address '%s'. Error: %v", address.String(), err))
		}
		return nil, err
	}

	tracer.Debug("signing")
	err = s.pkcsContext.SignInit(session, []*pkcs11.Mechanism{pkcs11.NewMechanism(pkcs11.CKM_ECDSA, nil)}, *private)
	if err != nil {
		return nil, toSignatureManagerErr(err).WithMessage(fmt.Sprintf("error initializing signature: %v", err))
	}
	sig, err := s.pkcsContext.Sign(session, payloadToSign)
	if err != nil {
		return nil, toSignatureManagerErr(err).WithMessage(fmt.Sprintf("error signing data: %v", err))
	}
	return sig, nil
}

// setLabel sets the label for the given object.
func (s *PKCS11HSMSignatureManager) setLabel(session pkcs11.SessionHandle, objectHandle pkcs11.ObjectHandle, label string) error {
	attributeTemplate := []*pkcs11.Attribute{
		pkcs11.NewAttribute(pkcs11.CKA_LABEL, label),
	}
	return s.pkcsContext.SetAttributeValue(session, objectHandle, attributeTemplate)
}

func (s *PKCS11HSMSignatureManager) getLabel(session pkcs11.SessionHandle, objectHandle pkcs11.ObjectHandle) (*string, error) {
	attributeTemplate := []*pkcs11.Attribute{
		pkcs11.NewAttribute(pkcs11.CKA_LABEL, nil),
	}
	as, err := s.pkcsContext.GetAttributeValue(session, objectHandle, attributeTemplate)
	if err != nil {
		return nil, signererrors.InternalFromErr(err).WithMessage(fmt.Sprintf("error retrieving attribute value: %v", err))
	}
	result := string(as[0].Value)
	return &result, nil
}

func (s *PKCS11HSMSignatureManager) getEllipticCurveParameters() []byte {
	// GetECCurveParams returns an elliptic curve parameters for the Ethereum curve
	switch s.connectionDetails.Configuration.Curve {
	case CurveSecp256k1:
		fallthrough
	default:
		return []byte{0x06, 0x05, 0x2B, 0x81, 0x04, 0x00, 0x0A} // The value of this parameter is the DER encoding of an ANSI X9.62 Parameters value for the curve secp256k1. Reference: OID 1.3.132.0.10
	}
}

// getAddress returns an address for the given public key.
func (s *PKCS11HSMSignatureManager) getAddress(session pkcs11.SessionHandle, publicKeyHandle pkcs11.ObjectHandle) (*address.Address, error) {
	ecp, err := s.getDecodedECPoint(session, publicKeyHandle)
	if err != nil {
		return nil, signaturemanager.NewKeyGenerationError().WithMessage(fmt.Sprintf("%v", err))
	}

	pubKey, err := curves.ParsePubKey(ecp)
	if err != nil {
		return nil, signaturemanager.NewInternalError().WithMessage(fmt.Sprintf("unable to parse public key. Error: %v", err))
	}

	derivedAddr, err := signaturemanager.DeriveAddressFromPublicKey(pubKey.SerializeUncompressed())
	if err != nil {
		return nil, err
	}
	return derivedAddr, nil
}

// getDecodedECPoint decodes the EC point and removes the DER encoding.
func (s *PKCS11HSMSignatureManager) getDecodedECPoint(session pkcs11.SessionHandle, publicKeyHandle pkcs11.ObjectHandle) ([]byte, error) {
	ecPoint, err := s.getECPoint(session, publicKeyHandle)
	if err != nil {
		return nil, err
	}
	var ecp []byte
	_, err = asn1.Unmarshal(ecPoint, &ecp)
	if err != nil {
		return nil, signaturemanager.NewKeyGenerationError().WithMessage(fmt.Sprintf("from error: %v", err))
	}
	return ecp, nil
}

// getECPoint returns the CKA_EC_POINT of the given public key.
func (s *PKCS11HSMSignatureManager) getECPoint(session pkcs11.SessionHandle, publicKeyHandle pkcs11.ObjectHandle) ([]byte, error) {
	attributeTemplate := []*pkcs11.Attribute{
		pkcs11.NewAttribute(pkcs11.CKA_EC_POINT, nil),
	}
	attributes, err := s.pkcsContext.GetAttributeValue(session, publicKeyHandle, attributeTemplate)
	if err != nil {
		return nil, err
	}

	var ecPoint []byte
	for _, attribute := range attributes {
		if attribute.Type == pkcs11.CKA_EC_POINT {
			ecPoint = attribute.Value
			break
		}
	}
	if ecPoint != nil {
		return ecPoint, nil
	}
	return nil, signaturemanager.NewKeyGenerationError().WithMessage("unable to get EC Point")
}

// findObjects returns a list of objects for the given attributes template.
func (s *PKCS11HSMSignatureManager) findObjects(session pkcs11.SessionHandle, template []*pkcs11.Attribute) ([]pkcs11.ObjectHandle, error) {
	err := s.pkcsContext.FindObjectsInit(session, template)
	if err != nil {
		return nil, signaturemanager.NewInternalError().WithMessage(fmt.Sprintf("PKCS11 find objects init failed. Error: %v", err))
	}
	var done = false
	var objects = make([]pkcs11.ObjectHandle, 0)
	for !done {
		os, _, errFindObjects := s.pkcsContext.FindObjects(session, 10000) //the second returned parameter is deprecated and must be ignored
		if errFindObjects != nil {
			return nil, signaturemanager.NewInternalError().WithMessage(fmt.Sprintf("PKCS11 find objects failed. Error: %v", errFindObjects))
		}
		objects = append(objects, os...)
		done = len(os) == 0 // when the library has found all the objects, it returns an empty array
	}
	err = s.pkcsContext.FindObjectsFinal(session)
	if err != nil {
		return nil, signaturemanager.NewInternalError().WithMessage(fmt.Sprintf("PKCS11 find objects final failed. Error: %v", err))
	}
	return objects, nil
}

// FindObject returns an object for the given attribute template.
func (s *PKCS11HSMSignatureManager) findObject(session pkcs11.SessionHandle, template []*pkcs11.Attribute) (*pkcs11.ObjectHandle, error) {
	err := s.pkcsContext.FindObjectsInit(session, template)
	if err != nil {
		return nil, err
	}
	os, _, err := s.pkcsContext.FindObjects(session, 1)
	if err != nil {
		return nil, err
	}
	err = s.pkcsContext.FindObjectsFinal(session)
	if err != nil {
		return nil, err
	}
	if len(os) == 0 {
		return nil, signaturemanager.NewNotFoundError().WithMessage("object not found")
	}
	return &os[0], nil
}

// getTimestamp gets the timestamp from the CKA_ID attribute.
func (s *PKCS11HSMSignatureManager) getTimestamp(session pkcs11.SessionHandle, objectHandle pkcs11.ObjectHandle) (uint64, error) {
	attributeTemplate := []*pkcs11.Attribute{
		pkcs11.NewAttribute(pkcs11.CKA_ID, nil),
	}
	as, err := s.pkcsContext.GetAttributeValue(session, objectHandle, attributeTemplate)
	if err != nil {
		return 0, err
	}
	if len(as[0].Value) == 8 {
		return binary.BigEndian.Uint64(as[0].Value), nil
	}
	return 0, err
}

// logOutAndClose logout and close the pkcs11 session
func (s *PKCS11HSMSignatureManager) logOut(tracer logger.Tracer, session pkcs11.SessionHandle) {
	tracer.Debug("logging out")
	err := s.pkcsContext.Logout(session)
	if err != nil {
		tracer.Errorf("logout failed. Error: %v", err)
	}
}

// logOutAndClose logout and close the pkcs11 session
func (s *PKCS11HSMSignatureManager) closeSession(tracer logger.Tracer, session pkcs11.SessionHandle) {
	tracer.Debug("closing session")
	err := s.pkcsContext.CloseSession(session)
	if err != nil {
		tracer.Errorf("closing session failed. Error: %v", err)
	}
}

// GenerateTimestampId returns a timestamp to use as CKA_ID so keys can be sorted chronologically
func generateTimestampId() []byte {
	t := time.Now().UnixNano()
	ts := make([]byte, 8)
	binary.BigEndian.PutUint64(ts, uint64(t))
	return ts
}

// calculatePublicKeyLabel calculates the label used for the underlying PKCS11 storage when storing the public key.
func calculatePublicKeyLabel(addr address.Address) string {
	return addr.String()
}

// calculatePrivateKeyLabel calculates the label used for the underlying PKCS11 storage when storing the private key. It is the EIP55-compliant hex encoded string representation of the address.
func calculatePrivateKeyLabel(addr address.Address) string {
	return addr.String()
}

// addressTime defines a structure to be used when sorting accounts by time.
type addressTime struct {
	Address address.Address
	Time    uint64
}

// ByTime defines a type to be used when sorting accounts by time.
type ByTime []addressTime

// Len is implementing the required element interface to be used by sort.Sort.
func (a ByTime) Len() int { return len(a) }

// Less is implementing the required element interface to be used by sort.Sort.
func (a ByTime) Less(i, j int) bool { return a[i].Time < a[j].Time }

// Swap is implementing the required element interface to be used by sort.Sort.
func (a ByTime) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

// Addresses is implementing the required element interface to be used by sort.Sort.
func (a ByTime) Addresses() []address.Address {
	addresses := make([]address.Address, 0)
	for _, at := range a {
		addresses = append(addresses, at.Address)
	}
	return addresses
}

// toSignatureManagerErr translates pkcs11 errors into signature manager errors.
func toSignatureManagerErr(originalErr error, message ...string) *signaturemanager.Error {
	errMsg := originalErr.Error()
	if len(message) > 0 {
		errMsg = fmt.Sprintf("%s: %s", message, errMsg)
	}

	var pkcs11Err pkcs11.Error
	if errors.As(originalErr, &pkcs11Err) {
		err, ok := pkcsErrTranslator[pkcs11Err]
		if ok {
			return err.WithMessage(errMsg)
		}
	}
	return signaturemanager.NewInternalError().WithMessage(errMsg)
}
