module github.com/eduardoths/sandbox/consensus-simulator

go 1.19

replace (
	github.com/eduardoths/sandbox/go-utils/http => ../go-utils/http
	github.com/eduardoths/sandbox/go-utils/worker-pool => ../go-utils/worker-pool
)

require (
	github.com/eduardoths/sandbox/go-utils/http v0.0.0-20230313222343-01c5b0975626
	github.com/eduardoths/sandbox/go-utils/worker-pool v0.0.0-20230315155155-714cd518e5d6
	github.com/gofiber/fiber/v2 v2.42.0
	github.com/google/uuid v1.3.0
	gopkg.in/yaml.v3 v3.0.1
)

require (
	github.com/andybalholm/brotli v1.0.4 // indirect
	github.com/klauspost/compress v1.15.9 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.17 // indirect
	github.com/mattn/go-runewidth v0.0.14 // indirect
	github.com/philhofer/fwd v1.1.1 // indirect
	github.com/rivo/uniseg v0.2.0 // indirect
	github.com/savsgio/dictpool v0.0.0-20221023140959-7bf2e61cea94 // indirect
	github.com/savsgio/gotils v0.0.0-20220530130905-52f3993e8d6d // indirect
	github.com/tinylib/msgp v1.1.6 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasthttp v1.44.0 // indirect
	github.com/valyala/tcplisten v1.0.0 // indirect
	golang.org/x/sys v0.0.0-20220811171246-fbc7d0a398ab // indirect
)
