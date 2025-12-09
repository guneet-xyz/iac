package secret

import (
	"iac/config"
	"iac/utils/errors"
	"iac/utils/fs"
	"log/slog"
)

func DoesMasterKeyExist() (bool, error) {
	slog.Debug("Checking if master key exists")
	conf, err := config.GetConfig()
	if err != nil {
		slog.Error("Error while trying to get config", "error", err)
		return false, err
	}

	stat, err := fs.Stat(conf.MasterKey.MasterKeyPath)
	if err != nil {
		slog.Error("Error while trying to stat master key file", "error", err)
		return false, err
	}
	if stat != fs.StatResultFile {
		slog.Debug("Master key file does not exist")
		return false, nil
	}
	slog.Debug("Master key file exists")
	masterKeyB64, err := GetMasterKeyB64()
	if err != nil {
		slog.Error("Error while trying to get master key", "error", err)
		return false, err
	}
	if masterKeyB64 == "" {
		slog.Debug("Master key is empty")
		return false, nil
	}
	slog.Debug("Master key exists and is non-empty")
	return true, nil
}

func GetMasterKeyB64() (string, error) {
	slog.Debug("Retrieving master key")
	conf, err := config.GetConfig()
	if err != nil {
		slog.Error("Error while trying to get config", "error", err)
		return "", err
	}
	keyB64, err := fs.ReadFileAsString(conf.MasterKey.MasterKeyPath)
	slog.Debug("Master key read from file", "length", len(keyB64))
	if err != nil {
		slog.Error("Error while trying to read master key from file", "error", err)
		return "", err
	}
	good, err := verifyKeyWithKcv(keyB64)
	if err != nil {
		slog.Error("Error while trying to verify master key from file against KCV", "error", err)
		return "", err
	}
	if !good {
		return "", errors.New("Master key in keyring does not match KCV")
	}
	return keyB64, nil
}

func SetMasterKey(plainTextValue string) error {
	slog.Debug("Setting master key")
	var err error
	kcvExists, err := doesKcvExist()
	if err != nil {
		slog.Error("Error while trying to check if KCV exists", "error", err)
		return err
	}
	var keyB64 string
	if kcvExists {
		kcv, err := getKcv()
		if err != nil {
			return err
		}
		saltB64 := kcv.SaltBase64
		keyB64, err = GenerateKeyB64(plainTextValue, saltB64)
		if err != nil {
			return err
		}
		valid, err := verifyKeyWithKcv(keyB64)
		if err != nil {
			slog.Error("Error while trying to verify provided master key against existing KCV", "error", err)
			return err
		}
		if !valid {
			return errors.New("provided master key does not match existing KCV")
		}
		slog.Info("Master key verified.")
	} else {
		saltB64 := GenerateSaltB64()
		ivB64 := GenerateIVB64()
		keyB64, err = GenerateKeyB64(plainTextValue, saltB64)
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
		slog.Info("KCV created.")
	}

	conf, err := config.GetConfig()
	if err != nil {
		slog.Error("Error while trying to get config", "error", err)
		return err
	}

	err = fs.WriteFileFromString(conf.MasterKey.MasterKeyPath, keyB64)
	if err != nil {
		slog.Error("Error while trying to write master key to file", "error", err)
		return err
	}

	slog.Debug("Master key set successfully")
	return nil
}
