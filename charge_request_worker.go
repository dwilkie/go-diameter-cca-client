package main

import (
  "fmt"
  "github.com/dwilkie/go-diameter-cca-client/client/beeline"
  "github.com/benmanns/goworker"
)

func init() {
  goworker.Register("ChargeRequestWorker", ChargeRequestWorker)
}

func ChargeRequestWorker(queue string, args ...interface{}) error {
  fmt.Println("Hello, world!")
  return nil
}
