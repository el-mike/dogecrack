package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
)

// NullableTime - an alias for time.Time struct, providing
// nullable functionality for JSON.
type NullableTime struct {
	time.Time `bson:",inline"`
}

// NullableTimeNow - returns new NullableTime with time.Now().
func NullableTimeNow() NullableTime {
	return NullableTime{time.Now()}
}

// MarshalJSON - JSON Marshaler implementation, for handling time.Time
// zero values. It changes time's zero value to "null".
func (nt NullableTime) MarshalJSON() ([]byte, error) {
	if nt.IsZero() {
		return []byte("null"), nil
	}

	return nt.Time.MarshalJSON()
}

// UnarshalJSON - JSON Unmarshaler implementation.
func (nt *NullableTime) UnmarshalJSON(data []byte) error {
	var time time.Time

	if err := time.UnmarshalJSON(data); err != nil {
		return nil
	}

	*nt = NullableTime{time}

	return nil
}

// MarshalBSONValue - ValueMarshaler implementation. Helps with converting
// between NullableTime and basic Time type.
func (nt *NullableTime) MarshalBSONValue() (bsontype.Type, []byte, error) {
	return bson.MarshalValue(nt.Time)
}

// MarshalBSONValue - ValueUnmarshaler implementation. Helps with converting
// between NullableTime and basic Time type.
func (nt *NullableTime) UnmarshalBSONValue(t bsontype.Type, data []byte) error {
	raw := bson.RawValue{Type: t, Value: data}

	return raw.Unmarshal(&nt.Time)
}
