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
	"fmt"
	"strconv"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"
)

// Parse accepts contextful date format and returns absolute time.Time value.
func Parse(value string) (time.Time, error) {
	value = strings.TrimSpace(value)
	switch {
	case strings.HasSuffix(value, "ago"):
		return ParseAgo(value)
	case hasRelative(value):
		return ParseRelative(value)
	default:
		return ParseAbsolute(value)
	}
	return time.Now().In(time.Local), nil // TODO(ymotongpoo): replace actual time.
}

// splitTokens splits value with commas, periods and spaces.
// Currently, it only expects single byte character tokenizer.
func splitTokens(value string) []string {
	f := func(c rune) bool {
		return c == rune(' ') || c == rune(',') || c == rune('.')
	}
	return strings.FieldsFunc(value, f)
}

// hasRelative confirms if value contains relative datatime words, such as
// "now", "today", "last xxx", "noon", "pm", "am"  and so on.
func hasRelative(value string) bool {
	keywords := []string{"now", "today", "yesterday", "last"}
	for _, k := range keywords {
		if strings.Contains(value, k) {
			return true
		}
	}
	return false
}

// ParseAgo parse "xxxx ago" format and returns corresponding absolute datetime.
func ParseAgo(value string) (time.Time, error) {
	tokens := splitTokens(value)
	now := time.Now().In(time.Local)
	for i := 0; i < len(tokens); i++ {
		t := tokens[i]
		if t == "ago" {
			return now, nil
		}
		if i%2 == 0 {
			var err error
			n, err := strconv.Atoi(t)
			if err != nil {
				return now, fmt.Errorf("Format error: %v", t)
			}
			now, err = subDate(now, n, tokens[i+1])
			if err != nil {
				return now, err
			}
			i++
		}
	}
	return now, nil
}

// subDate subtracts n*unit duration from t and return the result.
// supportes units are "year", "month", "week", "day", "hour", "minute", "second", and those plurals.
func subDate(t time.Time, n int, unit string) (time.Time, error) {
	if strings.HasSuffix(unit, "s") {
		unit = string([]byte(unit)[:len(unit)-1])
	}
	switch unit {
	case "year":
		return t.AddDate(-1*n, 0, 0), nil
	case "month":
		return t.AddDate(0, -1*n, 0), nil
	case "week":
		return t.AddDate(0, 0, -7*n), nil
	case "day":
		return t.AddDate(0, 0, -1*n), nil
	case "hour":
		return t.Add(time.Duration(-1*n) * time.Hour), nil
	case "minute":
		return t.Add(time.Duration(-1*n) * time.Minute), nil
	case "second":
		return t.Add(time.Duration(-1*n) * time.Second), nil
	default:
		return t, fmt.Errorf("Unsupported time unit: %v", unit)
	}
}

// ParseRelative returns absolute datetime corresponding to relative date expressed in value.
func ParseRelative(value string) (time.Time, error) {
	tokens := splitTokens(value)
	_ = tokens
	var t time.Time
	t = time.Now() // TODO(ymotongpoo): implement me
	return t, nil
}

func ParseTimeWord(word string) (time.Time, error) {
	now := time.Now().In(time.Local)
	switch word {
	case "now", "today":
		return now, nil
	case "noon":
		return time.Date(now.Year(), now.Month(), now.Day(), 12, 0, 0, 0, time.Local), nil
	case "tea":
		return time.Date(now.Year(), now.Month(), now.Day(), 15, 0, 0, 0, time.Local), nil
	case "midnight":
		return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local), nil
	}
	return now, fmt.Errorf("Unsupported time word: %v", word)
}

func parse12HourClock(word string) (time.Time, error) {
	lower := strings.ToLower(word)
	now := time.Now().In(time.Local)

	start := 0
	hour := 0
	var err error
	for width := 0; start < len(lower); start += width {
		var r rune
		r, width = utf8.DecodeRuneInString(lower[start:])
		if !unicode.IsNumber(r) {
			hour, err = strconv.Atoi(lower[:start])
			if err != nil || hour > 12 || hour < 0 {
				return time.Now(), fmt.Errorf("Wrong hour: %v", word)
			}
			if string(lower[start:]) == "am" {
				break
			}
			if string(lower[start:]) == "pm" {
				hour += 12
				break
			}
			return time.Now(), fmt.Errorf("Unsupported 12 hour clock notation: %v", word)
		}
	}
	return time.Date(now.Year(), now.Month(), now.Day(), hour, 0, 0, 0, time.Local), nil
}

func ParseAbsolute(value string) (time.Time, error) {
	return time.Now(), nil // TODO(ymotongpoo): implement me
}
