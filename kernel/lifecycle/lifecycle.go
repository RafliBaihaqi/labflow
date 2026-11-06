package lifecycle

import "errors"

type Status string

const (
	Pending        Status = "pending"
	SampleReceived Status = "sample_received"
	TestOnProgress Status = "test_on_progress"
	TestCompleted  Status = "test_completed"
	CoaSubmitted   Status = "coa_submitted"
	HeadApproved   Status = "head_approved"
	OrderCompleted Status = "order_completed"
	Rejected       Status = "rejected"
)

var (
	ErrInvalidTransition   = errors.New("invalid status transition")
	ErrSamplePhotoRequired = errors.New("sample photo required to mark sample received")
)

var forward = map[Status][]Status{
	Pending:        {SampleReceived, Rejected},
	SampleReceived: {TestOnProgress, Rejected},
	TestOnProgress: {TestCompleted, Rejected},
	TestCompleted:  {CoaSubmitted, Rejected},
	CoaSubmitted:   {HeadApproved, Rejected},
	HeadApproved:   {OrderCompleted, Rejected},
}

// Transitions returns a copy of the forward transition map for API exposure
func Transitions() map[Status][]Status {
	cp := make(map[Status][]Status, len(forward))
	for k, v := range forward {
		dst := make([]Status, len(v))
		copy(dst, v)
		cp[k] = dst
	}
	return cp
}

// AllStatuses returns all known statuses
func AllStatuses() []Status {
	return []Status{
		Pending,
		SampleReceived,
		TestOnProgress,
		TestCompleted,
		CoaSubmitted,
		HeadApproved,
		OrderCompleted,
		Rejected,
	}
}

func IsTerminal(s Status) bool {
	return s == OrderCompleted || s == Rejected
}

func Next(current, requested Status, hasSamplePhoto bool) error {
	allowed, ok := forward[current]
	if !ok {
		return ErrInvalidTransition
	}

	found := false
	for _, s := range allowed {
		if s == requested {
			found = true
			break
		}
	}
	if !found {
		return ErrInvalidTransition
	}

	if requested == SampleReceived && !hasSamplePhoto {
		return ErrSamplePhotoRequired
	}

	return nil
}
