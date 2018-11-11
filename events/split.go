package events

import (
	"strings"
)

// Seperator is the token on which event keys are seperated
const Seperator = "."

// Split takes a chan of events and splits it based on the keys provided
// it returns a chan of the event belonging to the keys events and all other events
func Split(ch <-chan *Event, keys ...string) (<-chan *Event, <-chan *Event) {
	split := make(chan *Event)
	rest := make(chan *Event)
	go func(ch <-chan *Event, split, rest chan<- *Event, keys ...string) {
	splitLoop:
		for evnt := range ch {
			parts := strings.Split(evnt.Key, Seperator)
			for i := range keys {
				if i >= len(parts) {
					rest <- evnt
					continue splitLoop
				} else if keys[i] != parts[i] {
					rest <- evnt
					continue splitLoop
				}
			}
			split <- evnt
		}

		close(split)
		close(rest)
	}(ch, split, rest, keys...)

	return split, rest
}

// Null clears the channel given to it
// It is required because if it is not used goroutines will leak
func Null(ch <-chan *Event) {
	go func(ch <-chan *Event) {
		for evnt := range ch {
			_ = evnt
		}
	}(ch)
}
