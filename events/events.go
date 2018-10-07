package events

import (
	"reflect"

	"github.com/fvdveen/mu2-config"
)

const (
	// Change is the EventType for a change
	Change EventType = iota
	// Add is the EventType for a addition
	Add
	// Remove is the EventType for a removal
	Remove
)

// EventType is the type of change that happened in a event
type EventType uint8

// Event represents a change in the config
type Event struct {
	// EventType shows what happened
	EventType EventType
	// The config key that got changes e.g. discord.token
	Key string

	Change    string
	Additions []string
	Removals  []string
	Database  config.Database
	Log       config.Log
}

// Watch puts all changes between the configs given by in into Events
func Watch(in <-chan *config.Config) <-chan *Event {
	ch := make(chan *Event)

	go func(in <-chan *config.Config, ch chan<- *Event) {
		last := &config.Config{}
		for conf := range in {
			if !reflect.DeepEqual(conf.Bot, last.Bot) {
				botChanges(ch, conf, last)
			}
			if !reflect.DeepEqual(conf.Log, last.Log) {
				logChanges(ch, conf, last)
			}
			if !reflect.DeepEqual(conf.Database, last.Database) {
				ch <- &Event{
					EventType: Change,
					Key:       "database",
					Change:    conf.Database.Type,
					Database:  conf.Database,
				}
			}
			if !reflect.DeepEqual(conf.Youtube, last.Youtube) {
				ytChanges(ch, conf, last)
			}
			if !reflect.DeepEqual(conf.Services, last.Services) {
				serviceChanges(ch, conf, last)
			}

			last = conf
		}

		close(ch)
	}(in, ch)

	return ch
}

func logChanges(ch chan<- *Event, conf *config.Config, last *config.Config) {
	if conf.Log.Level != last.Log.Level {
		ch <- &Event{
			EventType: Change,
			Key:       "log.level",
			Change:    conf.Log.Level,
		}
	}
	if !reflect.DeepEqual(conf.Log.Discord, last.Log.Discord) {
		ch <- &Event{
			EventType: Change,
			Key:       "log.discord",
			Change:    "hook",
			Log:       conf.Log,
		}
	}
}

func ytChanges(ch chan<- *Event, conf *config.Config, last *config.Config) {
	if conf.Youtube.APIKey != last.Youtube.APIKey {
		ch <- &Event{
			EventType: Change,
			Key:       "youtube.apikey",
			Change:    conf.Youtube.APIKey,
		}
	}
}

func botChanges(ch chan<- *Event, conf *config.Config, last *config.Config) {
	if conf.Bot.Discord.Token != last.Bot.Discord.Token {
		ch <- &Event{
			EventType: Change,
			Key:       "bot.discord.token",
			Change:    conf.Bot.Discord.Token,
		}
	}
	if conf.Bot.Prefix != last.Bot.Prefix {
		ch <- &Event{
			EventType: Change,
			Key:       "bot.prefix",
			Change:    conf.Bot.Prefix,
		}
	}

	if !reflect.DeepEqual(conf.Bot.Commands, last.Bot.Commands) {
		a, r := changes(conf.Bot.Commands, last.Bot.Commands)
		if len(a) == 0 && len(r) == 0 {
		} else if len(r) == 0 {
			ch <- &Event{
				EventType: Add,
				Key:       "bot.commands",
				Additions: a,
			}
		} else if len(a) == 0 {
			ch <- &Event{
				EventType: Remove,
				Key:       "bot.commands",
				Removals:  r,
			}
		} else {
			ch <- &Event{
				EventType: Change,
				Key:       "bot.commands",
				Additions: a,
				Removals:  r,
			}
		}
	}
}

func serviceChanges(ch chan<- *Event, conf *config.Config, last *config.Config) {
	if !reflect.DeepEqual(conf.Services.Search, last.Services.Search) {
		searchChanges(ch, conf, last)
	}
	if !reflect.DeepEqual(conf.Services.Encode, last.Services.Encode) {
		encodeChanges(ch, conf, last)
	}
}

func searchChanges(ch chan<- *Event, conf *config.Config, last *config.Config) {
	if conf.Services.Search.Location != last.Services.Search.Location {
		ch <- &Event{
			EventType: Change,
			Key:       "services.search.location",
			Change:    conf.Services.Search.Location,
		}
	}
}

func encodeChanges(ch chan<- *Event, conf *config.Config, last *config.Config) {
	if conf.Services.Encode.Location != last.Services.Encode.Location {
		ch <- &Event{
			EventType: Change,
			Key:       "services.encode.location",
			Change:    conf.Services.Encode.Location,
		}
	}
}

func changes(new []string, old []string) (additions []string, removals []string) {
	oldComs := map[string]bool{}
	for _, com := range old {
		oldComs[com] = true
	}
	newComs := map[string]bool{}
	for _, com := range new {
		newComs[com] = true
	}

	for x := range oldComs {
		found := false
		for y := range newComs {
			if x == y {
				found = true
				break
			}
		}
		if found {
			continue
		}
		removals = append(removals, x)
	}

	for x := range newComs {
		double := false
		for y := range oldComs {
			if x == y {
				double = true
				break
			}
		}
		if double {
			continue
		}
		additions = append(additions, x)
	}

	return additions, removals
}
