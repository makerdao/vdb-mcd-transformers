package backfill

import (
	"errors"
	"fmt"
)

var (
	ErrDuplicateEvent   = errors.New("must pass distinct events when passing more than one")
	ErrTooFewEvents     = errors.New("must specify at least one event to backfill")
	ErrTooManyEvents    = fmt.Errorf("cannot backfill more than %d events", MaxEvents)
	ErrUnsupportedEvent = fmt.Errorf("invalid event detected: only %s, %s, %s events supported", ForkEvent, FrobEvent, GrabEvent)
	ForkEvent           = "fork"
	FrobEvent           = "frob"
	GrabEvent           = "grab"
	MaxEvents           = 3
)

func ValidateArgs(eventsToBackFill []string) error {
	lenEvents := len(eventsToBackFill)
	if lenEvents < 1 {
		return ErrTooFewEvents
	}
	if lenEvents > MaxEvents {
		return ErrTooManyEvents
	}
	eventOne := eventsToBackFill[0]
	if lenEvents == 1 {
		return validateOneArg(eventOne)
	}
	eventTwo := eventsToBackFill[1]
	if lenEvents == 2 {
		return validateTwoArgs(eventOne, eventTwo)
	}
	eventThree := eventsToBackFill[2]
	if lenEvents == 3 {
		return validateThreeArgs(eventOne, eventTwo, eventThree)
	}
	return fmt.Errorf("backfill argument validation not supported for %d events; please file an issue", lenEvents)
}

func validateOneArg(event string) error {
	if event != ForkEvent && event != FrobEvent && event != GrabEvent {
		return ErrUnsupportedEvent
	}
	return nil
}

func validateTwoArgs(eventOne, eventTwo string) error {
	if !argIsSupportedEvent(eventOne) || !argIsSupportedEvent(eventTwo) {
		return ErrUnsupportedEvent
	}
	if eventOne == eventTwo {
		return ErrDuplicateEvent
	}
	return nil
}

func validateThreeArgs(eventOne, eventTwo, eventThree string) error {
	if !argIsSupportedEvent(eventOne) || !argIsSupportedEvent(eventTwo) || !argIsSupportedEvent(eventThree) {
		return ErrUnsupportedEvent
	}
	if eventOne == eventTwo || eventOne == eventThree || eventTwo == eventThree {
		return ErrDuplicateEvent
	}
	return nil
}

func argIsSupportedEvent(arg string) bool {
	return arg == ForkEvent || arg == FrobEvent || arg == GrabEvent
}
