package main

import (
	"context"
	"fmt"
	"log"

	"github.com/coreos/go-systemd/v22/dbus"
)

func main() {
	conn, err := dbus.NewWithContext(context.Background())
	if err != nil {
		log.Fatalln(err.Error())
	}
	if err := conn.Subscribe(); err != nil {
		log.Fatalln(err.Error())
	}
	updateCh := make(chan *dbus.PropertiesUpdate, 256)
	errorCh := make(chan error, 256)
	conn.SetPropertiesSubscriber(updateCh, errorCh)

	for {
		select {
		case change := <-updateCh:
			fmt.Printf("Service: %s, ActiveState: %s\n", change.UnitName, change.Changed["ActiveState"].String())
		case err = <-errorCh:
			log.Fatalln(err)
		}
	}
}
