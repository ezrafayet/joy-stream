module client

go 1.25.1

require (
	github.com/eiannone/keyboard v0.0.0-20220611211555-0d226195f203
	github.com/joy-stream/protocol v0.0.0
)

require golang.org/x/sys v0.0.0-20220520151302-bc2c85ada10a // indirect

replace github.com/joy-stream/protocol => ../protocol
