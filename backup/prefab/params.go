package prefab

import (
	"encoding/xml"
	"fmt"
)

type XMLParam struct {
	Key   string
	Value string
}

func (param *XMLParam) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var ready bool
	var i interface{} = start
	var err error
	for ; err == nil; i, err = d.Token() {
		if i == start.End() {
			return nil
		}
		switch token := i.(type) {
		case xml.StartElement:
			if !ready {
				param.Key = token.Name.Local
				ready = true
			} else {
				return fmt.Errorf("unexpected token %s", token.Name.Local)
			}
		case xml.EndElement:
			if fmt.Sprintf("/%s", param.Key) != token.Name.Local {
				return fmt.Errorf("unexpected token %s", token)
			}
		case xml.CharData:
			if ready == true {
				param.Value += string(token)
			}
		}
	}
	return fmt.Errorf("unexpected end of file")
}
