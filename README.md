# goFiddleFastly

[![Build Status](https://travis-ci.org/GannettDigital/goFiddleFastly.svg?branch=master)](https://travis-ci.org/GannettDigital/goFiddleFastly)
[![codecov](https://codecov.io/gh/GannettDigital/goFiddleFastly/branch/master/graph/badge.svg)](https://codecov.io/gh/GannettDigital/goFiddleFastly)

This is a simple Go client for using Fastly Fiddle to do unit testing. To use goFastlyFiddle, create a client and your first Fiddle:
```go
client, _ := fiddle.DefaultClient()
```

You can now create a new Fiddle:
```go
myFastlyFiddle, err := client.CreateFiddle(&fiddle.CreateFiddleInput{
    Origins: []string{"https://httpbin.org"},
    Vcl: fiddle.Vcl{
        Recv:  "set req.url = querystring.sort(req.url);\nif (req.url.path == \"/robots.txt\") {\nerror 901;\n}",
        Error: "if (obj.status == 901) {\nset obj.status = 200;\nset obj.response = \"OK\";synthetic \"User-agent: BadBot\" LF \"Disallow: /\";\nreturn(deliver);\n}",
    },
})
```

The Fiddle will be associated with an ID that is returned to you. To update your Fiddle, provide the ID:
```go
myFastlyFiddle, err = client.UpdateFiddle(&fiddle.UpdateFiddleInput{
    ID:      myFastlyFiddle.ID,
    Origins: []string{"https://httpbin.org"},
    Vcl: fiddle.Vcl{
        Recv:  "if (req.http.Fastly-FF) {set req.max_stale_while_revalidate = 0s;}",
        Hit:   "if (!obj.cacheable) {return(pass);}",
        Error: "if (obj.status == 901) {\nset obj.status = 200;\nset obj.response = \"OK\";synthetic \"User-agent: BadBot\" LF \"Disallow: /\";\nreturn(deliver);\n}",
    },
    ReqURL: "/some/uri/to/hit",
})
```

Once you've adjusted your Fiddle, you can run it and use the results in your tests:
```go
results, err := client.ExecuteFiddle(&fiddle.ExecuteFiddleInput{
	ID: myFastlyFiddle.ID,
})
```

License
-------
```
Copyright 2018 Gannett Co., Inc

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
```