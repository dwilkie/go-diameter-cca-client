# go-diameter-cca-client

A Diameter Credit-Control Application Client written in Go

## Development

Use [forego](https://github.com/ddollar/forego) for development.

### Start the test server

```
cd src/github.com/fiorix/go-diameter/examples/server
go run server.go
```

### Run the client

```
cd go-diameter-cca-client
go get
forego start
```

### Update the dependencies

```
go get -u github.com/dwilkie/go-diameter/diam
godep update github.com/dwilkie/go-diameter/diam
```

### Restore the dependencies back to upstream branch

Edit `client.go` and replace

```
"github.com/dwilkie/go-diameter/diam"
"github.com/dwilkie/go-diameter/diam/datatypes"
```

with

```
"github.com/fiorix/go-diameter/diam"
"github.com/fiorix/go-diameter/diam/datatypes"
```

Then run:

```
go get -u github.com/fiorix/go-diameter-cca-client
godep restore
godep update github.com/fiorix/go-diameter/diam
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
SERVER_ADDRESS=192.168.3.20:3868 go-diameter-cca-client
```

### Inspecting Packets with Wireshark

On the server:

```
tcpdump -i eth0 -nnvvS host 192.168.3.20 -w beeline_diameter.cap
````

On your local machine:

```
sftp -i ~/.ssh/aws/dwilkie.pem ubuntu@nuntium.chibitxt.me:beeline_diameter.cap .
```

Then open `beeline_diameter.cap` with Wireshark and inspect the response.

### Foreman

See also [forego#20](https://github.com/ddollar/forego/issues/20)

We might be able to create something similar to what [foreman](https://github.com/ddollar/foreman/blob/master/lib/foreman/export/upstart.rb) does
