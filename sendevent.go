package eslgo

import (
	"bytes"
	"fmt"
	"strings"
	"time"
)

// EVT is the container used by SendEvent to store event sent to FreeSWITCH.
// It's support sendevent command only
//
// See https://freeswitch.org/confluence/display/FREESWITCH/mod_event_socket#mod_event_socket-3.9sendmsg for details.
type EVT map[string]string

// SendEvent sends events to FreeSWITCH and returns a response Event.
//
// Example:
//
//	SendEvent("CUSTOM", EVT{
//	"profile":        "internal",
//	"content-length": string(contentLen),
//	"content-type":   contentType,
//	"host":           "192.168.160.10",
//	}, content)
//
// eventName is the name of the sendevent, such as CUSTOM, SEND_MESSAGE, NOTIFY
// Keys with empty values are ignored; appData is optional.
// If appData is set, "content-length" and "content-type" headers are expected (lower case!).
//
// See https://freeswitch.org/confluence/display/FREESWITCH/mod_event_socket#mod_event_socket-3.9sendmsg for details.
func (h *Connection) SendEvent(eventName string, evt EVT, appData string) (*Event, error) {
	b := bytes.NewBufferString("sendevent")

	// append event name
	if eventName == "" {
		return nil, errEmptySendEvent
	} else {
		// Make sure there's no \r or \n in the event name.
		if strings.IndexAny(eventName, "\r\n") > 0 {
			return nil, errInvalidSendEvent
		}
		b.WriteString(" " + eventName)
	}

	// append header
	b.WriteString("\n")
	for k, v := range evt {
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

	// append content
	b.WriteString("\n")
	if evt["content-length"] != "" && appData != "" {
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

//// SendEventCustom execute sendevent with CUSTOM event.
////
//// --------------------------------------------------
////	sendevent CUSTOM
////	Event-Subclass: foo::bar
////	content-type: text/plain
////	content-length: 2
////
////	OK
//// --------------------------------------------------
//func (h *Connection) SendEventCustom(subclass, contentType, content string) (*Event, error) {
//	// calculate content length
//	contentLen := strings.NewReader(content).Len()
//
//	return h.SendEvent("CUSTOM", EVT{
//		"Event-Subclass": subclass,
//		"profile":        "internal",
//		"content-length": fmt.Sprintf("%d", contentLen),
//		"content-type":   contentType,
//		"host":           "192.168.160.10",
//	}, content)
//}

//// SendEventMessage execute sendevent with SEND_MESSAGE event.
////
//// --------------------------------------------------
////	sendevent SEND_MESSAGE
////	profile: internal
////	content-type: text/plain
////	content-length: 2
////	user: 1005
////	host: 99.157.44.194
////
////	OK
//// --------------------------------------------------
//func (h *Connection) SendEventMessage(user, contentType, content string) (*Event, error) {
//	// calculate content length
//	contentLen := strings.NewReader(content).Len()
//
//	return h.SendEvent("SEND_MESSAGE", EVT{
//		"profile":        "internal",
//		"content-type":   contentType,
//		"content-length": fmt.Sprintf("%d", contentLen),
//		"user":           user,
//		"host":           "192.168.160.10",
//	}, content)
//}

// SendEventNotify execute sendevent with NOTIFY event.
//
// --------------------------------------------------
//	sendevent NOTIFY
//	profile: internal
//	event-string: check-sync
//	content-type: application/simple-message-summary
//	content-length: 2
//	user: 1005
//	host: 192.168.10.4
//
//	OK
// --------------------------------------------------
func (h *Connection) SendEventNotify(user, contentType, content string) (*Event, error) {
	// calculate content length
	contentLen := strings.NewReader(content).Len()

	return h.SendEvent("NOTIFY", EVT{
		"profile":        "internal",
		"event-string":   "check-sync",
		"content-type":   contentType,
		"content-length": fmt.Sprintf("%d", contentLen),
		"user":           user,
		"host":           "192.168.160.10",
	}, content)
}
