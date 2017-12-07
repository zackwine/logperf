package main

import (
  "bytes"
  "io/ioutil"
  "log"
  "net/http"
  "time"
)


type HttpOutputPool struct {
  httpClient *http.Client
}


func NewHttpOutputPool(maxConnectsPerHost int, requestTimeout int) (*HttpOutputPool) {
  h := &HttpOutputPool{}
  h.httpClient = &http.Client{
        Transport: &http.Transport{
            MaxIdleConnsPerHost: maxConnectsPerHost,
        },
        Timeout: time.Duration(requestTimeout) * time.Second,
  }
  return h
}

func(h *HttpOutputPool) SendMessage( url string, method string, msg string) {

  req, err := http.NewRequest(method, url, bytes.NewBuffer([]byte(msg)))
  if err != nil {
      log.Fatalf("Error Occured. %+v", err)
  }
  req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
  response, err := h.httpClient.Do(req)
  if err != nil && response == nil {
      log.Fatalf("Error sending request to url (%s). %+v", url, err)
  } else {
      // Close the connection to reuse it
      defer response.Body.Close()

      // Let's check if the work actually is done
      // We have seen inconsistencies even when we get 200 OK response
      body, err := ioutil.ReadAll(response.Body)
      if err != nil {
          log.Fatalf("Couldn't parse response body. %+v", err)
      }

      log.Println("Response Body:", string(body))
  }
}
