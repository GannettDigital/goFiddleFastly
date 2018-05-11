package fiddle

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/donovanhide/eventsource"
)

// RequestInput contains options for each request. There are no default values provided.
type RequestInput struct {
	// URI is the URI that you want to hit, which will be appended to your Client Address.
	URI string

	// Headers is the map of all the headers needed for each request.
	Headers map[string]string

	// Body is a byte array of the body passed with the request.
	Body []byte
}

// Request creates a new request and executes it.
func (c *Client) Request(method string, r *RequestInput) (*http.Response, error) {
	req, err := http.NewRequest(method, fmt.Sprintf("%s%s", c.Address, r.URI), bytes.NewReader(r.Body))
	if err != nil {
		return nil, err
	}

	for header, value := range r.Headers {
		req.Header.Add(header, value)
	}

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// StreamResponse contains all of the information returned from a Fiddle execution stream.
type StreamResponse struct {
	ID          string `json:"id"`
	StartTime   string `json:"startTime"`
	Status      int    `json:"status"`
	ClientFetch struct {
		Req  string `json:"req"`
		Resp string `json:"resp"`
	} `json:"clientFetch"`
	OriginFetches struct {
	} `json:"originFetches"`
	RespBodyIsText        bool   `json:"respBodyIsText"`
	RespBodyType          string `json:"respBodyType"`
	RespBodyPreview       string `json:"respBodyPreview"`
	RespBodyBytesReceived int    `json:"respBodyBytesReceived"`
	RespBodyChunkCount    int    `json:"respBodyChunkCount"`
	RespComplete          bool   `json:"respComplete"`
	Events                []struct {
		Name   string `json:"name"`
		Reqkey string `json:"reqkey"`
		Time   string `json:"time"`
		Server struct {
			Datacenter string `json:"datacenter"`
			NodeID     string `json:"nodeID"`
		} `json:"server"`
		Restarts int    `json:"restarts,omitempty"`
		URL      string `json:"url,omitempty"`
		IsESI    bool   `json:"isESI,omitempty"`
		Return   string `json:"return,omitempty"`
		Status   int    `json:"status,omitempty"`
		State    string `json:"state,omitempty"`
		Hits     int    `json:"hits,omitempty"`
		EdgeDC   string `json:"edgeDC,omitempty"`
	} `json:"events"`
}

// StreamInput contains fields needed to get the Fiddle execute results stream
type StreamInput struct {
	// ID is returned from a POST to execute the Fiddle
	ID string `json:"sessionID"`

	// FiddleID is the ID of the Fiddle being run.
	FiddleID string
}

// Stream gets the streamed results from a Fiddle execution.
func (c *Client) Stream(s *StreamInput) (*StreamResponse, error) {
	stream, err := eventsource.Subscribe(fmt.Sprintf("%s/fiddle/%s/result-stream/%s", c.Address, s.FiddleID, s.ID), "")
	if err != nil {
		return nil, err
	}

	var results StreamResponse

	for i := 0; i < 100; i++ {
		ev := <-stream.Events
		if strings.Contains(ev.Event(), "updateResult") {
			json.Unmarshal([]byte(ev.Data()), &results)
			break
		}
		time.Sleep(5 * time.Second)
	}

	return &results, nil
}
