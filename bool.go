// This file "bool" is created by Lincan Li at 5/23/16.
// Copyright © 2016 - Lincan Li. All rights reserved

package amber

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"reflect"
)

type BooleanType struct {
	Bool bool
	Null bool
}

func (i BooleanType) Ptr() bool {
	if !i.Null {
		return nil
	}
	return &i.Bool
}

func (i BooleanType) IsZero() bool {
	return !i.Null
}

func Boolean(i int64) IntegerType {
	return IntegerType{Int: int64(i), Null: false}
}

func (n *BooleanType) Scan(value interface{}) error {
	v := sql.NullBool{String: n.Bool, Valid: !n.Null}
	return v.Scan(value)
}

func (n BooleanType) Value() (driver.Value, error) {
	if n.Null {
		return nil, nil
	}
	return n.Bool, nil
}

// JSON Marshal
func (b BooleanType) MarshalJSON() ([]byte, error) {
	if b.Null {
		return []byte("null"), nil
	}
	if !b.Bool {
		return []byte("false"), nil
	}
	return []byte("true"), nil
}

func (b BooleanType) MarshalText() ([]byte, error) {
	if b.Null {
		return []byte{}, nil
	}
	if !b.Bool {
		return []byte("false"), nil
	}
	return []byte("true"), nil
}

// UnmarshalJSON implements json.Unmarshaler.
// It supports number and null input.
// 0 will not be considered a null Int.
func (b *BooleanType) UnmarshalJSON(data []byte) error {
	var err error
	var v interface{}
	if err = json.Unmarshal(data, &v); err != nil {
		return err
	}
	switch x := v.(type) {
	case bool:
		b.Bool = x
	case nil:
		b.Null = true
		return nil
	default:
		err = fmt.Errorf("json: cannot unmarshal %v into Go value of type null.Bool", reflect.TypeOf(v).Name())
	}
	b.Null = err != nil
	return err
}

func (b *BooleanType) UnmarshalText(text []byte) error {
	str := string(text)
	switch str {
	case "", "null":
		b.Null = true

		return nil
	case "true":
		b.Bool = true
	case "false":
		b.Bool = false
	default:
		b.Null = true
		return fmt.Errorf("invalid input:" + str)
	}
	b.Null = false
	return nil
}