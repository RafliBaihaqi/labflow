package lifecycle

import (
	"errors"
	"testing"
)

func TestNext(t *testing.T) {
	cases := []struct {
		name           string
		current        Status
		requested      Status
		hasSamplePhoto bool
		want           error
	}{
		{"pending -> sample received (with photo)", Pending, SampleReceived, true, nil},
		{"sample received -> test on progress", SampleReceived, TestOnProgress, false, nil},
		{"test on progress -> test completed", TestOnProgress, TestCompleted, false, nil},
		{"test completed -> coa submitted", TestCompleted, CoaSubmitted, false, nil},
		{"coa submitted -> head approved", CoaSubmitted, HeadApproved, false, nil},
		{"head approved -> order completed", HeadApproved, OrderCompleted, false, nil},

		{"pending -> rejected", Pending, Rejected, false, nil},
		{"sample received -> rejected", SampleReceived, Rejected, false, nil},
		{"test on progress -> rejected", TestOnProgress, Rejected, false, nil},
		{"test completed -> rejected", TestCompleted, Rejected, false, nil},
		{"coa submitted -> rejected", CoaSubmitted, Rejected, false, nil},
		{"head approved -> rejected", HeadApproved, Rejected, false, nil},

		{"pending -> sample received without photo", Pending, SampleReceived, false, ErrSamplePhotoRequired},

		{"sample received -> sample received", SampleReceived, SampleReceived, true, ErrInvalidTransition},
		{"test on progress -> test on progress", TestOnProgress, TestOnProgress, false, ErrInvalidTransition},

		{"pending -> test completed", Pending, TestCompleted, false, ErrInvalidTransition},
		{"pending -> test on progress", Pending, TestOnProgress, false, ErrInvalidTransition},

		{"test on progress -> sample received", TestOnProgress, SampleReceived, true, ErrInvalidTransition},

		{"order completed is terminal", OrderCompleted, Rejected, false, ErrInvalidTransition},
		{"order completed -> coa submitted is invalid", OrderCompleted, CoaSubmitted, false, ErrInvalidTransition},
		{"rejected is terminal", Rejected, TestOnProgress, false, ErrInvalidTransition},

		{"unknown current", Status("garbage"), SampleReceived, true, ErrInvalidTransition},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := Next(tc.current, tc.requested, tc.hasSamplePhoto)
			if !errors.Is(got, tc.want) {
				t.Fatalf("Next(%q, %q, %v) = %v, want %v",
					tc.current, tc.requested, tc.hasSamplePhoto, got, tc.want)
			}
		})
	}
}
