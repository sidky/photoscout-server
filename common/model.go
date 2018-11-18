package common

import (
	"encoding/json"
	"strconv"
)

type FlexInt struct {
	value int32
}

func (v *FlexInt) UnmarshalJSON(bytes []byte) error {
	s := string(bytes)
	if l := len(s); l > 1 && s[0] == '"' && s[l-1] == '"' {
		s = s[1 : l-1]
	}
	value, err := strconv.Atoi(s)
	if err != nil {
		return err
	}
	v.value = int32(value)
	return nil
}

func (v *FlexInt) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *FlexInt) Int() *int32 {
	if v == nil {
		return nil
	} else {
		return &v.value
	}
}

type FlexFloat struct {
	value float64
}

func (v *FlexFloat) UnmarshalJSON(bytes []byte) error {
	s := string(bytes)
	if l := len(s); l > 1 && s[0] == '"' && s[l-1] == '"' {
		s = s[1 : l-1]
	}
	value, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return err
	}
	v.value = value
	return nil
}

func (v *FlexFloat) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *FlexFloat) Float() float64 {
	return v.value
}
