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
		}
	}
}

func TestParseAgo(t *testing.T) {
	for i, test := range agoTests {
		parsed, err := ParseAgo(test)
		if err != nil {
			t.Errorf("#%v: %v", i, err)
		}
		t.Logf("#%v: parsed: %v (%v)", i, parsed, test)
	}
}

var relativeTests = []string{
	"now",
	"today",
	"yesterday",
	"last friday",
}

func TestIsRelative(t *testing.T) {
	for i, test := range relativeTests {
		if !isRelative(test) {
			t.Errorf("#%v: %v", i, test)
		}
	}
}
