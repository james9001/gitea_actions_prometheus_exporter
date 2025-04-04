package collector

// Status represents the status of ActionRun, ActionRunJob, ActionTask, or ActionTaskStep.
type Status int

// These are the possible statuses of an ActionRun.
const (
	StatusUnknown   Status = iota // 0, consistent with runnerv1.Result_RESULT_UNSPECIFIED
	StatusSuccess                 // 1, consistent with runnerv1.Result_RESULT_SUCCESS
	StatusFailure                 // 2, consistent with runnerv1.Result_RESULT_FAILURE
	StatusCancelled               // 3, consistent with runnerv1.Result_RESULT_CANCELLED
	StatusSkipped                 // 4, consistent with runnerv1.Result_RESULT_SKIPPED
	StatusWaiting                 // 5, isn't a runnerv1.Result
	StatusRunning                 // 6, isn't a runnerv1.Result
	StatusBlocked                 // 7, isn't a runnerv1.Result
)

func getStatusNames() map[Status]string {
	return map[Status]string{
		StatusUnknown:   "unknown",
		StatusWaiting:   "waiting",
		StatusRunning:   "running",
		StatusSuccess:   "success",
		StatusFailure:   "failure",
		StatusCancelled: "cancelled",
		StatusSkipped:   "skipped",
		StatusBlocked:   "blocked",
	}
}

// String returns the string name of the Status.
func (s Status) String() string {
	return getStatusNames()[s]
}

// MarshalJSON implements the json.Marshaler interface.
// It converts the Status to its string representation when serializing to JSON.
func (s Status) MarshalJSON() ([]byte, error) {
	return []byte(`"` + s.String() + `"`), nil
}

// IsDone returns whether the Status is final.
func (s Status) IsDone() bool {
	return s.In(StatusSuccess, StatusFailure, StatusCancelled, StatusSkipped)
}

// HasRun returns whether the Status is a result of running.
func (s Status) HasRun() bool {
	return s.In(StatusSuccess, StatusFailure)
}

// In returns whether s is one of the given statuses.
func (s Status) In(statuses ...Status) bool {
	for _, v := range statuses {
		if s == v {
			return true
		}
	}

	return false
}
