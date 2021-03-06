// Copyright 2013-2014 go-diameter authors.  All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

// Diameter client.

package main

import (
  "flag"
  "log"
  "math/rand"
  "os"
  "net"
  "bytes"
  "fmt"

  "github.com/fiorix/go-diameter/diam"
  "github.com/fiorix/go-diameter/diam/diamtype"
  "github.com/fiorix/go-diameter/diam/diamdict"
)

const (
  Identity    = diamtype.DiameterIdentity("teletech1.client.com")
  Realm       = diamtype.DiameterIdentity("teletech.com")
  DestinationRealm = diamtype.DiameterIdentity("comverse.com")
  VendorId    = diamtype.Unsigned32(0)
  ProductName = diamtype.UTF8String("Chibi")
  AuthApplicationId = diamtype.Unsigned32(4)
  ServiceContextId = diamtype.UTF8String("CMVT-SVC@comverse.com")
  CCRequestType = diamtype.Enumerated(0x04)
  CCRequestNumber = diamtype.Unsigned32(0)
  RequestedAction = diamtype.Enumerated(0x00)
  SubscriptionIdType = diamtype.Enumerated(0x00) // E164
  SubscriptionIdData = diamtype.UTF8String("85560201158")
  ServiceIdentifier = diamtype.Unsigned32(0)
  ServiceParameterType1 = diamtype.Unsigned32(1)
  ServiceParameterValue1 = diamtype.OctetString("401")
  ServiceParameterType2 = diamtype.Unsigned32(2)
  ServiceParameterValue2 = diamtype.OctetString("401")
)

func main() {
  ssl := flag.Bool("ssl", false, "connect using SSL/TLS")
  parser, _ := diamdict.NewParser()
  parser.Load(bytes.NewReader(diamdict.DefaultXML))
  parser.Load(bytes.NewReader(diamdict.CreditControlXML))
  flag.Parse()
  var result_code string
  var session_id string
  // CCA incoming messages are handled here.
  diam.HandleFunc("CCA", func(c diam.Conn, m *diam.Message) {
    result_code_avp, err := m.FindAVP(268)
    result_code = result_code_avp.Data.String()
    if err != nil {
      log.Fatal(err)
    }

    session_id_avp, err := m.FindAVP(263)
    session_id = session_id_avp.Data.String()
    if err != nil {
      log.Fatal(err)
    }

    c.Close()
    log.Printf("Receiving message from %s", c.RemoteAddr().String())
    log.Println(m)
  })

  // Connect using the default handler and base.Dict.
  addr := os.Getenv("SERVER_ADDRESS")
  log.Println("Connecting to", addr)
  var (
    c   diam.Conn
    err error
  )
  if *ssl {
    c, err = diam.DialTLS(addr, "", "", nil, nil)
  } else {
    c, err = diam.Dial(addr, nil, parser)
  }
  if err != nil {
    log.Fatal(err)
  }
  go NewClient(c)
  // Wait until the server kick us out.
  <-c.(diam.CloseNotifier).CloseNotify()
  log.Println(result_code)
  log.Println(session_id)
  log.Println("Server disconnected.")
}

// NewClient sends a CER to the server and then a DWR every 10 seconds.
func NewClient(c diam.Conn) {
  // Build CCR

  parser, _ := diamdict.NewParser()
  parser.Load(bytes.NewReader(diamdict.DefaultXML))
  parser.Load(bytes.NewReader(diamdict.CreditControlXML))

  m := diam.NewRequest(257, 0, parser)
  // Add AVPs
  m.NewAVP("Origin-Host", 0x40, 0x00, Identity)
  m.NewAVP("Origin-Realm", 0x40, 0x00, Realm)
  m.NewAVP("Origin-State-Id", 0x40, 0x00, diamtype.Unsigned32(rand.Uint32()))
  m.NewAVP("Auth-Application-Id", 0x40, 0x00, AuthApplicationId)
  laddr := c.LocalAddr()
  ip, _, _ := net.SplitHostPort(laddr.String())
  m.NewAVP("Host-IP-Address", 0x40, 0x0, diamtype.Address(net.ParseIP(ip)))
  m.NewAVP("Vendor-Id", 0x40, 0x0, VendorId)
  m.NewAVP("Product-Name", 0x00, 0x0, ProductName)

  log.Printf("Sending message to %s", c.RemoteAddr().String())
  log.Println(m.String())
  // Send message to the connection
  if _, err := m.WriteTo(c); err != nil {
    log.Fatal("Write failed:", err)
  }

  m = diam.NewRequest(272, 4, parser)
  // Add AVPs
  m.NewAVP("Session-Id", 0x40, 0x00, diamtype.UTF8String(fmt.Sprintf("%v", rand.Uint32())))
  m.NewAVP("Origin-Host", 0x40, 0x00, Identity)
  m.NewAVP("Origin-Realm", 0x40, 0x00, Realm)
  m.NewAVP("Destination-Realm", 0x40, 0x00, DestinationRealm)
  m.NewAVP("Auth-Application-Id", 0x40, 0x0, AuthApplicationId)
  m.NewAVP("CC-Request-Type", 0x40, 0x0, CCRequestType)
  m.NewAVP("Service-Context-Id", 0x40, 0x0, ServiceContextId)
  m.NewAVP("Service-Identifier", 0x40, 0x0, ServiceIdentifier)
  m.NewAVP("CC-Request-Number", 0x40, 0x0, CCRequestNumber)
  m.NewAVP("Requested-Action", 0x40, 0x0, RequestedAction)
  m.NewAVP("Subscription-Id", 0x40, 0x00, &diam.Grouped{
    AVP: []*diam.AVP{
      // Subscription-Id-Type
      diam.NewAVP(450, 0x40, 0x0, SubscriptionIdType),
      // Subscription-Id-Data
      diam.NewAVP(444, 0x40, 0x0, SubscriptionIdData),
    },
  })
  m.NewAVP("Service-Parameter-Info", 0x40, 0x00, &diam.Grouped{
    AVP: []*diam.AVP{
      // Service-Parameter-Type
      diam.NewAVP(441, 0x40, 0x0, ServiceParameterType1),
      // Service-Parameter-Value
      diam.NewAVP(442, 0x40, 0x0, ServiceParameterValue1),
    },
  })
  m.NewAVP("Service-Parameter-Info", 0x40, 0x00, &diam.Grouped{
    AVP: []*diam.AVP{
      // Service-Parameter-Type
      diam.NewAVP(441, 0x40, 0x0, ServiceParameterType2),
      // Service-Parameter-Value
      diam.NewAVP(442, 0x40, 0x0, ServiceParameterValue2),
    },
  })

  log.Printf("Sending message to %s", c.RemoteAddr().String())
  log.Println(m.String())
  // Send message to the connection
  if _, err := m.WriteTo(c); err != nil {
    log.Fatal("Write failed:", err)
  }
}
