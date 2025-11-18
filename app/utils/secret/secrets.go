package secret

import (
	"errors"
	"iac/config"
	"iac/utils/fs"
	"iac/utils/json"
	"log/slog"
	"path/filepath"
)

type Secret struct {
	IvBase64    string
	ValueBase64 string
}

func GetSecret(secretName string) (string, error) {
	slog.Debug("Getting secret", "secretName", secretName)
	keyB64, err := GetMasterKeyB64()
	if err != nil {
		slog.Error("Failed to get master key", "error", err)
		return "", err
	}
	conf := config.GetConfig()
	secretFilePath := filepath.Join(conf.SecretsDirPath, secretName+".json")
	stat, err := fs.Stat(secretFilePath)
	if err != nil {
		slog.Error("Failed to stat secret file", "error", err, "secretName", secretName)
		return "", err
	}
	if stat != fs.StatResultFile {
		slog.Error("Secret does not exist", "secretName", secretName)
		return "", errors.New("secret does not exist")
	}
	secret := Secret{}
	bytes, err := fs.ReadFileAsBytes(secretFilePath)
	if err != nil {
		slog.Error("Failed to read secret file", "error", err, "secretName", secretName)
		return "", err
	}
	err = json.Unmarshal(bytes, &secret)
	if err != nil {
		slog.Error("Failed to unmarshal secret JSON", "error", err, "secretName", secretName)
		return "", err
	}
	value, err := decryptAES256GCMB64(secret.ValueBase64, keyB64, secret.IvBase64)
	if err != nil {
		slog.Error("Failed to decrypt secret value", "error", err, "secretName", secretName)
		return "", err
	}
	slog.Debug("Got secret", "secretName", secretName)
	return value, nil
}

func SetSecret(secretName string, secretValue string) error {
	slog.Debug("Setting secret", "secretName", secretName)
	keyB64, err := GetMasterKeyB64()
	if err != nil {
		slog.Error("Failed to get master key", "error", err)
		return err
	}
	ivB64 := GenerateIVB64()
	encryptedValueB64, err := encryptAES256GCMB64(secretValue, keyB64, ivB64)
	if err != nil {
		slog.Error("Failed to encrypt secret value", "error", err)
		return err
	}
	secret := Secret{
		IvBase64:    ivB64,
		ValueBase64: encryptedValueB64,
	}
	bytes, err := json.Marshal(secret)
	if err != nil {
		slog.Error("Failed to marshal secret to JSON", "error", err)
		return err
	}
	conf := config.GetConfig()
	secretFilePath := filepath.Join(conf.SecretsDirPath, secretName+".json")
	err = fs.WriteFileFromBytes(secretFilePath, bytes)
	if err != nil {
		slog.Error("Failed to write secret file", "error", err, "secretFilePath", secretFilePath)
		return err
	}
	slog.Debug("Secret has been set", "secretName", secretName)
	return nil
}

func GetSecrets() ([]string, error) {
	slog.Debug("Listing secrets")
	conf := config.GetConfig()
	secrets, err := fs.ReadDir(conf.SecretsDirPath)
	if err != nil {
		slog.Error("Failed to read secrets directory", "error", err)
		return nil, err
	}
	var secretNames []string
	for _, entry := range secrets {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if filepath.Ext(name) == ".json" {
			secretName := name[:len(name)-len(".json")]
			secretNames = append(secretNames, secretName)
		}
	}
	slog.Debug("Listed secrets", "count", len(secretNames))
	return secretNames, nil
}

func DeleteSecret(secretName string) error {
	slog.Debug("Deleting secret", "secretName", secretName)
	conf := config.GetConfig()
	secretFilePath := filepath.Join(conf.SecretsDirPath, secretName+".json")
	err := fs.DeleteFile(secretFilePath)
	if err != nil {
		slog.Error("Failed to delete secret file", "error", err, "secretName", secretName)
		return err
	}
	slog.Debug("Deleted secret", "secretName", secretName)
	return nil
}
