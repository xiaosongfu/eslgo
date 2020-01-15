package eslgo

import (
	"bytes"
	"fmt"
	"strings"
	"time"
)

// MSG is the container used by SendMsg to store messages sent to FreeSWITCH.
// It's supposed to be populated with directives supported by the sendmsg
// command only, like "call-command: execute".
//
// See https://freeswitch.org/confluence/display/FREESWITCH/mod_event_socket#mod_event_socket-3.9sendmsg for details.
type MSG map[string]string

// SendMsg sends messages to FreeSWITCH and returns a response Event.
//
// Examples:
//
//	SendMsg(MSG{
//		"call-command": "hangup",
//		"hangup-cause": "we're done!",
//	}, "", "")
//
//	SendMsg(MSG{
//		"call-command":     "execute",
//		"execute-app-name": "playback",
//		"execute-app-arg":  "/tmp/test.wav",
//	}, "", "")
//
// Keys with empty values are ignored; uuid and appData are optional.
// If appData is set, a "content-length" header is expected (lower case!).
//
// See https://freeswitch.org/confluence/display/FREESWITCH/mod_event_socket#mod_event_socket-3.9sendmsg for details.
func (h *Connection) SendMsg(m MSG, uuid, appData string) (*Event, error) {
	b := bytes.NewBufferString("sendmsg")
	if uuid != "" {
		// Make sure there's no \r or \n in the UUID.
		if strings.IndexAny(uuid, "\r\n") > 0 {
			return nil, errInvalidCommand
		}
		b.WriteString(" " + uuid)
	}
	b.WriteString("\n")
	for k, v := range m {
		// Make sure there's no \r or \n in the key, and value.
		if strings.IndexAny(k, "\r\n") > 0 {
			return nil, errInvalidCommand
		}
		if v != "" {
			if strings.IndexAny(v, "\r\n") > 0 {
				return nil, errInvalidCommand
			}
			b.WriteString(fmt.Sprintf("%s: %s\n", k, v))
		}
	}
	b.WriteString("\n")
	if m["content-length"] != "" && appData != "" {
		b.WriteString(appData)
	}
	if _, err := b.WriteTo(h.conn); err != nil {
		return nil, err
	}
	var (
		ev  *Event
		err error
	)
	select {
	case err = <-h.err:
		return nil, err
	case ev = <-h.cmd:
		return ev, nil
	case <-time.After(timeoutPeriod):
		return nil, errTimeout
	}
}

// Execute is a shortcut to SendMsg with call-command: execute without UUID,
// suitable for use on outbound event socket connections (acting as server).
//
// Example:
//
//	Execute("playback", "/tmp/test.wav", false)
//
// See https://freeswitch.org/confluence/display/FREESWITCH/mod_event_socket#mod_event_socket-3.9sendmsg for details.
func (h *Connection) Execute(appName, appArg string, lock bool) (*Event, error) {
	var evlock string
	if lock {
		// Could be strconv.FormatBool(lock), but we don't want to
		// send event-lock when it's set to false.
		evlock = "true"
	}
	return h.SendMsg(MSG{
		"call-command":     "execute",
		"execute-app-name": appName,
		"execute-app-arg":  appArg,
		"event-lock":       evlock,
	}, "", "")
}

// ExecuteUUID is similar to Execute, but takes a UUID and no lock. Suitable
// for use on inbound event socket connections (acting as client).
func (h *Connection) ExecuteUUID(uuid, appName, appArg string) (*Event, error) {
	return h.SendMsg(MSG{
		"call-command":     "execute",
		"execute-app-name": appName,
		"execute-app-arg":  appArg,
	}, uuid, "")
}
