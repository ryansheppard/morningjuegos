package utils

import "testing"

func TestGetTimeBounday(t *testing.T) {
	t.Parallel()

	start, end := GetTimeBoundary(1695170605)
	wantStart, wantEnd := int64(1695096000), int64(1695182399)

	if start != wantStart {
		t.Errorf("start = %d, want %d", start, wantStart)
	}

	if end != wantEnd {
		t.Errorf("end = %d, want %d", end, wantEnd)
	}
}
