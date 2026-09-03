package cmd

import (
	"testing"
	"time"
)

func TestParseRepeatUnit(t *testing.T) {
	tests := []struct {
		in      string
		want    repeatUnit
		wantErr bool
	}{
		{"", repeatNone, false},
		{"daily", repeatDaily, false},
		{"day", repeatDaily, false},
		{"weekly", repeatWeekly, false},
		{"monthly", repeatMonthly, false},
		{"yearly", repeatYearly, false},
		{"annually", repeatYearly, false},
		{"hourly", repeatNone, true},
	}

	for _, tt := range tests {
		got, err := parseRepeatUnit(tt.in)
		if tt.wantErr {
			if err == nil {
				t.Errorf("parseRepeatUnit(%q) expected error", tt.in)
			}
			continue
		}
		if err != nil {
			t.Errorf("parseRepeatUnit(%q): %v", tt.in, err)
			continue
		}
		if got != tt.want {
			t.Errorf("parseRepeatUnit(%q) = %v, want %v", tt.in, got, tt.want)
		}
	}
}

func TestExpandRepeatScheduleCount(t *testing.T) {
	start := time.Date(2025, 2, 10, 0, 0, 0, 0, time.UTC)
	got := expandRepeatSchedule(start, repeatDaily, 5, nil)
	if len(got) != 5 {
		t.Fatalf("got %d dates, want 5", len(got))
	}
	if !got[0].Equal(start) || !got[4].Equal(start.AddDate(0, 0, 4)) {
		t.Fatalf("unexpected range: %v .. %v", got[0], got[4])
	}
}
