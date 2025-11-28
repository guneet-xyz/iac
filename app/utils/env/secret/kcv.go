package secret

import (
	"iac/config"
	"iac/utils/fs"
	"iac/utils/json"
)

type Kcv struct {
	SaltBase64  string
	IvBase64    string
	ValueBase64 string
}

func doesKcvExist() bool {
	conf := config.GetConfig()
	stat, err := fs.Stat(conf.MasterKey.KcvPath)
	if err != nil {
		return false
	}
	return stat == fs.StatResultFile
}

func getKcv() (Kcv, error) {
	conf := config.GetConfig()
	bytes, err := fs.ReadFileAsBytes(conf.MasterKey.KcvPath)
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
	conf := config.GetConfig()
	bytes, err := json.Marshal(kcv)
	if err != nil {
		return err
	}
	return fs.WriteFileFromBytes(conf.MasterKey.KcvPath, bytes)
}

var KcvPlainText = "Be My Baby - The Ronettes"

func generateKcv(saltBase64 string, ivBase64 string, keyBase64 string) (Kcv, error) {
	encryptedKcvValue, err := EncryptAES256GCMB64(KcvPlainText, keyBase64, ivBase64)
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
	decryptedKcvValue, err := DecryptAES256GCMB64(kcv.ValueBase64, keyBase64, kcv.IvBase64)
	if err != nil {
		return false, err
	}
	return decryptedKcvValue == KcvPlainText, nil
}

func VerifyPassphraseWithKcv(passphrase string) (bool, error) {
	kcv, err := getKcv()
	if err != nil {
		return false, err
	}
	keyBase64, err := GenerateKeyB64(passphrase, kcv.SaltBase64)
	if err != nil {
		return false, err
	}
	return verifyKeyWithKcv(keyBase64)
}
