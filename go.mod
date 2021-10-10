module github.com/teramono/engine-api

go 1.16

replace github.com/teramono/utilities v0.0.0-20210919081101-b247dd3f53c0 => ../utilities

require (
	github.com/gin-contrib/location v0.0.2
	github.com/gin-contrib/static v0.0.1
	github.com/gin-gonic/gin v1.7.4
	github.com/go-playground/validator/v10 v10.9.0
	github.com/teramono/utilities v0.0.0-20210919081101-b247dd3f53c0
	gorm.io/gorm v1.21.15
)
