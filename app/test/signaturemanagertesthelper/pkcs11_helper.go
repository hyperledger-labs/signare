package signaturemanagertesthelper

import (
	"errors"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/hyperledger-labs/signare/app/pkg/commons/time"
)

const (
	SoftHSMLib         = "/usr/local/lib/softhsm/libsofthsm2.so"
	SlotPin            = "userpin"
	ImportedKeyAddress = "0xa2c16184fA76cD6D16685900292683dF905e4Bf2"

	tokenLabel                       = "WALLET-000"
	relativePathToPrivateKeyFileName = "app/test/signaturemanagertesthelper/testdata/private.pem"
	importedKeyID                    = "0001"
	privateKey                       = `-----BEGIN PRIVATE KEY-----
MIGEAgEAMBAGByqGSM49AgEGBSuBBAAKBG0wawIBAQQgEACcTiQhuQ9kiLjdLx6N
G2EcQ7i3aST+P0si+B3QomahRANCAAQbSGsG/Mmtv3Tqh9h5RClPY3UGAom2+kAv
/E9AReYIpjXFumvF4SxujDHZA/y1VRSPcb+b/LVY/s2Oo9nyk3MC
-----END PRIVATE KEY-----`
)

// InitializeSoftHSMSlot initializes a slot in SoftHSM and returns its ID. It assumes that slot 0 is not yet initialized.
func InitializeSoftHSMSlot() (*string, *string, error) {
	// 1. Make sure that the token from a previous execution is removed so that this is as idempotent as possible
	cmdListSlots := exec.Command("softhsm2-util",
		"--show-slots",
	)
	listOutput, err := cmdListSlots.CombinedOutput()
	if err != nil {
		return nil, nil, err
	}

	serialNumbers, err := findSerialNumbers(string(listOutput))
	if err != nil {
		return nil, nil, err
	}

	for _, serialNumber := range serialNumbers {
		cmdDeleteToken := exec.Command("softhsm2-util",
			"--delete-token",
			"--serial", serialNumber,
		)
		_, deleteErr := cmdDeleteToken.CombinedOutput()
		if deleteErr != nil {
			return nil, nil, deleteErr
		}
	}

	// 2. Initialize a new token in the slot 0
	initFirstTokenOutput, err := initToken("0")
	if err != nil {
		return nil, nil, err
	}
	slotOne, err := findSlotToUse(string(initFirstTokenOutput))
	if err != nil {
		return nil, nil, err
	}

	// 3. Initialize a new token in the slot 1
	initSecondTokenOutput, err := initToken("1")
	if err != nil {
		return nil, nil, err
	}
	slotTwo, err := findSlotToUse(string(initSecondTokenOutput))
	if err != nil {
		return nil, nil, err
	}

	// 3. Import private key in slot 0. It is used to test the sign operation
	privateKeyPath, err := createPrivateKey()
	if err != nil {
		return nil, nil, err
	}
	cmdImportKey := exec.Command("softhsm2-util", //nolint:gosec
		"--import", *privateKeyPath,
		"--slot", *slotOne,
		"--label", ImportedKeyAddress,
		"--id", importedKeyID,
		"--pin", SlotPin,
	)
	importOutput, err := cmdImportKey.CombinedOutput()
	if err != nil {
		return nil, nil, err
	}
	re := regexp.MustCompile(`The key pair has been imported.\n$`)
	match := re.FindStringSubmatch(string(importOutput))
	if len(match) != 1 {
		return nil, nil, errors.New("private key was not imported")
	}

	return slotOne, slotTwo, nil
}

func findSlotToUse(input string) (*string, error) {
	re := regexp.MustCompile(`The token has been initialized and is reassigned to slot (\d+)\n$`)
	match := re.FindStringSubmatch(input)
	if len(match) == 2 {
		return &match[1], nil
	}
	return nil, errors.New("no match found")
}

func findSerialNumbers(input string) ([]string, error) {
	re := regexp.MustCompile(`Serial number:\s*(\w+)`)
	match := re.FindAllStringSubmatch(input, -1)
	if len(match) == 0 {
		return nil, errors.New("no match found")
	}
	var serialNumbers []string
	for _, result := range match {
		if len(result) != 2 {
			continue
		}
		serialNumber := result[1]
		if serialNumber == "Initialized" {
			continue
		}
		serialNumbers = append(serialNumbers, serialNumber)
	}
	return serialNumbers, nil
}

// createPrivateKey creates a private key in a /tmp dir and returns its path
func createPrivateKey() (*string, error) {
	tmpFile, err := os.CreateTemp("", strings.ReplaceAll(time.Now().String(), " ", ""))
	if err != nil {
		return nil, err
	}

	defer func(tmpFile *os.File) {
		closeErr := tmpFile.Close()
		if closeErr != nil {
			panic(closeErr)
		}
	}(tmpFile)

	content := []byte(privateKey)

	if _, writeErr := tmpFile.Write(content); writeErr != nil {
		return nil, writeErr
	}

	tmpFilePath := tmpFile.Name()
	return &tmpFilePath, nil
}

func initToken(slotNumber string) ([]byte, error) {
	cmdInitToken := exec.Command("softhsm2-util",
		"--init-token",
		"--slot", slotNumber,
		"--label", tokenLabel,
		"--pin", SlotPin,
		"--so-pin", "superpin",
	)
	return cmdInitToken.CombinedOutput()
}
