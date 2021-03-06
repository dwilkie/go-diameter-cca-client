# go-diameter-cca-client

A Diameter Credit-Control Application Client written in Go

## How it works

The `Procfile` contains the command to run the binary which is compiled from the `main` package and named after the project name `go-diameter-cca-client`.

The the `main()` function is containted within `worker.go` which calls `goworker.Work()`.

In `beeline_charge_request.go` which is also in the `main` package, the `init()` function registers the `BeelineChargeRequest` worker. The `BeelineChargeRequest` function is also defined within `beeline_charge_request.go` and accepts the `queue` and the `args` for the job. From the args it extracts the `transaction_id` and `mobile_number` and calls `beeline.Charge(transaction_id, mobile_number)` which returns the `session_id` and `result_code` of the charge request. It then uses the `redisurl` package to connect to the redis server and enqueues a job to the `beeline_charge_request_updater` queue with the `session_id` and `result_code`.

Inside the file `client/beeline.go` the package `beeline` is defined. The function `Charge()` which is called by the function `BeelineChargeRequest` (explained above) accepts a `transaction_id` and a `mobile_number` and returns the `session_id` and `result_code`. It first creates a `Parser` which loads the definitions of `diamdict.DefaultXML` and `diamdict.CreditControlXML`. It then defines a handler which extracts the `session_id` and `result_code` from the `CCA` response. It then connects to the `ServerAddress` passing the defined `handler` to the `Dial()` command. The `Dial()` command return the `connection` and `NewClient()` is called with the `connection`, `transaction_id` and `mobile_number`.

`NewClient` builds a `CER` request and writes the request to the `connection`. It then builds a `CCR` request using the `transaction_id` and the `mobile_number` and writes the request to the `connection`. The responses are handled by the `handler` which extracts the `result_code` and `session_id` from the `CCA`.

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

#### Beeline

##### Example Client

```
SERVER_ADDRESS=192.168.3.20:3868 go run go/src/github.com/dwilkie/go-diameter-cca-client/examples/beeline.go
```

##### Worker

###### Using Foreman

```
forego start -e .env.production
```

###### Manually

Manually start the worker with the following command:

```
REDIS_URL=redis_uri BEELINE_CHARGE_REQUEST_UPDATER_QUEUE=beeline_charge_request_updater_queue BEELINE_CHARGE_REQUEST_UPDATER_WORKER=BeelineChargeRequestUpdater ./go-diameter-cca-client -queues="beeline_charge_request_queue" -uri $REDIS_URL
```

where:

* `REDIS_URL` is the full Redis URL including Authorization, Scheme and PORT. e.g. `redis://redis-user:redis_password@redis-host:port`
* `BEELINE_CHARGE_REQUEST_UPDATER_QUEUE` is the name of the queue for which jobs are pushed which handles updating charge requests for Beeline.
* `BEELINE_CHARGE_REQUEST_UPDATER_WORKER` is the name of the worker class which handles updating charge requests for Beeline.
* `-queues` is the name of the queue which handles sending charge requests to Beeline.

#### Inspecting Packets with Wireshark

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
