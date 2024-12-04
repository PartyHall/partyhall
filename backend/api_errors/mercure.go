package api_errors

var MERCURE_PUBLISH_FAILURE = JsonProblem{
	Type:   "mercure-issue",
	Title:  "Failed to publish to mercure",
	Detail: "Unable to publish to mercure. Clients may not have the latest data.",
}
