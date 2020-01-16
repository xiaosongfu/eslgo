## eslgo 

The ESL client of golang for FreeSWTICH. (FreeSWTICH ESL 的 Golang 客户端)

## source file

* esl.go : core file,handle connect,read and write data.
* apicommand.go : send api commands to FreeSWITCH
* sendevent.go : send events to FreeSWITCH
* sendmsg.go : send messages to FreeSWITCH

## reference

* [mod_event_socket](https://freeswitch.org/confluence/display/FREESWITCH/mod_event_socket)
* [mod_event_socket api command](https://freeswitch.org/confluence/display/FREESWITCH/mod_event_socket#mod_event_socket-3.CommandDocumentation)
* [mod_event_socket sendmsg](https://freeswitch.org/confluence/display/FREESWITCH/mod_event_socket#mod_event_socket-3.9sendmsg)
* [mod_event_socket sendevent](https://freeswitch.org/confluence/display/FREESWITCH/mod_event_socket#mod_event_socket-3.8sendevent)

## the same project

* [https://github.com/0x19/goesl](https://github.com/0x19/goesl)
* [https://github.com/fiorix/go-eventsocket](https://github.com/fiorix/go-eventsocket)
* [https://github.com/cgrates/fsock](https://github.com/cgrates/fsock)
