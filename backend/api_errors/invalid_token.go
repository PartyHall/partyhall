package api_errors

var INVALID_TOKEN = JsonProblem{
	Type:   "invalid-token",
	Title:  "Invalid authentication token",
	Detail: "The given token is invalid",
}

var NO_TOKEN = JsonProblem{
	Type:   "no-token",
	Title:  "No authorization token",
	Detail: "The request did no contain any authorization token",
}
