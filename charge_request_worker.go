package main

import (
  "fmt"
  "github.com/dwilkie/go-diameter-cca-client/client"
  "github.com/benmanns/goworker"
)

func init() {
  goworker.Register("ChargeRequestWorker", ChargeRequestWorker)
}

func ChargeRequestWorker(queue string, args ...interface{}) error {
  beeline.Charge("foo");
  fmt.Println("Hello, world!")
  return nil
}
