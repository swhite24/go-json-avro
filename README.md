# go-json-avro-converter

Utility to convert traditional JSON data into avro JSON data.

Union types in avro require a special encoding, see http://avro.apache.org/docs/1.9.0/spec.html#json_encoding.

Loosely inspired by https://github.com/allegro/json-avro-converter

## Installation

```shell
go get github.com/swhite24/go-json-avro-converter/converter
```

## Usage

```go
package main

import (
  "github.com/swhite24/go-json-avro-converter/converter"
)

var schema = `
{
  "type": "record",
  "name": "TestSchema",
  "namespace": "com.punisher",
  "fields": [
    {
      "name": "name",
      "type": "string"
    },
    {
      "name": "nested",
      "type": [
        "null",
        {
          "type": "record",
          "name": "NestedType",
          "namespace": "com.punisher",
          "fields": [
            {
              "name": "id",
              "type": "string"
            },
            {
              "name": "val",
              "type": "string"
            }
          ]
        }
      ],
      "default": null
    }
  ]
}
`

func main() {
  input := map[string]interface{}{
    "name": "foo",
    "nested": map[string]interface{}{
      "id": "1234",
      "val": "abcd",
    },
  }

  out, err := converter.Convert(input, schema)
  fmt.Println(out, err)

  // map[name:foo nested:map[com.punisher.NestedType:map[id:1234 val:abcd]]] <nil>
}
```
