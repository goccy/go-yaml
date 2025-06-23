module demo

go 1.22.8

require (
	github.com/goccy/go-graphviz v0.2.10-0.20250109095217-4ceff9e58e1a
	github.com/goccy/go-json v0.10.4
	github.com/goccy/go-yaml v1.16.0
)

require (
	github.com/disintegration/imaging v1.6.2 // indirect
	github.com/flopp/go-findfont v0.1.0 // indirect
	github.com/fogleman/gg v1.3.0 // indirect
	github.com/golang/freetype v0.0.0-20170609003504-e2365dfdc4a0 // indirect
	github.com/tetratelabs/wazero v1.8.1 // indirect
	golang.org/x/image v0.21.0 // indirect
	golang.org/x/text v0.19.0 // indirect
)

replace github.com/goccy/go-yaml => ../../

replace github.com/flopp/go-findfont => github.com/goccy/go-findfont v0.0.0-20250109093214-c2e12b298c75
