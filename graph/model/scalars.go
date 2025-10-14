package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/99designs/gqlgen/graphql"
	"github.com/shopspring/decimal"
)

// Type aliases for GraphQL scalars
type Time = time.Time
type JSON = map[string]interface{}
type Decimal = decimal.Decimal

// MarshalTime marshals time.Time to RFC3339 string for GraphQL
func MarshalTime(t time.Time) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		w.Write([]byte(`"`))
		w.Write([]byte(t.Format(time.RFC3339)))
		w.Write([]byte(`"`))
	})
}

// UnmarshalTime unmarshals RFC3339 string to time.Time for GraphQL
func UnmarshalTime(v interface{}) (time.Time, error) {
	if tmpStr, ok := v.(string); ok {
		return time.Parse(time.RFC3339, tmpStr)
	}
	return time.Time{}, errors.New("time must be a string")
}

// MarshalJSON marshals map[string]interface{} to JSON for GraphQL
func MarshalJSON(m map[string]interface{}) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		if err := json.NewEncoder(w).Encode(m); err != nil {
			panic(err)
		}
	})
}

// UnmarshalJSON unmarshals JSON to map[string]interface{} for GraphQL
func UnmarshalJSON(v interface{}) (map[string]interface{}, error) {
	m, ok := v.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("JSON must be an object")
	}
	return m, nil
}

// MarshalDecimal marshals decimal.Decimal to string for GraphQL
func MarshalDecimal(d decimal.Decimal) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		w.Write([]byte(`"`))
		w.Write([]byte(d.String()))
		w.Write([]byte(`"`))
	})
}

// UnmarshalDecimal unmarshals string/number to decimal.Decimal for GraphQL
func UnmarshalDecimal(v interface{}) (decimal.Decimal, error) {
	var dec decimal.Decimal
	var err error

	switch v := v.(type) {
	case string:
		dec, err = decimal.NewFromString(v)
		if err != nil {
			return decimal.Zero, err
		}
	case int:
		dec = decimal.NewFromInt(int64(v))
	case int64:
		dec = decimal.NewFromInt(v)
	case float64:
		dec = decimal.NewFromFloat(v)
	default:
		return decimal.Zero, fmt.Errorf("decimal must be a string or number, got %T", v)
	}

	return dec, nil
}

// Type alias for BigInt scalar
type BigInt = int64

// MarshalBigInt marshals int64 to number for GraphQL
func MarshalBigInt(i int64) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		w.Write([]byte(fmt.Sprintf("%d", i)))
	})
}

// UnmarshalBigInt unmarshals number/string to int64 for GraphQL
func UnmarshalBigInt(v interface{}) (int64, error) {
	switch v := v.(type) {
	case string:
		var i int64
		_, err := fmt.Sscanf(v, "%d", &i)
		if err != nil {
			return 0, err
		}
		return i, nil
	case int:
		return int64(v), nil
	case int32:
		return int64(v), nil
	case int64:
		return v, nil
	case float64:
		return int64(v), nil
	default:
		return 0, fmt.Errorf("BigInt must be a number or string, got %T", v)
	}
}
