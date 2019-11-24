# goflux

![Code Analysis](https://github.com/deltacat/goflux/workflows/CA/badge.svg)
![Unit Test](https://github.com/deltacat/goflux/workflows/CI/badge.svg)
[![GoDoc](https://godoc.org/github.com/deltacat/goflux?status.svg)](https://godoc.org/github.com/deltacat/goflux)

an influx client for golang

## Usage

```go
package main
import (
    "github.com/deltacat/goflux"
)

//Setup setup storage
func Setup() error {
	cli, err := goflux.CreateClient(
		"http://locahost:8086",
		"username",
		"password",
		"database",
		"precision")
	if err != nil {
		return err
	}
    
	return nil
}
```