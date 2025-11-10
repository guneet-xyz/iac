package json

import (
	"encoding/json"
	"iac/utils/errors"
	"log/slog"
)

func Unmarshal(data []byte, v any) error {
	slog.Debug("Unmarshal called", "data length", len(data))
	err := json.Unmarshal(data, v)
	if err != nil {
		return errors.New("Failed to unmarshal JSON data", "error", err)
	}
	slog.Debug("Unmarshal successful", "unmarshaled", v)
	return nil
}
