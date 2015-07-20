// Copyright 2015 Yoshi Yamaguchi
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package datemaki

import (
	"strings"
	"testing"
	"time"
)

var agoTests = []string{
	"2 seconds ago",
	"3 minutes ago",
	"4 hours ago",
	"5 days ago",
	"1 week ago",
	"2 months ago",
	"1 year, 3 months ago",
	"1.year.4.months.ago",
	"2.years.ago",
}

func TestSplitTokens(t *testing.T) {
	for i, test := range agoTests {
		pre1 := strings.Replace(test, ",", " ", -1)
		pre2 := strings.Replace(pre1, ".", " ", -1)
		words := strings.Fields(pre2)
		tokens := splitTokens(test)
		if len(words) != len(tokens) {
			t.Errorf("#%d: word counts are different, %d is expected, got %d", i, len(words), len(tokens))
			continue
		}
	}
}

func TestParseAgo(t *testing.T) {
	for i, test := range agoTests {
		parsed, err := ParseAgo(test)
		if err != nil {
			t.Errorf("#%v: %v", i, err)
			continue
		}
		t.Logf("#%v: parsed: %v (%v)", i, parsed, test)
	}
}

var relativeTests = []string{
	"now",
	"today",
	"yesterday",
	"last friday",
	"noon yesterday",
	"tea yesterday",
	"midnight today",
	"3pm today",
	"2am last friday",
	"19:00 yesterday",
}

func TestHasRelative(t *testing.T) {
	for i, test := range relativeTests {
		if !hasRelative(test) {
			t.Errorf("#%v: %v", i, test)
			continue
		}
	}
}

func TestParseRelative(t *testing.T) {
	for i, test := range relativeTests {
		parsed, err := ParseRelative(test)
		if err != nil {
			t.Errorf("#%v: %v", i, err)
			continue
		}
		t.Logf("#%v: %v (%v)", i, parsed, test)
	}
}

var TwelveHourClockTests = []string{
	"10am",
	"3pm",
	"1AM",
	"5PM",
}

func TestParse12HourClock(t *testing.T) {
	for i, test := range TwelveHourClockTests {
		parsed, err := parse12HourClock(test)
		if err != nil {
			t.Errorf("#%v: %v", i, err)
			continue
		}
		t.Logf("#%v: %v (%v)", i, parsed, test)
	}
}

func TestNumericDate(t *testing.T) {
	now := time.Now().In(time.Local)
	tests := map[string]time.Time{
		"2008-12-01":       time.Date(2008, 12, 1, 0, 0, 0, 0, time.Local),
		"06/05/2009":       time.Date(2009, 6, 5, 0, 0, 0, 0, time.Local),
		"06.05.2009":       time.Date(2009, 6, 5, 0, 0, 0, 0, time.Local),
		"06 05 2009":       time.Date(2009, 6, 5, 0, 0, 0, 0, time.Local),
		"10/30":            time.Date(now.Year(), 10, 30, 0, 0, 0, 0, time.Local),
		"01 02 2010 11:12": time.Date(2010, 1, 2, 11, 12, 0, 0, time.Local),
		"8 9 1999 1:22:33": time.Date(1999, 8, 9, 1, 22, 33, 0, time.Local),
	}

	for test, expected := range tests {
		parsed, err := parseNumeric(test)
		if err != nil {
			t.Errorf("%v: error parsing: %v", test, err)
			continue
		}
		if !parsed.Equal(expected) {
			t.Errorf("%v: wrongly parsed, got %v, %v expected", test, parsed, expected)
			continue
		}
		t.Logf("%v: %v (%v)", test, parsed, expected)
	}
}
