package payloads

import (
	"net/http"
)

func XmlEndpoint(rw http.ResponseWriter, _ *http.Request) {
	rw.Header().Add("Content-Type", "application/xml; charset=utf-8")
	rw.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?>
<user type="admin">
  <name>Elliot</name>
  <social>
    <facebook>https://facebook.com</facebook>
    <twitter>https://twitter.com</twitter>
    <youtube>https://youtube.com</youtube>
  </social>
</user>`))
}
