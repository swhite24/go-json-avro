package converter

import (
	"encoding/json"
	"fmt"
	"testing"
)

var ()

func TestFoo(t *testing.T) {
	schemas := []string{
		"{\"type\":\"record\",\"name\":\"TestSchema\",\"namespace\":\"com.punisher\",\"fields\":[{\"name\":\"name\",\"type\":\"string\"},{\"name\":\"number\",\"type\":[\"null\",\"int\"],\"default\":null},{\"name\":\"list\",\"type\":\"array\",\"items\":\"string\"}]}",
		"{\"type\":\"record\",\"name\":\"InStoreNavigationProcessingCompletedEvent\",\"namespace\":\"com.kroger.desp.events.fulfillment\",\"doc\":\"Event fired when in-store navigation processing completes\",\"fields\":[{\"name\":\"eventHeader\",\"type\":{\"type\":\"record\",\"name\":\"EventHeader\",\"namespace\":\"com.kroger.desp.commons\",\"doc\":\"The below fields include header information and should be included on every event in the DESP. Inspired by: https://github.com/cloudevents/spec/blob/v0.2/spec.md\",\"fields\":[{\"name\":\"id\",\"type\":{\"type\":\"string\",\"avro.java.string\":\"String\"},\"doc\":\"A unique identifier of the event - for example, a randomly generated GUID\"},{\"name\":\"time\",\"type\":\"long\",\"doc\":\"Time the event occurred in milliseconds since epoch, UTC timezone.\"},{\"name\":\"type\",\"type\":{\"type\":\"string\",\"avro.java.string\":\"String\"},\"doc\":\"Type of occurrence which has happened. Reference the domain.event registered in schema-registry.\"},{\"name\":\"source\",\"type\":{\"type\":\"string\",\"avro.java.string\":\"String\"},\"doc\":\"Service that produced the event. Future: reference to producer registry.\"}]}},{\"name\":\"summary\",\"type\":{\"type\":\"record\",\"name\":\"InStoreNavigationProcessSummary\",\"namespace\":\"com.kroger.desp.commons.fulfillment\",\"fields\":[{\"name\":\"numProcessed\",\"type\":\"int\",\"doc\":\"Number of input items processed during the request\"},{\"name\":\"timeBegin\",\"type\":\"long\",\"doc\":\"A millisecond timestamp indicating when processing began\"},{\"name\":\"timeEnd\",\"type\":\"long\",\"doc\":\"A millisecond timestamp indicating when processing ended\"}]}}],\"version\":\"1\"}",
	}
	inputs := []map[string]interface{}{
		map[string]interface{}{
			"name":   "foo",
			"number": 5,
			"list":   []string{"foo", "bar"},
		},
		map[string]interface{}{
			"eventHeader":     map[string]interface{}{},
			"summary":         map[string]interface{}{},
			"itemDescription": "foobar",
		},
		map[string]interface{}{
			"eventHeader":     map[string]interface{}{},
			"itemDescription": "foobar",
		},
	}

	for i, schema := range schemas {
		input := inputs[i]
		data, err := Convert(input, schema)
		fmt.Println("out: ", data, err)
		out, _ := json.Marshal(data)
		fmt.Println(string(out))
	}

	t.Fail()
}
