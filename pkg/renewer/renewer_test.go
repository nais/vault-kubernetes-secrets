package renewer

import (
	"testing"
	"github.com/nbio/st"
)

func TestSuggestedRefreshTimeLessThanOneMinute(t *testing.T) {
	st.Assert(t, suggestedRefreshTime(20.000), 10.000)
	st.Assert(t, suggestedRefreshTime(50.000), 25.000)
}

func TestSuggestedRefreshTimeMoreThanOneMinute(t *testing.T) {
	st.Assert(t, suggestedRefreshTime(60.000), 30.000)
	st.Assert(t, suggestedRefreshTime(61.000), 31.000)
	st.Assert(t, suggestedRefreshTime(90.000), 60.000)
}
