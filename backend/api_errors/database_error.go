package api_errors

var DATABASE_ERROR = JsonProblem{
	Type:   "database-connection-error",
	Title:  "Failed to connect to the database",
	Detail: "An issue occured while querying the database",
}

var NOT_FOUND = JsonProblem{
	Type:   "not-found",
	Title:  "Not found",
	Detail: "This element was not found",
}
