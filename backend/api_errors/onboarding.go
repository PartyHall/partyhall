package api_errors

var SHOULD_BE_ONBOARDED = JsonProblem{
	Type:   "not-onboarded",
	Title:  "Not onboarded",
	Detail: "This appliance is not setup yet",
}

var ALREADY_ONBOARDED = JsonProblem{
	Type:   "already-onboarded",
	Title:  "Already onboarded",
	Detail: "This appliance is already setup",
}
