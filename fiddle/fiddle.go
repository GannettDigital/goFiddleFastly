package fiddle

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

// Fiddle represents an existing Fastly Fiddle.
type Fiddle struct {
	ID          string `json:"id"`
	StartTime   string `json:"startTime"`
	Status      int    `json:"status"`
	ClientFetch struct {
		Req  string `json:"req"`
		Resp string `json:"resp"`
	} `json:"clientFetch"`
	OriginFetches         map[string]interface{} `json:"originFetches"`
	RespBodyIsText        bool                   `json:"respBodyIsText"`
	RespBodyType          string                 `json:"respBodyType"`
	RespBodyPreview       string                 `json:"respBodyPreview"`
	RespBodyBytesReceived int                    `json:"respBodyBytesReceived"`
	RespBodyChunkCount    int                    `json:"respBodyChunkCount"`
	RespComplete          bool                   `json:"respComplete"`
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

// Vcl contains information about the VCL you need uploaded to your Fiddle object.
type Vcl struct {
	Recv    string `json:"recv"`
	Hit     string `json:"hit"`
	Miss    string `json:"miss"`
	Pass    string `json:"pass"`
	Fetch   string `json:"fetch"`
	Error   string `json:"error"`
	Deliver string `json:"deliver"`
	Init    string `json:"init"`
}

// CreateFiddleInput contains the values you can specify to create a new Fiddle.
type CreateFiddleInput struct {
	Origins       []string `json:"origins"`
	ReqURL        string   `json:"reqUrl"`
	ReqMethod     string   `json:"reqMethod"`
	ReqHeaders    string   `json:"reqHeaders"`
	ReqBody       string   `json:"reqBody"`
	Vcl           Vcl      `json:"vcl"`
	PurgeFirst    bool     `json:"purgeFirst"`
	EnableCluster bool     `json:"enableCluster"`
	EnableShield  bool     `json:"enableShield"`
}

// NewFiddleResponse contains
type NewFiddleResponse struct {
	Fiddle Fiddle      `json:"fiddle"`
	Valid  bool        `json:"valid"`
	Errors interface{} `json:"errors"`
}

// CreateFiddle creates a new Fiddle and returns a Fiddle object.
func (c *Client) CreateFiddle(f *CreateFiddleInput) (*Fiddle, error) {
	headers := make(map[string]string)
	headers["Content-Type"] = "application/json"
	headers["Accept"] = "application/json"

	body, err := json.Marshal(f)
	if err != nil {
		return nil, err
	}

	resp, err := c.Post(&RequestInput{
		URI:     "/fiddle",
		Headers: headers,
		Body:    body,
	})
	if err != nil {
		return nil, err
	}

	readBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, nil
	}

	var createResp NewFiddleResponse
	json.Unmarshal(readBody, &createResp)

	if !createResp.Valid {
		return nil, fmt.Errorf("Fiddle not valid! %v", createResp.Errors)
	}

	return &createResp.Fiddle, nil
}

// UpdateFiddleInput contains the values you can specify to create a new Fiddle.
type UpdateFiddleInput struct {
	ID            string   `json:"id"`
	ReqURL        string   `json:"reqUrl"`
	ReqMethod     string   `json:"reqMethod"`
	ReqHeaders    string   `json:"reqHeaders"`
	ReqBody       string   `json:"reqBody"`
	Origins       []string `json:"origins"`
	Vcl           Vcl      `json:"vcl"`
	PurgeFirst    bool     `json:"purgeFirst"`
	EnableCluster bool     `json:"enableCluster"`
	EnableShield  bool     `json:"enableShield"`
}

// UpdateFiddle updates an existing Fiddle and returns a Fiddle object.
func (c *Client) UpdateFiddle(f *UpdateFiddleInput) (*Fiddle, error) {
	headers := make(map[string]string)
	headers["Content-Type"] = "application/json"
	headers["Accept"] = "application/json"

	body, err := json.Marshal(f)
	if err != nil {
		return nil, err
	}

	resp, err := c.Put(&RequestInput{
		URI:     fmt.Sprintf("/fiddle/%s", f.ID),
		Headers: headers,
		Body:    body,
	})
	if err != nil {
		return nil, err
	}

	readBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, nil
	}

	var updateResp NewFiddleResponse
	json.Unmarshal(readBody, &updateResp)

	if !updateResp.Valid {
		return nil, fmt.Errorf("Fiddle not valid! %v", updateResp.Errors)
	}

	return &updateResp.Fiddle, nil
}

// ExecuteFiddleInput contains what is needed to execute a request (currently only Fiddle ID).
type ExecuteFiddleInput struct {
	ID string `json:"id"`
}

// ExecuteFiddle uploads a request to the specified Fiddle and then runs it.
func (c *Client) ExecuteFiddle(e *ExecuteFiddleInput) (*StreamResponse, error) {
	headers := make(map[string]string)
	headers["Accept"] = "application/json"

	resp, err := c.Post(&RequestInput{
		URI:     fmt.Sprintf("/fiddle/%s/execute?cacheID=1", e.ID),
		Headers: headers,
	})
	if err != nil {
		return nil, err
	}

	readBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, nil
	}

	var session StreamInput
	json.Unmarshal(readBody, &session)

	session.FiddleID = e.ID

	return c.Stream(&session)
}

// DeleteFiddle removes code from the existing Fiddle to "delete" it.
func (c *Client) DeleteFiddle(ID string) bool {
	_, err := c.UpdateFiddle(&UpdateFiddleInput{
		ID:         ID,
		ReqURL:     "",
		ReqMethod:  "",
		ReqHeaders: "",
		ReqBody:    "",
		Origins:    []string{},
		Vcl: Vcl{
			Recv:    "",
			Hit:     "",
			Miss:    "",
			Pass:    "",
			Fetch:   "",
			Error:   "",
			Deliver: "",
			Init:    "",
		},
		PurgeFirst:    true,
		EnableCluster: true,
		EnableShield:  true,
	})

	return err == nil
}
