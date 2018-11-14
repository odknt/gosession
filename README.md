# gosession

[![license](http://img.shields.io/badge/license-MIT-red.svg?style=flat)](https://raw.githubusercontent.com/odknt/gosession/master/LICENSE) [![Coverage](http://gocover.io/_badge/github.com/odknt/gosession)](http://gocover.io/github.com/odknt/gosession) [![Build Status](https://api.travis-ci.com/odknt/gosession.svg?branch=master)](https://travis-ci.com/odknt/gosession)

A simple implementation for net/http.

```go
package main

import (
	"net/http"

	"github.com/odknt/gosession"
	_ "github.com/odknt/gosession/memory"
)

var manager *session.Manager

func helloOnce(w http.ResponseWriter, r *http.Request) {
	s, _ := manager.Start(w, r)
	check, ok := s.Get("suspicious person").(bool)
	if ok && check {
		w.Write([]byte("<h1>NEVER COME BACK HERE!!!!</h1>"))
		return
	}
	s.Set("suspicious person", true)
	w.Write([]byte("<p>Hello, what did you come here for?</p>"))
}

func main() {
	manager, _ = session.New("memory", "gosession", 86400)
	http.ListenAndServe(":8000", http.HandlerFunc(helloOnce))
}
```
