package secret

import (
	"iac/config"
	"iac/utils/fs"
	"iac/utils/json"
	"path/filepath"
)

type Kcv struct {
	SaltBase64  string
	IvBase64    string
	ValueBase64 string
}

func doesKcvExist() bool {
	config := config.GetConfig()
	kcvPath := filepath.Join(config.SecretsDirPath, "kcv.json")
	stat, err := fs.Stat(kcvPath)
	if err != nil {
		return false
	}
	return stat == fs.StatResultFile
}

func getKcv() (Kcv, error) {
	config := config.GetConfig()
	kcvPath := filepath.Join(config.SecretsDirPath, "kcv.json")
	bytes, err := fs.ReadFileAsBytes(kcvPath)
	if err != nil {
		return Kcv{}, err
	}
	var kcv Kcv
	err = json.Unmarshal(bytes, &kcv)
	if err != nil {
		return Kcv{}, err
	}
	return kcv, nil
}

func saveKcv(kcv Kcv) error {
	config := config.GetConfig()
	kcvPath := filepath.Join(config.SecretsDirPath, "kcv.json")
	bytes, err := json.Marshal(kcv)
	if err != nil {
		return err
	}
	return fs.WriteFileFromBytes(kcvPath, bytes)
}

var KcvPlainText = "Be My Baby - The Ronettes"

func generateKcv(saltBase64 string, ivBase64 string, keyBase64 string) (Kcv, error) {
	encryptedKcvValue, err := encryptAES256GCMB64(KcvPlainText, keyBase64, ivBase64)
	if err != nil {
		return Kcv{}, err
	}

	kcv := Kcv{
		SaltBase64:  saltBase64,
		IvBase64:    ivBase64,
		ValueBase64: encryptedKcvValue,
	}
	return kcv, nil
}

func verifyKeyWithKcv(keyBase64 string) (bool, error) {
	kcv, err := getKcv()
	if err != nil {
		return false, err
	}
	decryptedKcvValue, err := decryptAES256GCMB64(kcv.ValueBase64, keyBase64, kcv.IvBase64)
	if err != nil {
		return false, err
	}
	return decryptedKcvValue == KcvPlainText, nil
}
