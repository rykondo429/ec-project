module ec-sample/services/order-service

go 1.23

toolchain go1.24.1

replace ec-sample/shared => ../../shared

require (
	ec-sample/shared v0.0.0-00010101000000-000000000000
	github.com/google/uuid v1.5.0
	github.com/labstack/echo/v4 v4.11.4
	github.com/oapi-codegen/runtime v1.1.2
	go.uber.org/zap v1.26.0
	gorm.io/driver/mysql v1.5.4
	gorm.io/gorm v1.25.7-0.20240204074919-46816ad31dde
)

require (
	github.com/apapsch/go-jsonmerge/v2 v2.0.0 // indirect
	github.com/go-sql-driver/mysql v1.7.0 // indirect
	github.com/golang-jwt/jwt v3.2.2+incompatible // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/labstack/gommon v0.4.2 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasttemplate v1.2.2 // indirect
	go.uber.org/multierr v1.10.0 // indirect
	golang.org/x/crypto v0.17.0 // indirect
	golang.org/x/net v0.19.0 // indirect
	golang.org/x/sys v0.15.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	golang.org/x/time v0.5.0 // indirect
)
