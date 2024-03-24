package model

import (
	"encoding/json"
	"fmt"
	"strconv"
)

const idEncodingBase = 36

type Identifier struct {
	value int64
}

type (
	TitleId = Identifier
)

var NoId int64 = 0

func IdentifierFromString(valueStr string) (Identifier, error) {
	value, err := strconv.ParseInt(valueStr, idEncodingBase, 64)
	if err != nil {
		return Identifier{}, fmt.Errorf("invalid identifier format `%s`", valueStr)
	}
	return Identifier{value: value}, nil
}

//goland:noinspection GoMixedReceiverTypes
func (id Identifier) Int() int64 {
	return id.value
}

//goland:noinspection GoMixedReceiverTypes
func (id Identifier) String() string {
	return strconv.FormatInt(id.value, idEncodingBase)
}

//goland:noinspection GoMixedReceiverTypes
func (id Identifier) MarshalText() (text []byte, err error) {
	return []byte(id.String()), nil
}

//goland:noinspection GoMixedReceiverTypes
func (id Identifier) MarshalJSON() ([]byte, error) {
	return json.Marshal(id.String())
}

//goland:noinspection GoMixedReceiverTypes
func (id *Identifier) UnmarshalJSON(b []byte) error {
	var idString string
	if err := json.Unmarshal(b, &idString); err != nil {
		return err
	}
	idInt, err := IdentifierFromString(idString)
	if err != nil {
		return err
	}
	*id = idInt
	return nil
}
