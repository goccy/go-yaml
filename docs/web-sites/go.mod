module demo

go 1.22.8

require (
	github.com/goccy/go-graphviz v0.2.9
	github.com/goccy/go-yaml v1.15.13
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

replace github.com/goccy/go-graphviz => ../../../go-graphviz

replace github.com/flopp/go-fontdir => ../../../go-findfont

replace github.com/flopp/go-findfont => ../../../go-findfont
