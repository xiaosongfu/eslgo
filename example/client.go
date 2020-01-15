// Copyright 2013 Alexandre Fiori
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

// Event Socket client that connects to FreeSWITCH to originate a new call.
package main

import (
	"fmt"
	"log"

	"github.com/xiaosongfu/eslgo"
)

const dest = "user/1011"
const dialplan = "&socket(192.168.150.235:9090 async)"

func main() {
	c, err := eslgo.Dial("192.168.160.46:8021", "ClueCon")
	if err != nil {
		log.Fatal(err)
	}
	c.Command("events json ALL")
	c.Command(fmt.Sprintf("bgapi originate %s %s", dest, dialplan))
	for {
		ev, err := c.ReadEvent()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("\nNew event")
		ev.PrettyPrint()
		if ev.Get("Answer-State") == "hangup" {
			break
		}
	}
	c.Close()
}
