package main

import (
  "github.com/dwilkie/go-diameter-cca-client/client"
  "github.com/benmanns/goworker"
)

func init() {
  goworker.Register("ChargeRequestWorker", ChargeRequestWorker)
}

func ChargeRequestWorker(queue string, args ...interface{}) error {
  beeline.Charge("85560201158");
  return nil
}
