package main

import (
  "fmt"
  "errors"
  "os"
  "github.com/dwilkie/go-diameter-cca-client/client"
  "github.com/benmanns/goworker"
  "github.com/soveran/redisurl"
)

var (
  errorInvalidParam = errors.New("invalid param")
)

func init() {
  goworker.Register("BeelineChargeRequest", BeelineChargeRequest)
}

func BeelineChargeRequest(queue string, args ...interface{}) error {
  charge_request_updater_queue := os.Getenv("BEELINE_CHARGE_REQUEST_UPDATER_QUEUE")
  charge_request_updater_worker := os.Getenv("BEELINE_CHARGE_REQUEST_UPDATER_WORKER")

  c, err := redisurl.Connect()
  defer c.Close()

  transaction_id, ok := args[0].(string)
  if !ok {
    fmt.Println(errorInvalidParam)
    return errorInvalidParam
  }

  mobile_number, ok := args[1].(string)
  if !ok {
    fmt.Println(errorInvalidParam)
    return errorInvalidParam
  }

  session_id, result_code := beeline.Charge(transaction_id, mobile_number)

  json := fmt.Sprintf("{\"class\":\"%s\",\"args\":[\"%s\",\"%s\"]}", charge_request_updater_worker, session_id, result_code)
  queue_key := fmt.Sprintf("resque:queue:%s", charge_request_updater_queue)

  n, err := c.Do("RPUSH", queue_key, json)
  _ = n

  if err != nil {
    fmt.Println(err)
    return err
  }

  return nil
}
