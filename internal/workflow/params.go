package workflow

import (
	"encoding/json"
	"fmt"
	"reflect"
)

type parameter interface {
	validate() error
}

func unmarshalParams(marshalled []byte, params any) (err error) {
	t := reflect.TypeOf(params)
	if t.Kind() != reflect.Ptr {
		return fmt.Errorf("type %T is not a pointer", params)
	}

	err = json.Unmarshal(marshalled, params)
	if err != nil {
		return fmt.Errorf("error unmarshalling params to type %T got: %w ", params, err)
	}

	return nil
}
