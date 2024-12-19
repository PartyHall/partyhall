package api_errors

var MQTT_PUBLISH_FAILURE = JsonProblem{
	Type:   "mqtt-issue",
	Title:  "Failed to publish to MQTT",
	Detail: "Unable to publish to MQTT.",
}
