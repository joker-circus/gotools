package models

import (
	"strconv"
	"testing"
	"time"
)

func TestComparators(t *testing.T) {
	t1a := TimeFromUnix(0)
	t1b := TimeFromUnix(0)
	t2 := TimeFromUnix(2*second - 1)

	if !t1a.Equal(t1b) {
		t.Fatalf("Expected %s to be equal to %s", t1a, t1b)
	}
	if t1a.Equal(t2) {
		t.Fatalf("Expected %s to not be equal to %s", t1a, t2)
	}

	if !t1a.Before(t2) {
		t.Fatalf("Expected %s to be before %s", t1a, t2)
	}
	if t1a.Before(t1b) {
		t.Fatalf("Expected %s to not be before %s", t1a, t1b)
	}

	if !t2.After(t1a) {
		t.Fatalf("Expected %s to be after %s", t2, t1a)
	}
	if t1b.After(t1a) {
		t.Fatalf("Expected %s to not be after %s", t1b, t1a)
	}
}

func TestTimeConversions(t *testing.T) {
	unixSecs := int64(1136239445)
	unixNsecs := int64(123456789)
	unixNano := unixSecs*1e9 + unixNsecs

	t1 := time.Unix(unixSecs, unixNsecs-unixNsecs%nanosPerTick)
	t2 := time.Unix(unixSecs, unixNsecs)

	ts := TimeFromUnixNano(unixNano)
	if !ts.Time().Equal(t1) {
		t.Fatalf("Expected %s, got %s", t1, ts.Time())
	}

	// Test available precision.
	ts = TimeFromUnixNano(t2.UnixNano())
	if !ts.Time().Equal(t1) {
		t.Fatalf("Expected %s, got %s", t1, ts.Time())
	}

	if ts.UnixNano() != unixNano-unixNano%nanosPerTick {
		t.Fatalf("Expected %d, got %d", unixNano, ts.UnixNano())
	}
}

func TestDuration(t *testing.T) {
	duration := time.Second + time.Minute + time.Hour
	goTime := time.Unix(1136239445, 0)

	ts := TimeFromUnix(goTime.Unix())
	if !goTime.Add(duration).Equal(ts.Add(duration).Time()) {
		t.Fatalf("Expected %s to be equal to %s", goTime.Add(duration), ts.Add(duration))
	}

	earlier := ts.Add(-duration)
	delta := ts.Sub(earlier)
	if delta != duration {
		t.Fatalf("Expected %s to be equal to %s", delta, duration)
	}
}

func TestTimeJSON(t *testing.T) {
	tests := []struct {
		in  Time
		out string
	}{
		{Time(1), `0.001`},
		{Time(-1), `-0.001`},
	}

	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			b, err := test.in.MarshalJSON()
			if err != nil {
				t.Fatalf("Error marshaling time: %v", err)
			}

			if string(b) != test.out {
				t.Errorf("Mismatch in marshal expected=%s actual=%s", test.out, b)
			}

			var tm Time
			if err := tm.UnmarshalJSON(b); err != nil {
				t.Fatalf("Error Unmarshaling time: %v", err)
			}

			if !test.in.Equal(tm) {
				t.Fatalf("Mismatch after Unmarshal expected=%v actual=%v", test.in, tm)
			}

		})
	}

}
