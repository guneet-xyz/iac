package env

import (
	"iac/config"
	"iac/utils/fs"
	"log/slog"
	"path/filepath"
	"strings"
)

func ListEnvEntries() ([]string, error) {
	var conf = config.GetConfig()
	var dirEntries, err = fs.ReadDir(conf.Environment.DirectoryPath)
	if err != nil {
		slog.Error("Failed to list env entries", "error", err)
		return nil, err
	}

	var entries []string

	for _, entry := range dirEntries {
		var entryName = entry.Name()
		if entry.Type().IsRegular() && strings.HasSuffix(entryName, ".json") {
			entryNameWithoutExt := strings.TrimSuffix(entryName, ".json")
			entries = append(entries, entryNameWithoutExt)
		}
	}

	return entries, nil
}

func GetEnv(name string) (Env, error) {
	var envType, err = readStoredEnvType(name)
	if err != nil {
		slog.Error("Failed to read env type", "error", err, "name", name)
		return Env{}, err
	}

	switch envType {
	case EnvTypeSecret:
		storedEnv, err := readStoredEnvSecret(name)
		if err != nil {
			slog.Error("Failed to read stored env secret", "error", err, "name", name)
			return Env{}, err
		}

		plaintext, err := decryptSecret(storedEnv)
		if err != nil {
			slog.Error("Failed to decrypt stored env secret", "error", err, "name", name)
			return Env{}, err
		}

		return Env{
			Type:  EnvTypeSecret,
			Name:  name,
			Value: plaintext,
		}, nil

	case EnvTypeVariable:
		storedEnv, err := readStoredEnvVariable(name)
		if err != nil {
			slog.Error("Failed to read stored env variable", "error", err, "name", name)
			return Env{}, err
		}

		return Env{
			Type:  EnvTypeVariable,
			Name:  name,
			Value: storedEnv.Value,
		}, nil

	default:
		slog.Error("Unknown env type", "type", envType, "name", name)
		return Env{}, nil
	}
}

func SetEnv(env Env) error {
	var err error
	var storedEnv any

	switch env.Type {
	case EnvTypeSecret:
		storedEnv, err = encryptSecret(env.Value)
		if err != nil {
			slog.Error("Failed to encrypt secret", "error", err, "name", env.Name)
			return err
		}
	case EnvTypeVariable:
		storedEnv = StoredEnvVariable{
			Value: env.Value,
		}
	default:
		slog.Error("Unknown env type", "type", env.Type, "name", env.Name)
		return nil
	}

	err = writeStoredEnv(env.Name, env.Type, storedEnv)
	if err != nil {
		slog.Error("Failed to write stored env variable", "error", err, "name", env.Name)
		return err
	}

	return nil

}

func DeleteEnv(name string) error {
	var conf = config.GetConfig()
	var filePath = filepath.Join(conf.Environment.DirectoryPath, name+".json")
	err := fs.DeleteFile(filePath)
	if err != nil {
		slog.Error("Failed to delete env entry", "error", err, "name", name)
		return err
	}
	return nil
}
