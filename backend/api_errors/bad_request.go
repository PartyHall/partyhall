package api_errors

var BAD_REQUEST = JsonProblem{
	Type:   "bad-request",
	Title:  "Bad request",
	Detail: "The request body is invalid",
}
