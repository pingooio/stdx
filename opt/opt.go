// package opt provides optional types. It can be useful for optional configuration paramteters
package opt

import "time"

func String(str string) *string {
	return &str
}

func Int(i int) *int {
	return &i
}

func Int8(i int8) *int8 {
	return &i
}

func Int16(i int16) *int16 {
	return &i
}

func Int32(i int32) *int32 {
	return &i
}

func Int64(i int64) *int64 {
	return &i
}

func Uint(i uint) *uint {
	return &i
}

func Uint8(i uint8) *uint8 {
	return &i
}

func Uint16(i uint16) *uint16 {
	return &i
}

func Uint32(i uint32) *uint32 {
	return &i
}

func Uint64(i uint64) *uint64 {
	return &i
}

func Time(t time.Time) *time.Time {
	return &t
}

func Bool(v bool) *bool {
	return &v
}
