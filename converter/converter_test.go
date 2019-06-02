package converter

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSimpleSchema(t *testing.T) {
	schema := "{\"type\":\"record\",\"name\":\"TestSchema\",\"namespace\":\"com.punisher\",\"fields\":[{\"name\":\"name\",\"type\":\"string\"},{\"name\":\"number\",\"type\":\"int\"},{\"name\":\"list\",\"type\":\"array\",\"items\":\"string\"},{\"name\":\"map\",\"type\":\"map\",\"values\":\"string\"}]}"
	input := map[string]interface{}{
		"name":   "foo",
		"number": 10,
		"list":   []string{"foo", "bar"},
		"map": map[string]string{
			"foo": "bar",
		},
	}

	out, err := Convert(input, schema)

	assert.Nil(t, err)
	assert.Equal(t, input["name"], out["name"])
	assert.Equal(t, input["number"], out["number"])
	assert.Equal(t, input["list"], out["list"])
	assert.Equal(t, input["map"], out["map"])
}

func TestSchemaRecordField(t *testing.T) {
	schema := "{\"type\":\"record\",\"name\":\"TestSchema\",\"namespace\":\"com.punisher\",\"fields\":[{\"name\":\"name\",\"type\":\"string\"},{\"name\":\"nested\",\"type\":{\"type\":\"record\",\"name\":\"NestedType\",\"fields\":[{\"name\":\"id\",\"type\":\"string\"},{\"name\":\"val\",\"type\":\"string\"}]}}]}"
	input := map[string]interface{}{
		"name": "foo",
		"nested": map[string]interface{}{
			"id":  "1234",
			"val": "abcd",
		},
	}

	out, err := Convert(input, schema)
	assert.Nil(t, err)
	assert.Equal(t, input["name"], out["name"])
	assert.Equal(t, input["nested"], out["nested"])
}

func TestSchemaUnionField(t *testing.T) {
	schema := "{\"type\":\"record\",\"name\":\"TestSchema\",\"namespace\":\"com.punisher\",\"fields\":[{\"name\":\"name\",\"type\":\"string\"},{\"name\":\"number\",\"type\":[\"null\",\"int\"],\"default\":null}]}"
	input := map[string]interface{}{
		"name": "foo",
	}

	out, err := Convert(input, schema)
	assert.Nil(t, err)
	assert.Equal(t, input["name"], out["name"])
	assert.Nil(t, input["number"])

	input = map[string]interface{}{
		"name":   "foo",
		"number": 10,
	}
	out, err = Convert(input, schema)
	assert.Nil(t, err)
	assert.Equal(t, input["name"], out["name"])
	assert.NotNil(t, out["number"])
	numval, _ := out["number"].(map[string]interface{})
	assert.Equal(t, input["number"], numval["int"])
}

func TestSchemaUnionRecordField(t *testing.T) {
	schema := "{\"type\":\"record\",\"name\":\"TestSchema\",\"namespace\":\"com.punisher\",\"fields\":[{\"name\":\"name\",\"type\":\"string\"},{\"name\":\"nested\",\"type\":[\"null\",{\"type\":\"record\",\"name\":\"NestedType\",\"namespace\":\"com.punisher\",\"fields\":[{\"name\":\"id\",\"type\":\"string\"},{\"name\":\"val\",\"type\":\"string\"}]}],\"default\":null}]}"
	input := map[string]interface{}{
		"name": "foo",
		"nested": map[string]interface{}{
			"id":  "abcd",
			"val": "1234",
		},
	}

	out, err := Convert(input, schema)
	assert.Nil(t, err)
	assert.Equal(t, input["name"], out["name"])
	assert.NotNil(t, out["nested"])
	nestedval, _ := out["nested"].(map[string]interface{})
	assert.Equal(t, nestedval["com.punisher.NestedType"], input["nested"])
}

func TestSchemaExpandedPrimitiveField(t *testing.T) {
	schema := "{\"type\":\"record\",\"name\":\"TestSchema\",\"namespace\":\"com.punisher\",\"fields\":[{\"name\":\"name\",\"type\":{\"type\":\"string\",\"avro.java.string\":\"String\"}},{\"name\":\"number\",\"type\":\"int\"},{\"name\":\"list\",\"type\":\"array\",\"items\":\"string\"},{\"name\":\"map\",\"type\":\"map\",\"values\":\"string\"}]}"
	input := map[string]interface{}{
		"name":   "foo",
		"number": 10,
	}

	out, err := Convert(input, schema)
	assert.Nil(t, err)
	assert.Equal(t, input["name"], out["name"])
	assert.Equal(t, input["number"], out["number"])
}
