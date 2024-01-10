package payloads

import (
	"net/http"
)

func SymmetricJWKEndpoint(rw http.ResponseWriter, _ *http.Request) {
	rw.Header().Add("Content-Type", "application/json")
	rw.Write([]byte(`{
  "keys": [
    {
      "kty": "oct",
      "alg": "A128KW",
      "k": "GawgguFyGrWKav7AX4VKUg",
      "kid": "sim1"
    },
    {
      "kty": "oct",
      "k": "AyM1SysPpbyDfgZld3umj1qzKObwVMkoqQ-EstJQLr_T-1qS0gZH75aKtMN3Yj0iPS4hcgUuTwjAzZr1Z9CAow",
      "kid": "sim2",
      "alg": "HS256"
    }
  ]
}`))
}
