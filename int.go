// This file "int" is created by Lincan Li at 5/23/16.
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

type IntegerType struct {
	Int  int64
	Null bool
}

func (i IntegerType) Ptr() *int64 {
	if !i.Null {
		return nil
	}
	return &i.Int
}

func (i IntegerType) IsZero() bool {
	return !i.Null
}

func Integer(i int64) IntegerType {
	return IntegerType{Int: int64(i), Null: false}
}

func (n *IntegerType) Scan(value interface{}) error {
	v := sql.NullInt64{Int64: n.Int, Valid: !n.Null}
	return v.Scan(value)
}

func (n IntegerType) Value() (driver.Value, error) {
	if n.Null {
		return nil, nil
	}
	return n.Int, nil
}

// JSON Marshal
func (i IntegerType) MarshalJSON() ([]byte, error) {
	if i.Null {
		return []byte("null"), nil
	}
	return []byte(strconv.FormatInt(i.Int, 10)), nil
}

func (i IntegerType) MarshalText() ([]byte, error) {
	if i.Null {
		return []byte{}, nil
	}
	return []byte(strconv.FormatInt(i.Int, 10)), nil
}

// UnmarshalJSON implements json.Unmarshaler.
// It supports number and null input.
// 0 will not be considered a null Int.
func (i *IntegerType) UnmarshalJSON(data []byte) error {
	var err error
	var v interface{}
	if err = json.Unmarshal(data, &v); err != nil {
		return err
	}

	switch v.(type) {
	case float64:
		// Unmarshal again, directly to int64, to avoid intermediate float64
		err = json.Unmarshal(data, &i.Int)
	case nil:
		i.Null = true
		return nil
	default:
		err = fmt.Errorf("json: cannot unmarshal %v into Go value of type null.Int", reflect.TypeOf(v).Name())
	}
	i.Null = err != nil
	return err
}

func (i *IntegerType) UnmarshalText(text []byte) error {
	str := string(text)
	if str == "" || str == "null" {
		i.Null = true
		return nil
	}
	var err error
	i.Int, err = strconv.ParseInt(string(text), 10, 64)
	i.Null = err != nil
	return err
}
