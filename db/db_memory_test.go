package db

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewDbInMemory(t *testing.T) {
	rawDb := `{"key1": "value1", "key2": "value2"}`
	db, err := NewDbInMemory(rawDb)

	assert.NoError(t, err)
	assert.NotNil(t, db)
}

func TestReadEntry(t *testing.T) {
	rawDb := `{"key1": "value1", "key2": "value2"}`
	db, _ := NewDbInMemory(rawDb)

	ok, value, err := ReadEntry[string](db, "key1")

	assert.True(t, ok)
	assert.Equal(t, "value1", value)
	assert.NoError(t, err)
}

func TestWriteEntry(t *testing.T) {
	rawDb := `{"key1": "value1", "key2": "value2"}`
	db, _ := NewDbInMemory(rawDb)

	err := WriteEntry(db, "key3", TestStruct{A: "value3", B: 3})

	assert.NoError(t, err)

	ok, value, err := ReadEntry[TestStruct](db, "key3")

	assert.True(t, ok)
	assert.Equal(t, TestStruct{A: "value3", B: 3}, value)
	assert.NoError(t, err)
}

func TestToString(t *testing.T) {
	rawDb := `{"key1": "value1", "key2": "value2"}`
	db, _ := NewDbInMemory(rawDb)

	err := WriteEntry(db, "key3", TestStruct{A: "value3", B: 3})

	assert.NoError(t, err)

	value, err := db.String()

	assert.Equal(t, `{"key1":"value1","key2":"value2","key3":{"A":"value3","B":3}}`, string(value))
	assert.NoError(t, err)
}

type TestStruct struct {
	A string
	B int
}
