module github.com/zakkbob/greenlight

go 1.25.1

require (
	github.com/golang-migrate/migrate/v4 v4.19.1
	github.com/julienschmidt/httprouter v1.3.0
	github.com/lib/pq v1.11.2
	github.com/wneessen/go-mail v0.7.2
	golang.org/x/crypto v0.48.0
	golang.org/x/time v0.15.0
)

require (
	github.com/BurntSushi/toml v1.6.0 // indirect
	golang.org/x/exp/typeparams v0.0.0-20231108232855-2478ac86f678 // indirect
	golang.org/x/mod v0.32.0 // indirect
	golang.org/x/sync v0.19.0 // indirect
	golang.org/x/text v0.34.0 // indirect
	golang.org/x/tools v0.41.0 // indirect
	honnef.co/go/tools v0.7.0 // indirect
)

tool honnef.co/go/tools/cmd/staticcheck
