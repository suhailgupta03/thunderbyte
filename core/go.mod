module github.com/suhailgupta03/thunderbyte/core

go 1.22.0

require (
	github.com/labstack/echo/v4 v4.11.4
	github.com/suhailgupta03/thunderbyte/common v0.0.0-00010101000000-000000000000
	github.com/suhailgupta03/thunderbyte/database v0.0.0-00010101000000-000000000000
)

require (
	github.com/jmoiron/sqlx v1.3.5 // indirect
	github.com/knadh/goyesql/v2 v2.2.0 // indirect
	github.com/labstack/gommon v0.4.2 // indirect
	github.com/lib/pq v1.10.9 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasttemplate v1.2.2 // indirect
	golang.org/x/crypto v0.17.0 // indirect
	golang.org/x/net v0.19.0 // indirect
	golang.org/x/sys v0.15.0 // indirect
	golang.org/x/text v0.14.0 // indirect
)

replace github.com/suhailgupta03/thunderbyte/common => ../common

replace github.com/suhailgupta03/thunderbyte/database => ../database
