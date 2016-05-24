// This file "string" is created by Lincan Li at 5/23/16.
// Copyright © 2016 - Lincan Li. All rights reserved

package amber

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"reflect"
)

type StringType struct {
	String string
	Null   bool
}

func (i StringType) Ptr() *string {
	if !i.Null {
		return nil
	}
	return &i.String
}

func (i StringType) IsZero() bool {
	return !i.Null
}

func String(i int64) IntegerType {
	return IntegerType{Int: int64(i), Null: false}
}

func (n *StringType) Scan(value interface{}) error {
	v := sql.NullString{String: n.String, Valid: !n.Null}
	return v.Scan(value)
}

func (n StringType) Value() (driver.Value, error) {
	if n.Null {
		return nil, nil
	}
	return n.String, nil
}

// JSON Marshal
func (i StringType) MarshalJSON() ([]byte, error) {
	if i.Null {
		return []byte("null"), nil
	}
	return json.Marshal(i.String)
}

func (i StringType) MarshalText() ([]byte, error) {
	if i.Null {
		return []byte{}, nil
	}
	return []byte(i.String), nil
}

// UnmarshalJSON implements json.Unmarshaler.
// It supports number and null input.
// 0 will not be considered a null Int.
func (s *StringType) UnmarshalJSON(data []byte) error {
	var err error
	var v interface{}
	if err = json.Unmarshal(data, &v); err != nil {
		return err
	}
	switch x := v.(type) {
	case string:
		s.String = x
	case nil:
		s.Null = true
		return nil
	default:
		err = fmt.Errorf("json: cannot unmarshal %v into Go value of type null.String", reflect.TypeOf(v).Name())
	}
	s.Null = err != nil
	return err
}
