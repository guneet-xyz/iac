package json

import (
	"encoding/json"
	"iac/utils/errors"
	"log/slog"
)

func Marshal(v any) ([]byte, error) {
	slog.Debug("Marshal called", "value", v)
	data, err := json.Marshal(v)
	if err != nil {
		return nil, errors.New("Failed to marshal JSON data", "error", err)
	}
	slog.Debug("Marshal successful", "data length", len(data))
	return data, nil
}
