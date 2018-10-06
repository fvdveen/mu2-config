package events

import (
	"strings"
)

// Bot takes a chan of events and splits it into a chan of bot events and all other events
func Bot(ch <-chan *Event) (<-chan *Event, <-chan *Event) {
	bot := make(chan *Event)
	rest := make(chan *Event)
	go func(ch <-chan *Event, bot, rest chan<- *Event) {
		for evnt := range ch {
			d := strings.Split(evnt.Key, ".")[0]
			switch d {
			case "bot":
				go func(bot chan<- *Event, evnt *Event) {
					bot <- evnt
				}(bot, evnt)
			default:
				go func(rest chan<- *Event, evnt *Event) {
					rest <- evnt
				}(rest, evnt)
			}
		}

		close(bot)
		close(rest)
	}(ch, bot, rest)

	return bot, rest
}

// Log takes a chan of events and splits it into a chan of log events and all other events
func Log(ch <-chan *Event) (<-chan *Event, <-chan *Event) {
	log := make(chan *Event)
	rest := make(chan *Event)
	go func(ch <-chan *Event, log, rest chan<- *Event) {
		for evnt := range ch {
			d := strings.Split(evnt.Key, ".")[0]
			switch d {
			case "log":
				go func(log chan<- *Event, evnt *Event) {
					log <- evnt
				}(log, evnt)
			default:
				go func(rest chan<- *Event, evnt *Event) {
					rest <- evnt
				}(rest, evnt)
			}
		}

		close(log)
		close(rest)
	}(ch, log, rest)

	return log, rest
}

// Database takes a chan of events and splits it into a chan of database events and all other events
func Database(ch <-chan *Event) (<-chan *Event, <-chan *Event) {
	db := make(chan *Event)
	rest := make(chan *Event)
	go func(ch <-chan *Event, db, rest chan<- *Event) {
		for evnt := range ch {
			d := strings.Split(evnt.Key, ".")[0]
			switch d {
			case "database":
				go func(db chan<- *Event, evnt *Event) {
					db <- evnt
				}(db, evnt)
			default:
				go func(rest chan<- *Event, evnt *Event) {
					rest <- evnt
				}(rest, evnt)
			}
		}

		close(db)
		close(rest)
	}(ch, db, rest)

	return db, rest
}

// Youtube takes a chan of events and splits it into a chan of youtube events and all other events
func Youtube(ch <-chan *Event) (<-chan *Event, <-chan *Event) {
	yt := make(chan *Event)
	rest := make(chan *Event)
	go func(ch <-chan *Event, yt, rest chan<- *Event) {
		for evnt := range ch {
			d := strings.Split(evnt.Key, ".")[0]
			switch d {
			case "youtube":
				go func(yt chan<- *Event, evnt *Event) {
					yt <- evnt
				}(yt, evnt)
			default:
				go func(rest chan<- *Event, evnt *Event) {
					rest <- evnt
				}(rest, evnt)
			}
		}

		close(yt)
		close(rest)
	}(ch, yt, rest)

	return yt, rest
}

// Services takes a chan of events and splits it into a chan of service events and all other events
func Services(ch <-chan *Event) (<-chan *Event, <-chan *Event) {
	s := make(chan *Event)
	rest := make(chan *Event)
	go func(ch <-chan *Event, s, rest chan<- *Event) {
		for evnt := range ch {
			d := strings.Split(evnt.Key, ".")[0]
			switch d {
			case "services":
				go func(s chan<- *Event, evnt *Event) {
					s <- evnt
				}(s, evnt)
			default:
				go func(rest chan<- *Event, evnt *Event) {
					rest <- evnt
				}(rest, evnt)
			}
		}

		close(s)
		close(rest)
	}(ch, s, rest)

	return s, rest
}

// SearchService takes a chan of events and splits it into a chan of search service events and all other events
func SearchService(ch <-chan *Event) (<-chan *Event, <-chan *Event) {
	ss := make(chan *Event)
	rest := make(chan *Event)
	go func(ch <-chan *Event, ss, rest chan<- *Event) {
		for evnt := range ch {
			d := strings.Split(evnt.Key, ".")
			if len(d) < 2 {
				go func(rest chan<- *Event, evnt *Event) {
					rest <- evnt
				}(rest, evnt)
				continue
			}

			if d[0] != "services" {
				go func(rest chan<- *Event, evnt *Event) {
					rest <- evnt
				}(rest, evnt)
				continue
			}

			switch d[1] {
			case "search":
				go func(ss chan<- *Event, evnt *Event) {
					ss <- evnt
				}(ss, evnt)
			default:
				go func(rest chan<- *Event, evnt *Event) {
					rest <- evnt
				}(rest, evnt)
			}
		}

		close(ss)
		close(rest)
	}(ch, ss, rest)

	return ss, rest
}

// Null clears the channel given to it
// It is required because if it is not used the provider will deadlock
func Null(ch <-chan *Event) {
	go func(ch <-chan *Event) {
		for evnt := range ch {
			_ = evnt
		}
	}(ch)
}
