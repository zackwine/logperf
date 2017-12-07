package main

import (
  "time"
  "math/rand"
)

// based on 
// https://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-golang

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
    letterIdxBits = 6                    // 6 bits to represent a letter index
    letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
    letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

type RandStringGen struct {
  Source rand.Source
}

func NewRandStringGen() (*RandStringGen) {
  r := &RandStringGen{}
  r.Source = rand.NewSource(time.Now().UnixNano())
  return r
}

func (r *RandStringGen) RandString(n int) string {
    b := make([]byte, n)
    // A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
    for i, cache, remain := n-1, r.Source.Int63(), letterIdxMax; i >= 0; {
        if remain == 0 {
            cache, remain = r.Source.Int63(), letterIdxMax
        }
        if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
            b[i] = letterBytes[idx]
            i--
        }
        cache >>= letterIdxBits
        remain--
    }

    return string(b)
}
