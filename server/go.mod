module server

go 1.25.1

require (
	github.com/joy-stream/gamepad v0.0.0
	udp v0.0.0
)

replace github.com/joy-stream/gamepad => ../gamepad
replace udp => ../udp
