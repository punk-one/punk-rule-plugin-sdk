package sdk_test

import (
	"testing"

	sdk "github.com/punk-one/punk-rule-plugin-sdk"
)

func TestEventCloneProducesIndependentMaps(t *testing.T) {
	event := sdk.NewEvent(map[string]interface{}{
		"device": map[string]interface{}{
			"id": "plc-1",
		},
		"values": []interface{}{1, "x"},
	}, map[string]string{
		"source": "s7",
	})
	sdk.SetTraceID(&event, "trace-1")

	cloned := event.Clone()

	cloned.Payload["device"].(map[string]interface{})["id"] = "plc-2"
	cloned.Payload["values"].([]interface{})[0] = 2
	cloned.Metadata["source"] = "modbus"
	sdk.SetTraceID(&cloned, "trace-2")

	if event.Payload["device"].(map[string]interface{})["id"] != "plc-1" {
		t.Fatalf("expected original nested payload map to stay unchanged")
	}
	if event.Payload["values"].([]interface{})[0] != 1 {
		t.Fatalf("expected original payload slice to stay unchanged")
	}
	if event.Metadata["source"] != "s7" {
		t.Fatalf("expected original metadata to stay unchanged")
	}
	if sdk.GetTraceID(event) != "trace-1" {
		t.Fatalf("expected original trace id to stay unchanged")
	}
}
