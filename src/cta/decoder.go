package cta

import (
  "encoding/json"
  "errors"
  "io"
  "net/http"
)

func DecodeCTA(k, stopID, url string) (CTABusFeedMessage, error) {
  bf := CTABusFeedMessage{}
  client := http.Client{}

  req, err := http.NewRequest("GET", url, nil)
  if err != nil {
    return bf, err
  }

  q := req.URL.Query()
  q.Add("key", k)
  q.Add("format", "json")
  q.Add("stpid", stopID)
  req.URL.RawQuery = q.Encode()

  resp, err := client.Do(req)
  if err != nil {
    return bf, err
  }
  defer resp.Body.Close()

  // read response code
  // TODO: make more robust
  if resp.StatusCode >= 400 {
    return bf, errors.New(http.StatusText(resp.StatusCode))
  }

  body, err := io.ReadAll(resp.Body)
  if err != nil {
    return bf, err
  }

  err = json.Unmarshal(body, &bf)

  return bf, err
}
