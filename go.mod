module weatherserver

go 1.20

// replace github.com/cleanshaven/noaaweather => /home/bruce/development/go/noaaweather

require (
	github.com/cleanshaven/noaaweather v1.0.2
	github.com/godbus/dbus/v5 v5.1.0
)

require github.com/davecgh/go-spew v1.1.1
