package api_errors

var DATABASE_ERROR = JsonProblem{
	Type:   "database-connection-error",
	Title:  "Failed to connect to the database",
	Detail: "An issue occured while querying the database",
}
