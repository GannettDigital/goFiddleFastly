package fiddle

import (
	"testing"

	"github.com/hashicorp/go-cleanhttp"
)

func TestDefaultClient(t *testing.T) {
	desired := Client{
		Address:    "https://fiddle.fastlydemo.net",
		HTTPClient: cleanhttp.DefaultClient(),
	}

	actual, err := DefaultClient()

	if err != nil {
		t.Errorf("Default client should not return error, got %v, want %v", err, nil)
	}

	if desired.Address != actual.Address {
		t.Errorf("Default client did not return expected response, got %s, want %s", desired.Address, actual.Address)
	}
}
