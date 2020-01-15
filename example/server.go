// Copyright 2013 Alexandre Fiori
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

// Server that accepts connections from FreeSWITCH and controls incoming calls.
package main

import (
	"fmt"
	"log"

	"github.com/xiaosongfu/eslgo"
)

const audioFile = "/tmp/freeswitch/welcome.wav"

func main() {
	eslgo.ListenAndServe(":9090", handler)
}

func handler(c *eslgo.Connection) {
	fmt.Println("new client:", c.RemoteAddr())
	c.Command("connect")
	c.Command("linger 10")
	c.Command("myevents")
	c.Execute("answer", "", false)
	ev, err := c.Execute("playback", audioFile, true)
	if err != nil {
		log.Fatal(err)
	}
	ev.PrettyPrint()
	for {
		ev, err = c.ReadEvent()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("\nNew event")
		ev.PrettyPrint()
		if ev.Get("Application") == "playback" {
			if ev.Get("Application-Response") == "FILE PLAYED" {
				c.Command("exit")
			}
		}
	}
}
