package main

import (
  "fmt"
  "github.com/dwilkie/go-diameter-cca-client/client"
  "github.com/benmanns/goworker"
  //"github.com/garyburd/redigo/redis"
)

func init() {
  goworker.Register("BeelineChargeRequest", BeelineChargeRequest)
}

func BeelineChargeRequest(queue string, args ...interface{}) error {
  result := beeline.Charge("85560201158")
  fmt.Printf("Result %s", result)
  fmt.Printf("From %s, %v\n", queue, args)
//  conn.Send("RPUSH", fmt.Sprintf("resque:queue:%s", queue), YOURJSON)
  return nil
}
