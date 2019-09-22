package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

var loc *time.Location

func main() {
	if len(os.Args) < 2 {
		printArgumentErrorMessage()
		os.Exit(1)
	}

	unixSeconds, err := strconv.ParseInt(os.Args[1], 10, 64)
	if err != nil {
		printArgumentErrorMessage()
		os.Exit(1)
	}

	loc, err = time.LoadLocation("UTC")
	if err != nil {
		fmt.Printf("Unable to load timezone (UTC): %s", err.Error())
		os.Exit(1)
	}

	then := time.Unix(unixSeconds, 0).In(loc)
	c := realClock{}
	fmt.Printf("%s (%s)\n", then.Format(time.RFC3339), differenceFromNow(then, &c))
}

// interface to fetch 'now'. Useful for testing.
type clock interface {
	Now() time.Time
}

// realClock implements 'clock' interface as a proxy to time.Now()
type realClock struct{}

// Now returns actual time using time.Now()
func (r *realClock) Now() time.Time {
	return time.Now().In(loc)
}

// differenceFromNow returns a human readable difference text of the form:
// "30 days, 23 hours, 52 minutes, 14 seconds ago"
func differenceFromNow(then time.Time, c clock) string {
	now := c.Now()
	diff := now.Sub(then)

	// return "now" if difference is less than a second(1E9nsec)
	if diff > -1E9 && diff < 1E9 {
		return "now"
	}

	past := diff > 0
	var start, end []int
	if past {
		start, end = segments(then), segments(now)
	} else {
		start, end = segments(now), segments(then)
	}

	// calculate difference by each segment
	// we don't use the 'diff' calculated above, since it only holds nanoseconds(int64)
	var seconds, minutes, hours, days, months, years int
	seconds = end[5] - start[5]
	if seconds < 0 {
		seconds = seconds + 60
		end[4] = end[4] - 1
	}

	minutes = end[4] - start[4]
	if minutes < 0 {
		minutes = minutes + 60
		end[3] = end[3] - 1
	}

	hours = end[3] - start[3]
	if hours < 0 {
		hours = hours + 24
		end[2] = end[2] - 1
	}

	days = end[2] - start[2]
	if days < 0 {
		// find the length of the month to carry (https://yourbasic.org/golang/last-day-month-date/)
		monthDays := time.Date(end[0], time.Month(end[1]), 0, 0, 0, 0, 0, time.UTC).Day()
		days = days + monthDays
	}

	months = end[1] - start[1]
	if months < 0 {
		months = months + 12
		end[0] = end[0] - 1
	}

	years = end[0] - start[0]

	// format the output
	labels := []string{"year", "month", "day", "hour", "minute", "second"}
	values := []int{years, months, days, hours, minutes, seconds}

	diffStrings := []string{}
	for i := 0; i < len(labels); i++ {
		if values[i] > 0 {
			diffString := fmt.Sprintf("%d %s", values[i], labels[i])
			if values[i] > 1 {
				diffString += "s"
			}
			diffStrings = append(diffStrings, diffString)
		}
	}

	readout := strings.Join(diffStrings, ", ")
	if past {
		readout += " ago"
	} else {
		readout += " in the future"
	}
	return readout
}

// segments take a time and returns []int{years, months, days, hours, minutes, seconds}
func segments(t time.Time) []int {
	return []int{t.Year(), int(t.Month()), t.Day(), t.Hour(), t.Minute(), t.Second()}
}

func printArgumentErrorMessage() {
	fmt.Fprintln(os.Stderr, "Usage error. Expecting a Unix Timestamp as an argument.")
	printUsage()
}

func printUsage() {
	fmt.Println("Summary:")
	fmt.Println("\twhen is a cmdline utility which accepts a unix timestamps and printlns RFC3339(UTC) datetime sting and a human readable relative difference from now.")
	fmt.Println("")
	fmt.Println("Usage:")
	fmt.Println("\t$> when 1569054942")
	fmt.Println("\t2019-09-21T01:35:42-07:00 (15 hours, 3 minutes, 56 seconds ago)")
	fmt.Println("")
	fmt.Println("\t$> when 15")
	fmt.Println("\t1970-01-01T00:00:15Z (49 years, 8 months, 21 days, 51 minutes, 33 seconds ago)")
	fmt.Println("")
	fmt.Println("\t$> when 2000000000")
	fmt.Println("\t2033-05-18T03:33:20Z (13 years, 8 months, 26 days, 2 hours, 40 minutes, 1 second in the future)")
}
