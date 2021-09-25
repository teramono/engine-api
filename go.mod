module github.com/teramono/engine-api

go 1.16

replace github.com/teramono/utilities v0.0.0-20210919081101-b247dd3f53c0 => ../utilities

replace github.com/teramono/engine-backend v0.0.0-20210921231255-d6d2e37aec06 => ../engine-backend

replace github.com/teramono/engine-fs v0.0.0-20210924140556-e15c34e7dbcd => ../engine-fs

replace github.com/teramono/engine-db v0.0.0-20210924140608-36967c4678af => ../engine-db

require (
	github.com/gin-gonic/gin v1.7.4
	github.com/teramono/engine-backend v0.0.0-20210921231255-d6d2e37aec06
	github.com/teramono/utilities v0.0.0-20210919081101-b247dd3f53c0
)
