package api_errors

var UNAUTHORIZED = JsonProblem{
	Type:   "unauthorized",
	Title:  "Unauthorized",
	Detail: "The user does not have the permissions to access this route",
}
