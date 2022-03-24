package main

import (
	"testing"
	"time"
)

func TestHumanDate(t *testing.T) {
	tests := []struct {
		name   string
		input  time.Time
		expect string
	}{
		{
			name:   "UTC",
			input:  time.Date(2020, 12, 17, 10, 0, 0, 0, time.UTC),
			expect: "17 Dec 2020 at 10:00",
		},
		{
			name:   "Empty",
			input:  time.Time{},
			expect: "",
		},
		{
			name:   "CET",
			input:  time.Date(2020, 12, 17, 10, 0, 0, 0, time.FixedZone("CET", 1*60*60)),
			expect: "17 Dec 2020 at 09:00",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hd := humanDate(tt.input)

			if hd != tt.expect {
				t.Errorf("expecting %q; got %q", tt.expect, hd)
			}
		})
	}
}
