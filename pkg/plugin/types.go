package plugin

import (
	"encoding/json"
	"fmt"
	"time"
)

type Level uint8

const (
	LevelUnknown Level = iota
	LevelTrace
	LevelDebug
	LevelInfo
	LevelWarn
	LevelError
	LevelFatal
)

func (l Level) String() string {
	switch l {
	case LevelTrace:
		return "TRACE"
	case LevelDebug:
		return "DEBUG"
	case LevelInfo:
		return "INFO"
	case LevelWarn:
		return "WARN"
	case LevelError:
		return "ERROR"
	case LevelFatal:
		return "FATAL"
	default:
		return ""
	}

}

func (l Level) MarshalText() (text []byte, err error) {
	s := l.String()
	if s == "" {
		return nil, fmt.Errorf("incorrect level %d", l)
	}
	return []byte(s), nil
}

func (l *Level) UnmarshalText(text []byte) error {
	switch string(text) {
	case "TRACE":
		*l = LevelTrace
	case "DEBUG":
		*l = LevelDebug
	case "INFO":
		*l = LevelInfo
	case "WARN":
		*l = LevelWarn
	case "ERROR":
		*l = LevelError
	case "FATAL":
		*l = LevelFatal
	default:
		*l = LevelUnknown
	}
	return nil
}

type logEntriesValues struct {
	extractFields []string

	timestamps    []time.Time
	contents      []string
	levels        []string
	ids           []string
	streams       []string
	resourceTypes []string
	resourceIDs   []string
	messages      []string
	payloads      []string
}

func (v *logEntriesValues) append(
	timestamp time.Time,
	level Level,
	id string,
	stream string,
	resourceType string,
	resourceID string,
	message string,
	payload map[string]any,
) {
	v.timestamps = append(v.timestamps, timestamp.UTC())
	v.levels = append(v.levels, level.String())
	v.ids = append(v.ids, id)
	v.streams = append(v.streams, stream)
	v.resourceTypes = append(v.resourceTypes, resourceType)
	v.resourceIDs = append(v.resourceIDs, resourceID)
	v.messages = append(v.messages, message)
	v.payloads = append(v.payloads, v.jsonValue(payload))

	content := "{}"
	if len(v.extractFields) > 0 {
		extractedFields := make(map[string]any, len(v.extractFields))
		for _, k := range v.extractFields {
			if v := payload[k]; v != nil {
				extractedFields[k] = v
			}
		}
		content = v.jsonValue(extractedFields)
	}
	v.contents = append(v.contents, message+" "+content)
}

func (logEntriesValues) jsonValue(val any) string {
	b, err := json.Marshal(val)
	if err != nil {
		return "{}"
	}
	return string(b)
}
