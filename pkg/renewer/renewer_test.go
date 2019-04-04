package renewer

import (
	"testing"
	"github.com/nbio/st"
)

func TestSuggestedRefreshTimeLessThanTenMinutes(t *testing.T) {
	st.Assert(t, suggestedRefreshTime(20.000), 10.000)
	st.Assert(t, suggestedRefreshTime(50.000), 25.000)
}

func TestSuggestedRefreshTimeMoreThanTenMinutse(t *testing.T) {
	st.Assert(t, suggestedRefreshTime(600.000), 300.000)
	st.Assert(t, suggestedRefreshTime(610.000), 310.000)
	st.Assert(t, suggestedRefreshTime(900.000), 600.000)
}
