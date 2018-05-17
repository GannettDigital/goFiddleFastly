package fiddle

import (
	"testing"
)

var (
	fiddleID string
	c        *Client
)

func init() {
	c, _ = DefaultClient()
}

func TestCreateFiddle(t *testing.T) {
	tests := []struct {
		TestDescription string
		Input           CreateFiddleInput
		ExpectedErr     error
	}{
		{
			TestDescription: "Should create a Fiddle",
			Input: CreateFiddleInput{
				Origins: []string{"https://google.com"},
				Vcl: Vcl{
					Recv:  "if (req.http.Fastly-FF) {set req.max_stale_while_revalidate = 0s;}",
					Hit:   "if (!obj.cacheable) {return(pass);}",
					Error: "if (obj.status == 901) {set obj.status = 200;set obj.response = \"OK\";synthetic \"User-agent: BadBot\" LF \"Disallow: /\";return(deliver);}",
				},
			},
			ExpectedErr: nil,
		},
	}

	for _, test := range tests {
		fiddle, err := c.CreateFiddle(&test.Input)

		fiddleID = fiddle.ID

		if err != nil {
			t.Errorf("Fiddle creation should not return error, got %v, want %v", err, test.ExpectedErr)
		}
	}
}

func TestUpdateFiddle(t *testing.T) {
	tests := []struct {
		TestDescription string
		Input           UpdateFiddleInput
		ExpectedErr     error
	}{
		{
			TestDescription: "Should update existing Fiddle",
			Input: UpdateFiddleInput{
				ID:      fiddleID,
				Origins: []string{"https://httpbin.org"},
			},
			ExpectedErr: nil,
		},
	}

	for _, test := range tests {
		_, err := c.UpdateFiddle(&test.Input)

		if err != nil {
			t.Errorf("Fiddle update should not return error, got %v, want %v", err, test.ExpectedErr)
		}
	}
}

func TestExecuteFiddle(t *testing.T) {
	tests := []struct {
		TestDescription string
		Input           ExecuteFiddleInput
		ExpectedErr     error
	}{
		{
			TestDescription: "Should execute existing Fiddle",
			Input: ExecuteFiddleInput{
				ID: fiddleID,
			},
			ExpectedErr: nil,
		},
	}

	for _, test := range tests {
		_, err := c.ExecuteFiddle(&test.Input)

		if err != nil {
			t.Errorf("Fiddle update should not return error, got %v, want %v", err, test.ExpectedErr)
		}
	}
}

func TestDeleteFiddle(t *testing.T) {
	tests := []struct {
		TestDescription string
		Input           string
		ExpectedResp    bool
	}{
		{
			TestDescription: "Should delete existing Fiddle",
			Input:           fiddleID,
			ExpectedResp:    true,
		},
	}

	for _, test := range tests {
		delete := c.DeleteFiddle(test.Input)

		if !delete {
			t.Errorf("Fiddle delete should finish successfully, got %v, want %v", delete, test.ExpectedResp)
		}
	}
}
