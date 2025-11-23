package json

import (
	"log/slog"
	"maps"
)

func Combine(objects ...any) (map[string]any, error) {
	var bytesArray [][]byte
	var combined map[string]any = make(map[string]any)

	for _, obj := range objects {
		bytes, err := Marshal(obj)
		if err != nil {
			slog.Error("Failed to marshal object to JSON", "error", err)
			return nil, err
		}
		bytesArray = append(bytesArray, bytes)
	}

	for _, bytes := range bytesArray {
		var partial map[string]any = make(map[string]any)
		err := Unmarshal(bytes, &partial)
		if err != nil {
			slog.Error("Failed to unmarshal JSON bytes", "error", err)
			return nil, err
		}
		maps.Copy(combined, partial)
	}

	return combined, nil
}
