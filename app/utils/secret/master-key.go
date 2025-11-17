package secret

import "iac/utils/errors"

const KeyringMasterKeyName = "master_key"

func GetMasterKey() (string, error) {
	keyB64, err := keyringGet(KeyringMasterKeyName)
	if err != nil {
		return "", err
	}
	good, err := verifyKeyWithKcv(keyB64)
	if err != nil {
		return "", err
	}
	if !good {
		return "", errors.New("master key in keyring does not match KCV")
	}
	return keyB64, nil
}

func SetMasterKey(plainTextValue string) error {
	kcvExists := doesKcvExist()
	if kcvExists {
		kcv, err := getKcv()
		if err != nil {
			return err
		}
		saltB64 := kcv.SaltBase64
		keyB64, err := GenerateKeyB64(plainTextValue, saltB64)
		if err != nil {
			return err
		}
		valid, err := verifyKeyWithKcv(keyB64)
		if err != nil {
			return err
		}
		if !valid {
			return errors.New("provided master key does not match existing KCV")
		}
		return nil
	}

	saltB64 := GenerateSaltB64()
	ivB64 := GenerateIVB64()
	keyB64, err := GenerateKeyB64(plainTextValue, saltB64)
	if err != nil {
		return err
	}

	kcv, err := generateKcv(saltB64, ivB64, keyB64)
	if err != nil {
		return err
	}
	err = saveKcv(kcv)
	if err != nil {
		return err
	}
	return keyringSet(KeyringMasterKeyName, keyB64)
}
