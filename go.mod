module github.com/teramono/engine-api

go 1.16

require (
	github.com/gin-gonic/gin v1.7.4
	github.com/teramono/utilities v0.0.0-20210919081101-b247dd3f53c0
	gorm.io/driver/sqlite v1.1.5 // indirect
)

replace github.com/teramono/utilities v0.0.0-20210919081101-b247dd3f53c0 => ../utilities
