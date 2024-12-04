package api_errors

var INVALID_PARAMETERS = JsonProblem{
	Type:   "invalid-parameters",
	Title:  "Invalid parameters",
	Detail: "The query params are invalid",
}
