package plugin

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/cbroglie/mustache"
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
	derivedRules  map[string]*mustache.Template

	timestamps    []time.Time
	contents      []string
	levels        []string
	ids           []string
	streams       []string
	resourceTypes []string
	resourceIDs   []string
	messages      []string
	payloads      []string
	derived       map[string][]*string
}

func (v *logEntriesValues) addDerivedRules(rules []derivedFieldRule) {
	if v.derivedRules == nil {
		v.derivedRules = make(map[string]*mustache.Template)
	}
	for _, rule := range rules {
		if tmpl, err := mustache.ParseString(rule.Template); err == nil {
			v.derivedRules[rule.Name] = tmpl
		}
	}
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

	var derived map[string]string
	if len(v.derivedRules) > 0 {
		derived = make(map[string]string, len(v.derivedRules))
		if len(v.derivedRules) > 0 && v.derived == nil {
			v.derived = make(map[string][]*string)
		}

		fields := renderCtx(timestamp, level, id, stream, resourceType, resourceID, message)
		for name, tmpl := range v.derivedRules {
			value, err := tmpl.Render(payload, fields)
			derived[name] = value
			if err != nil || value == "" {
				v.derived[name] = append(v.derived[name], nil)
				continue
			}
			v.derived[name] = append(v.derived[name], &value)
		}
	}

	content := "{}"
	if len(v.extractFields) > 0 {
		extractedFields := make(map[string]any, len(v.extractFields))
		for _, k := range v.extractFields {
			if v := payload[k]; v != nil {
				extractedFields[k] = v
				continue
			}
			if v := derived[k]; v != "" {
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

func renderCtx(
	timestamp time.Time,
	level Level,
	id string,
	stream string,
	resourceType string,
	resourceID string,
	message string,
) map[string]any {
	return map[string]any{
		"timestamp":     timestamp,
		"level":         level,
		"id":            id,
		"stream":        stream,
		"resource_type": resourceType,
		"resource_id":   resourceID,
		"message":       message,
	}
}
