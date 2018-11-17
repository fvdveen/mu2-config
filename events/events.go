package events

import (
	"reflect"
	"strings"

	"github.com/fvdveen/mu2-config"
)

const (
	// Change is the EventType for a change
	Change EventType = iota
	// Add is the EventType for a addition
	Add
	// Remove is the EventType for a removal
	Remove
	// Slice is the EventType for a slice or array event
	Slice
	// Map is the EventType for a map event
	Map
)

// EventType is the type of change that happened in a event
type EventType uint8

// Event represents a change in the config
type Event struct {
	// Type shows what happened
	Type EventType
	// The config key that got changes e.g. discord.token
	Key string

	Change interface{}
	Slice  SliceEvent
	Map    MapEvent
}

// SliceEvent represents a change in a slice
type SliceEvent struct {
	Type      EventType
	Additions []interface{}
	Removals  []interface{}
}

// MapEvent represents a change in a map
type MapEvent struct {
	Type      EventType
	Additions map[interface{}]interface{}
	Changes   map[interface{}]interface{}
	Removals  map[interface{}]interface{}
}

// Watch puts all changes between the configs given by in into Events
func Watch(in <-chan *config.Config) <-chan *Event {
	ch := make(chan *Event)

	go func(in <-chan *config.Config, ch chan<- *Event) {
		last := &config.Config{}
		for conf := range in {
			changes(reflect.ValueOf(last), reflect.ValueOf(conf), ch)
			last = conf
		}

		close(ch)
	}(in, ch)

	return ch
}

func changes(a, b reflect.Value, ch chan<- *Event, keys ...string) {
	if reflect.DeepEqual(a.Interface(), b.Interface()) {
		return
	}

	if !isValid(a, b) {
		return
	}

	switch a.Type().Kind() {
	case reflect.Struct:
		structChanges(a, b, ch, keys...)
	case reflect.Ptr:
		changes(a.Elem(), b.Elem(), ch, keys...)
	case reflect.Slice, reflect.Array:
		sliceChanges(a, b, ch, keys...)
	case reflect.Map:
		mapChanges(a, b, ch, keys...)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Bool, reflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128, reflect.String:
		var e *Event

		e = &Event{
			Type:   Change,
			Key:    strings.Join(keys, Seperator),
			Change: b.Interface(),
		}

		ch <- e
	}
}

func structChanges(a, b reflect.Value, ch chan<- *Event, keys ...string) {
	if reflect.DeepEqual(a.Interface(), b.Interface()) {
		return
	}

	if !a.IsValid() || !b.IsValid() {
		return
	}

	if a.Type().Kind() != reflect.Struct {
		return
	}

	if a.Type() != b.Type() {
		return
	}

	for i := 0; i < a.NumField(); i++ {
		changes(a.Field(i), b.Field(i), ch, append(keys, strings.ToLower(a.Type().Field(i).Name))...)
	}
}

func sliceChanges(a, b reflect.Value, ch chan<- *Event, keys ...string) {
	if reflect.DeepEqual(a.Interface(), b.Interface()) {
		return
	}

	if !a.IsValid() || !b.IsValid() {
		return
	}

	if a.Type().Kind() != reflect.Slice && a.Type().Kind() != reflect.Array {
		return
	}

	if a.Type() != b.Type() {
		return
	}

	if a.Type().Kind() == reflect.Slice {
		if a.IsNil() || b.IsNil() {
			return
		}
	}

	if !a.CanInterface() || !b.CanInterface() {
		return
	}

	aVals := make(map[interface{}]bool)
	bVals := make(map[interface{}]bool)

	for i := 0; i < a.Len(); i++ {
		aVals[a.Index(i).Interface()] = true
	}

	for i := 0; i < b.Len(); i++ {
		bVals[b.Index(i).Interface()] = true
	}

	adds, rems := sliceAddsRems(aVals, bVals)
	if len(adds) != 0 && len(rems) != 0 {
		ch <- &Event{
			Type: Slice,
			Key:  strings.Join(keys, Seperator),
			Slice: SliceEvent{
				Type:      Change,
				Additions: adds,
				Removals:  rems,
			},
		}
	} else if len(adds) != 0 {
		ch <- &Event{
			Type: Slice,
			Key:  strings.Join(keys, Seperator),
			Slice: SliceEvent{
				Type:      Add,
				Additions: adds,
			},
		}
	} else if len(rems) != 0 {
		ch <- &Event{
			Type: Slice,
			Key:  strings.Join(keys, Seperator),
			Slice: SliceEvent{
				Type:     Remove,
				Removals: rems,
			},
		}
	}
}

func mapChanges(a, b reflect.Value, ch chan<- *Event, keys ...string) {
	if reflect.DeepEqual(a.Interface(), b.Interface()) {
		return
	}

	if !a.IsValid() || !b.IsValid() {
		return
	}

	if a.Type().Kind() != reflect.Map {
		return
	}

	if a.Type() != b.Type() {
		return
	}

	if a.IsNil() || b.IsNil() {
		return
	}

	if !a.CanInterface() || !b.CanInterface() {
		return
	}

	aVals := make(map[interface{}]interface{})
	bVals := make(map[interface{}]interface{})

	for _, x := range a.MapKeys() {
		y := a.MapIndex(x)
		if !x.IsValid() || !y.IsValid() {
			continue
		}
		aVals[x.Interface()] = y.Interface()
	}

	for _, x := range b.MapKeys() {
		y := b.MapIndex(x)
		if !x.IsValid() || !y.IsValid() {
			continue
		}
		bVals[x.Interface()] = y.Interface()
	}

	add, change, rem := mapAddChangeRems(aVals, bVals)
	if len(change) != 0 || len(add) != 0 && len(rem) != 0 {
		ch <- &Event{
			Type: Map,
			Key:  strings.Join(keys, Seperator),
			Map: MapEvent{
				Type:      Change,
				Additions: add,
				Changes:   change,
				Removals:  rem,
			},
		}
	} else if len(add) != 0 {
		ch <- &Event{
			Type: Map,
			Key:  strings.Join(keys, Seperator),
			Map: MapEvent{
				Type:      Add,
				Additions: add,
			},
		}
	} else if len(rem) != 0 {
		ch <- &Event{
			Type: Map,
			Key:  strings.Join(keys, Seperator),
			Map: MapEvent{
				Type:     Remove,
				Removals: rem,
			},
		}
	}
}

func sliceAddsRems(a, b map[interface{}]bool) ([]interface{}, []interface{}) {
	removals := []interface{}{}
	additions := []interface{}{}

	for x := range a {
		found := false
		for y := range b {
			if x == y {
				found = true
				break
			}
		}

		if !found {
			removals = append(removals, x)
		}
	}

	for x := range b {
		found := false
		for y := range a {
			if x == y {
				found = true
				break
			}
		}

		if !found {
			additions = append(additions, x)
		}
	}

	return additions, removals
}

func mapAddChangeRems(a, b map[interface{}]interface{}) (
	map[interface{}]interface{},
	map[interface{}]interface{},
	map[interface{}]interface{},
) {
	add := make(map[interface{}]interface{})
	change := make(map[interface{}]interface{})
	rem := make(map[interface{}]interface{})

	for k, v := range a {
		x, ok := b[k]
		if !ok {
			rem[k] = v
			continue
		}

		if reflect.DeepEqual(x, v) {
			continue
		}

		change[k] = x
	}

	for k, v := range b {
		if _, ok := a[k]; ok {
			continue
		}

		add[k] = v
	}

	return add, change, rem
}

func isValid(a, b reflect.Value) bool {
	if !a.IsValid() || !b.IsValid() {
		return false
	}

	if a.Type() != b.Type() {
		return false
	}

	switch a.Type().Kind() {
	case reflect.Func, reflect.Chan, reflect.Interface, reflect.Ptr, reflect.Slice, reflect.Map:
		if a.IsNil() || b.IsNil() {
			return false
		}
	}

	if !a.CanInterface() || !b.CanInterface() {
		return false
	}

	return true
}
