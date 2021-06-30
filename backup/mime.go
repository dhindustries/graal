package graal

import (
	"fmt"
	"strconv"
	"strings"
)

type Mime string

func (mime Mime) SplitParams() Mime {
	fullType := mime.fullType()
	var t, s string
	t = strings.TrimSpace(fullType[0])
	if len(fullType) == 1 {
		s = "*"
	} else {
		s = strings.TrimSpace(fullType[1])
	}
	return Mime(fmt.Sprintf("%s/%s", t, s))
}

func (mime Mime) SplitSubType() Mime {
	return mime.WithSubType("*")
}

func (mime Mime) WithSubType(subtype string) Mime {
	if subtype == "" {
		subtype = "*"
	}
	parts := []string{
		fmt.Sprintf("%s/%s", mime.Type(), subtype),
	}
	parts = append(parts, mime.params()...)
	return Mime(strings.Join(parts, ";"))
}

func (mime Mime) WithType(ty string) Mime {
	if ty == "" {
		ty = "*"
	}
	parts := []string{
		fmt.Sprintf("%s/%s", ty, mime.SubType()),
	}
	parts = append(parts, mime.params()...)
	return Mime(strings.Join(parts, ";"))
}

func (mime Mime) Type() string {
	return strings.TrimSpace(mime.fullType()[0])
}

func (mime Mime) SubType() string {
	fullType := mime.fullType()
	if len(fullType) == 1 {
		return "*"
	}
	return strings.TrimSpace(fullType[1])
}

func (mime Mime) HasParam(name string) bool {
	_, exists := mime.Params()[name]
	return exists
}

func (mime Mime) ParamString(name, defaultValue string) string {
	value, exists := mime.Params()[name]
	if !exists {
		value = defaultValue
	}
	return value
}

func (mime Mime) ParamInt(name string, defaultValue int) int {
	var value int
	strValue, exists := mime.Params()[name]
	if exists {
		if v, err := strconv.ParseInt(strValue, 10, 32); err == nil {
			value = int(v)
			exists = false
		}
	}
	if !exists {
		value = defaultValue
	}
	return value
}

func (mime Mime) ParamUint(name string, defaultValue uint) uint {
	var value uint
	strValue, exists := mime.Params()[name]
	if exists {
		if v, err := strconv.ParseUint(strValue, 10, 32); err == nil {
			value = uint(v)
			exists = false
		}
	}
	if !exists {
		value = defaultValue
	}
	return value
}

func (mime Mime) ParamFloat(name string, defaultValue float32) float32 {
	var value float32
	strValue, exists := mime.Params()[name]
	if exists {
		if v, err := strconv.ParseFloat(strValue, 32); err == nil {
			value = float32(v)
			exists = false
		}
	}
	if !exists {
		value = defaultValue
	}
	return value
}

func (mime Mime) ParamBool(name string, defaultValue bool) bool {
	var value bool
	strValue, exists := mime.Params()[name]
	if exists {
		if v, err := strconv.ParseBool(strValue); err == nil {
			value = v
			exists = false
		}
	}
	if !exists {
		value = defaultValue
	}
	return value
}

func (mime Mime) Params() map[string]string {
	params := make(map[string]string)
	for _, param := range mime.params() {
		row := strings.SplitN(param, "=", 2)
		if len(row) == 2 {
			params[strings.TrimSpace(row[0])] = strings.TrimSpace(row[1])
		}
	}
	return params
}

func (mime Mime) fullType() []string {
	return strings.SplitN(mime.split()[0], "/", 2)
}

func (mime Mime) params() []string {
	return mime.split()[1:]
}

func (mime Mime) split() []string {
	return strings.Split(string(mime), ";")
}
