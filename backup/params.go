package graal

import (
	"strconv"
	"sync"
)

type Params struct {
	data map[string]interface{}
	l    sync.RWMutex
}

func (params *Params) Get(key string) (interface{}, bool) {
	params.l.RLock()
	defer params.l.RUnlock()
	if v := params.data[key]; v != nil {
		return v, true
	}
	return nil, false
}

func (params *Params) Has(key string) bool {
	params.l.RLock()
	defer params.l.RUnlock()
	if v := params.data[key]; v != nil {
		return true
	}
	return false
}

func (params *Params) GetInt(key string) (int, bool) {
	if v, ok := params.Get(key); ok {
		return params.toInt(v)
	}
	return 0, false
}

func (params *Params) GetUint(key string) (uint, bool) {
	if v, ok := params.Get(key); ok {
		return params.toUint(v)
	}
	return 0, false
}

func (params *Params) GetFloat(key string) (float64, bool) {
	if v, ok := params.Get(key); ok {
		return params.toFloat(v)
	}
	return 0, false
}

func (params *Params) GetString(key string) (string, bool) {
	if v, ok := params.Get(key); ok {
		return params.toString(v)
	}
	return "", false
}

func (params *Params) GetBool(key string) (bool, bool) {
	if v, ok := params.Get(key); ok {
		return params.toBool(v)
	}
	return false, false
}

func (params *Params) Keys() []string {
	params.l.RLock()
	defer params.l.RUnlock()
	keys := make([]string, 0)
	if params.data != nil {
		for key := range params.data {
			keys = append(keys, key)
		}
	}
	return keys
}

func (*Params) toString(v interface{}) (string, bool) {
	if v, ok := v.(string); ok {
		return v, true
	}
	return "", false
}

func (*Params) toInt(v interface{}) (int, bool) {
	switch v := v.(type) {
	case int:
		return int(v), true
	case uint:
		return int(v), true
	case float32:
		return int(v), true
	case float64:
		return int(v), true
	case string:
		if v, err := strconv.ParseInt(v, 10, 64); err == nil {
			return int(v), true
		}
	}
	return 0, false
}

func (*Params) toUint(v interface{}) (uint, bool) {
	switch v := v.(type) {
	case int:
		return uint(v), true
	case uint:
		return uint(v), true
	case float32:
		return uint(v), true
	case float64:
		return uint(v), true
	case string:
		if v, err := strconv.ParseUint(v, 10, 64); err == nil {
			return uint(v), true
		}
	}
	return 0, false
}

func (*Params) toFloat(v interface{}) (float64, bool) {
	switch v := v.(type) {
	case int:
		return float64(v), true
	case uint:
		return float64(v), true
	case float32:
		return float64(v), true
	case float64:
		return v, true
	case string:
		if v, err := strconv.ParseFloat(v, 64); err == nil {
			return v, true
		}
	}
	return 0, false
}

func (*Params) toBool(v interface{}) (bool, bool) {
	if v, ok := v.(bool); ok {
		return v, true
	}
	if v, ok := v.(string); ok {
		if v, err := strconv.ParseBool(v); err == nil {
			return v, true
		}
	}
	return false, false
}

type ParamsWriter struct {
	Params
}

func (writer *ParamsWriter) Set(key string, value interface{}) {
	writer.Params.l.Lock()
	defer writer.Params.l.Unlock()
	if writer.Params.data == nil {
		writer.Params.data = make(map[string]interface{})
	}
	writer.Params.data[key] = value
}

func (writer *ParamsWriter) Merge(params *Params) {
	writer.Params.l.Lock()
	defer writer.Params.l.Unlock()
	if writer.Params.data == nil {
		writer.Params.data = make(map[string]interface{})
	}
	for _, key := range params.Keys() {
		writer.Params.data[key], _ = params.Get(key)
	}
}
