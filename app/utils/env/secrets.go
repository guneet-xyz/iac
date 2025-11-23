package env

import (
	"iac/utils/env/secret"
	"log/slog"
)

func decryptSecret(storedSecret StoredEnvSecret) (string, error) {
	slog.Debug("Decrypting secret", "iv_length", len(storedSecret.IvBase64), "value_length", len(storedSecret.ValueBase64))
	keyB64, err := secret.GetMasterKeyB64()
	if err != nil {
		slog.Error("Failed to get master key", "error", err)
		return "", err
	}
	plain, err := secret.DecryptAES256GCMB64(storedSecret.ValueBase64, keyB64, storedSecret.IvBase64)
	if err != nil {
		slog.Error("Failed to decrypt secret value", "error", err, "iv_length", len(storedSecret.IvBase64), "value_length", len(storedSecret.ValueBase64))
		return "", err
	}
	slog.Debug("Decrypted secret", "plain_length", len(plain))
	return plain, nil
}

func encryptSecret(plain string) (StoredEnvSecret, error) {
	slog.Debug("Encrypting secret", "plain_length", len(plain))
	keyB64, err := secret.GetMasterKeyB64()
	if err != nil {
		slog.Error("Failed to get master key", "error", err)
		return StoredEnvSecret{}, err
	}
	ivB64 := secret.GenerateIVB64()
	encryptedValueB64, err := secret.EncryptAES256GCMB64(plain, keyB64, ivB64)
	if err != nil {
		slog.Error("Failed to encrypt secret value", "error", err)
		return StoredEnvSecret{}, err
	}
	return StoredEnvSecret{
		IvBase64:    ivB64,
		ValueBase64: encryptedValueB64,
	}, nil
}
