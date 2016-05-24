// This file "float" is created by Lincan Li at 5/23/16.
// Copyright © 2016 - Lincan Li. All rights reserved

package amber

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
)

type FloatType struct {
	Float float64
	Null  bool
}

func (i FloatType) Ptr() *float64 {
	if !i.Null {
		return nil
	}
	return &i.Float
}

func (i FloatType) IsZero() bool {
	return !i.Null
}

func Float(i int64) IntegerType {
	return IntegerType{Int: int64(i), Null: false}
}

func (n *FloatType) Scan(value interface{}) error {
	v := sql.NullFloat64{Float64: n.Float, Valid: !n.Null}
	return v.Scan(value)
}

func (n FloatType) Value() (driver.Value, error) {
	if n.Null {
		return nil, nil
	}
	return n.Float, nil
}

// JSON Marshal
func (i FloatType) MarshalJSON() ([]byte, error) {
	if i.Null {
		return []byte("null"), nil
	}
	return []byte(strconv.FormatFloat(i.Float, 'f', -1, 64)), nil
}

func (i FloatType) MarshalText() ([]byte, error) {
	if i.Null {
		return []byte{}, nil
	}
	return []byte(strconv.FormatFloat(i.Float, 'f', -1, 64)), nil
}

// UnmarshalJSON implements json.Unmarshaler.
// It supports number and null input.
// 0 will not be considered a null Int.
func (f *FloatType) UnmarshalJSON(data []byte) error {
	var err error
	var v interface{}
	if err = json.Unmarshal(data, &v); err != nil {
		return err
	}
	switch x := v.(type) {
	case float64:
		f.Float = float64(x)
	case nil:
		f.Null = true
		return nil
	default:
		err = fmt.Errorf("json: cannot unmarshal %v into Go value of type null.Float", reflect.TypeOf(v).Name())
	}
	f.Null = err != nil
	return err
}
