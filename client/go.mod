module client

go 1.25.1

require (
	github.com/joy-stream/gamepad v0.0.0
	keyboard                 v0.0.0
	udp                      v0.0.0
)

require (
	github.com/holoplot/go-evdev v0.0.0-20250804134636-ab1d56a1fe83 // indirect
	golang.org/x/sys v0.41.0 // indirect
	golang.org/x/term v0.40.0 // indirect
)

replace github.com/joy-stream/gamepad => ../gamepad
replace keyboard => ../keyboard
replace udp => ../udp
