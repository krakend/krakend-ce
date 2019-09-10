package krakend

import (
	"net/http"
	"github.com/devopsfaith/krakend/proxy"
)

func combineStatusAndData(total int, parts []*proxy.Response) *proxy.Response {
	isComplete := len(parts) == total
	var retResponse *proxy.Response
	var retStatus int

	for _, part := range parts {
		if part == nil || part.Data == nil {
			isComplete = false
			continue
		}

		isComplete = isComplete && part.IsComplete
		if retResponse == nil {
			retResponse = part
			continue
		}

		status := part.Metadata.StatusCode
		if status != 0 { // failed
			if retStatus == 0 {
				retStatus = status
			} else if retStatus != status {
				retStatus = http.StatusInternalServerError
			}
		}

		for k, v := range part.Data {
			retResponse.Data[k] = v
		}
	}

	if retStatus > http.StatusInternalServerError {
		retStatus = http.StatusInternalServerError
	}

	if retResponse == nil {
		// do not allow nil data in the response:
		return &proxy.Response{
			Data: make(map[string]interface{}, 0),
			IsComplete: isComplete,
			Metadata: proxy.Metadata{StatusCode: retStatus},
		}
	}

	retResponse.IsComplete = isComplete
	retResponse.Metadata.StatusCode = retStatus
	return retResponse
}
