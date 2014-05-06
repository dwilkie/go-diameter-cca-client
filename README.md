# go-diameter-cca-client

A Diameter Credit-Control Application Client written in Go

## Development

Use [forego](https://github.com/ddollar/forego) for development.

### Start the test server

```
cd src/github.com/fiorix/go-diameter/examples/server
go run server
```

### Run the client

```
cd go-diameter-cca-client
go get
forego start
```

## Resources

* [Installing Go](http://blog.labix.org/2013/06/15/in-flight-deb-packages-of-go)
* [Getting Started with Go on Heroku](http://mmcgrana.github.io/2012/09/getting-started-with-go-on-heroku.html)
* [Diameter Credit-Control Application RFC] (http://tools.ietf.org/html/rfc4006)
* [go-diameter](https://github.com/fiorix/go-diameter)
* [go-redis](https://github.com/fiorix/go-redis)

## Deployment

### Updating

```
cd go
go get -u github.com/dwilkie/go-diameter-cca-client
```

### Testing Live Connection

```
SERVER_ADDRESS=host:port go-diameter-cca-client
```

See also [forego#20](https://github.com/ddollar/forego/issues/20)

We might be able to create something similar to what [foreman](https://github.com/ddollar/foreman/blob/master/lib/foreman/export/upstart.rb) does
