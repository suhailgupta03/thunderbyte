module github.com/suhailgupta03/thunderbyte/core

go 1.22.0

require (
	github.com/jmoiron/sqlx v1.3.5
	github.com/knadh/koanf/v2 v2.1.0
	github.com/labstack/echo/v4 v4.11.4
	//github.com/suhailgupta03/thunderbyte/common v0.0.0-00010101000000-000000000000
	//github.com/suhailgupta03/thunderbyte/database v0.0.0-00010101000000-000000000000
	github.com/suhailgupta03/thunderbyte/otp v0.0.1
	github.com/zerodha/logf v0.5.5
)

require (
	github.com/suhailgupta03/thunderbyte/common v0.0.0-20240306185410-3ebf5146195a
	github.com/suhailgupta03/thunderbyte/database v0.0.0-20240306185410-3ebf5146195a
)

require (
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/go-viper/mapstructure/v2 v2.0.0-alpha.1 // indirect
	github.com/knadh/goyesql/v2 v2.2.0 // indirect
	github.com/knadh/koanf/maps v0.1.1 // indirect
	github.com/labstack/gommon v0.4.2 // indirect
	github.com/lib/pq v1.10.9 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mitchellh/copystructure v1.2.0 // indirect
	github.com/mitchellh/reflectwalk v1.0.2 // indirect
	github.com/redis/go-redis/v9 v9.5.1 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasttemplate v1.2.2 // indirect
	golang.org/x/crypto v0.21.0 // indirect
	golang.org/x/net v0.22.0 // indirect
	golang.org/x/sys v0.18.0 // indirect
	golang.org/x/text v0.14.0 // indirect
)

//replace github.com/suhailgupta03/thunderbyte/common => ../common
//
//replace github.com/suhailgupta03/thunderbyte/database => ../database
