package main

import (
	"fmt"
	"os"

	"github.com/cleanshaven/noaaweather"
	"github.com/godbus/dbus/v5"
	"github.com/godbus/dbus/v5/introspect"
)

const intro = `
<node>
    <interface name="com.github.cleanshaven.WeatherService">
         <method name="GetWeather">
             <arg direction="in" type="s"/>
             <arg direction="in" type="s"/>
             <arg direction="out" type="s"/>
             <arg direction="out" type="s"/>
         </method>
         <method name="GetIcon">
             <arg direction="in" type="s"/>
             <arg direction="out" type="s"/>
         </method>
    </interface>` + introspect.IntrospectDataString + `</node>`

func main() {
	conn, err := dbus.ConnectSessionBus()
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	server := serverObject{}

	conn.Export(server, noaaweather.ServerPath, noaaweather.ServerName)
	conn.Export(introspect.Introspectable(intro), "/com/github/cleanshaven/WeatherService",
		"org.freedesktop.DBus.Introspectable")

	reply, err := conn.RequestName(noaaweather.ServerName, dbus.NameFlagDoNotQueue)
	if err != nil {
		panic(err)
	}
	if reply != dbus.RequestNameReplyPrimaryOwner {
		fmt.Fprintln(os.Stderr, "name already taken")
		os.Exit(1)
	}
	fmt.Println("Listening on com.github.cleanshaven.WeatherService")
	select {} //block until done

}
