package eslgo

import (
	"fmt"
	"time"
)

// Command sends a single command to the server and returns a response Event.
//
// See https://freeswitch.org/confluence/display/FREESWITCH/mod_event_socket#mod_event_socket-3.CommandDocumentation for details.
func (h *Connection) Command(command string) (*Event, error) {
	// Sanity check to avoid breaking the parser
	//if strings.IndexAny(command, "\r\n") > 0 {
	//	return nil, errInvalidCommand
	//}
	fmt.Fprintf(h.conn, "%s\r\n\r\n", command)
	var (
		ev  *Event
		err error
	)
	select {
	case err = <-h.err:
		return nil, err
	case ev = <-h.cmd:
		return ev, nil
	case ev = <-h.api:
		return ev, nil
	case <-time.After(timeoutPeriod):
		return nil, errTimeout
	}
}
