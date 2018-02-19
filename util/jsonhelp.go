package util

import (
	"encoding/json"
	"reflect"
	"strings"
)

func UnmarshalJsonObjAndMap(jsonStr []byte, obj interface{}, otherFields map[string]json.RawMessage) (err error) {
	objValue := reflect.ValueOf(obj).Elem()
	knownFields := map[string]reflect.Value{}
	for i := 0; i != objValue.NumField(); i++ {
		jsonName := strings.Split(objValue.Type().Field(i).Tag.Get("json"), ",")[0]
		knownFields[jsonName] = objValue.Field(i)
	}
	err = json.Unmarshal(jsonStr, &otherFields)
	if err != nil {
		return err
	}

	for key, chunk := range otherFields {
		if field, found := knownFields[key]; found {
			err = json.Unmarshal(chunk, field.Addr().Interface())
			if err != nil {
				return err
			}
			delete(otherFields, key)
		}
	}
	return err
}

func MarshalJsonObjAndMap(obj interface{}, otherFields map[string]json.RawMessage) ([]byte, error) {
	a, err := json.Marshal(otherFields)
	if err != nil {
		return nil, err
	}
	b, err := json.Marshal(&obj)
	if err != nil {
		return nil, err
	}
	// the smallest json is {}
	if len(a) > 2 {
		a[len(a)-1] = ','
	} else {
		a = a[:1]
	}
	b = b[1:]

	c := append(a, b...)
	return c, nil
}
