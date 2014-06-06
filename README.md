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
go get -u github.com/fiorix/go-diameter/diam
godep update github.com/fiorix/go-diameter/diam
```

### Restore the dependencies back to upstream branch

Edit `client.go` and replace

```
"github.com/dwilkie/go-diameter/diam"
"github.com/dwilkie/go-diameter/diam/diamtype"
"github.com/dwilkie/go-diameter/diam/diamdict"
```

with

```
"github.com/fiorix/go-diameter/diam"
"github.com/fiorix/go-diameter/diam/diamtype"
"github.com/fiorix/go-diameter/diam/diamdict"
```

Then run:

```
go get -u github.com/fiorix/go-diameter/diam
godep save
```

### Checking the build

Push your changes then run:

```
cd go
go get -u github.com/dwilkie/go-diameter-cca-client
```

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

Note you'll only be able to see the response packets because the request packets are encrypted through the VPN tunnel. But you can see the request from the client after it builds the AVPs.

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

## Flags

When creating a new AVP in the client you set the flags like this:

```go
m.NewAVP("Vendor-Id", 0x40, 0x0, VendorId)
```

Here the flag is `0x40` which is `64` in decimal or `1000000` in binary. According to [this article](http://diameter-protocol.blogspot.com/2011/05/daimeter-avp-structure.html) the first bit is for the Manditory Flag. Therefore to toggle between Manditory ON and Manditory OFF change `0x40` to `0x00`

## Resources

* [Installing Go](http://blog.labix.org/2013/06/15/in-flight-deb-packages-of-go)
* [Getting Started with Go on Heroku](http://mmcgrana.github.io/2012/09/getting-started-with-go-on-heroku.html)
* [Diameter Credit-Control Application RFC] (http://tools.ietf.org/html/rfc4006)
* [go-diameter](https://github.com/fiorix/go-diameter)
* [go-redis](https://github.com/fiorix/go-redis)
* [go-worker](http://www.goworker.org/)
