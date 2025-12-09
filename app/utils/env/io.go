package env

import (
	"iac/config"
	"iac/utils/fs"
	"iac/utils/json"
	"log/slog"
	"path/filepath"
)

type EnvType string

var EnvTypeSecret EnvType = "secret"
var EnvTypeVariable EnvType = "variable"

type Env struct {
	Type  EnvType
	Name  string
	Value string
}

type StoredEnv struct {
	Type EnvType `json:"type"`
}

type StoredEnvSecret struct {
	IvBase64    string `json:"iv_b64"`
	ValueBase64 string `json:"value_b64"`
}

type StoredEnvVariable struct {
	Value string `json:"value"`
}

func storedEnvFilePath(name string) (string, error) {
	cfg, err := config.GetConfig()
	if err != nil {
		slog.Error("Failed to get config", "error", err)
		return "", err
	}
	return filepath.Join(cfg.Environment.DirectoryPath, name+".json"), nil
}

func readStoredEnvType(name string) (EnvType, error) {
	envFilePath, err := storedEnvFilePath(name)
	if err != nil {
		slog.Error("Failed to get stored env file path", "error", err, "name", name)
		return "", err
	}
	bytes, err := fs.ReadFileAsBytes(envFilePath)
	if err != nil {
		slog.Error("Failed to read stored env file", "error", err, "name", name)
		return "", err
	}

	var storedEnv StoredEnv
	err = json.Unmarshal(bytes, &storedEnv)
	if err != nil {
		slog.Error("Failed to unmarshal stored env JSON", "error", err, "name", name)
		return "", err
	}

	if storedEnv.Type != EnvTypeSecret && storedEnv.Type != EnvTypeVariable {
		slog.Error("Unknown stored env type", "type", storedEnv.Type, "name", name)
		return "", nil
	}

	return storedEnv.Type, nil
}

func readStoredEnvSecret(name string) (StoredEnvSecret, error) {
	envFilePath, err := storedEnvFilePath(name)
	if err != nil {
		slog.Error("Failed to get stored env file path", "error", err, "name", name)
		return StoredEnvSecret{}, err
	}
	bytes, err := fs.ReadFileAsBytes(envFilePath)
	if err != nil {
		slog.Error("Failed to read stored env secret file", "error", err, "name", name)
		return StoredEnvSecret{}, err
	}

	var storedEnvSecret StoredEnvSecret
	err = json.Unmarshal(bytes, &storedEnvSecret)
	if err != nil {
		slog.Error("Failed to unmarshal stored env secret JSON", "error", err, "name", name)
		return StoredEnvSecret{}, err
	}

	return storedEnvSecret, nil
}

func readStoredEnvVariable(name string) (StoredEnvVariable, error) {
	envFilePath, err := storedEnvFilePath(name)
	if err != nil {
		slog.Error("Failed to get stored env file path", "error", err, "name", name)
		return StoredEnvVariable{}, err
	}

	bytes, err := fs.ReadFileAsBytes(envFilePath)
	if err != nil {
		slog.Error("Failed to read stored env variable file", "error", err, "name", name)
		return StoredEnvVariable{}, err
	}

	var storedEnvVariable StoredEnvVariable
	err = json.Unmarshal(bytes, &storedEnvVariable)
	if err != nil {
		slog.Error("Failed to unmarshal stored env variable JSON", "error", err, "name", name)
		return StoredEnvVariable{}, err
	}

	return storedEnvVariable, nil
}

func writeStoredEnv(name string, envType EnvType, storedEnv any) error {
	combined, err := json.Combine(StoredEnv{Type: envType}, storedEnv)
	if err != nil {
		slog.Error("Failed to combine stored env data", "error", err, "name", name)
		return err
	}

	bytes, err := json.Marshal(combined)
	if err != nil {
		slog.Error("Failed to marshal stored env to JSON", "error", err, "name", name)
		return err
	}

	envFilePath, err := storedEnvFilePath(name)
	if err != nil {
		slog.Error("Failed to get stored env file path", "error", err, "name", name)
		return err
	}

	err = fs.WriteFileFromBytes(envFilePath, bytes)
	if err != nil {
		slog.Error("Failed to write stored env file", "error", err, "name", name)
		return err
	}

	return nil
}
