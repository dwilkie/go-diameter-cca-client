package main

import (
  "github.com/jrallison/go-workers"
  "github.com/dwilkie/go-diameter-cca-client/client"
)

func beelineChargeRequestJob(message *workers.Msg) {
  args := message.Args().GetIndex(0).Get("arguments").MustArray()
  transaction_id := args[0].(string)
  msisdn := args[1].(string)
  updater_queue := args[2].(string)
  updater_worker := args[3].(string)
  server_address := args[4].(string)

  session_id, result_code := beeline.Charge(transaction_id, msisdn, server_address)
  workers.Enqueue(updater_queue, updater_worker, []string{session_id, result_code})
}

func main() {
  workers.Configure(map[string]string{
    // location of redis instance
    "server":  "localhost:6379",
    // instance of the database
    "database":  "0",
    // number of connections to keep open with redis
    "pool":    "30",
    // unique process id for this instance of workers (for proper recovery of inprogress jobs on crash)
    "process": "1",
  })

  workers.Process("beeline_charge_request_queue", beelineChargeRequestJob, 10)
  // Add additional workers here

  workers.Run()
}
