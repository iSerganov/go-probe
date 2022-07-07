package ffprobe

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_HTTPErrorCode(t *testing.T) {
	err := &ProbeError{Message: "Server returned 404 Not Found"}

	validate := assert.New(t)
	validate.Equal(404, err.HTTPErrorCode())

	err = &ProbeError{Message: "Server returned 503 Service not available"}
	validate.Equal(503, err.HTTPErrorCode())

	err = &ProbeError{Message: "Server returned 5xx Internal server error"}
	validate.Equal(0, err.HTTPErrorCode())
}

func Test_IsInternalServerError(t *testing.T) {
	err := &ProbeError{Message: "Server returned 404 Not Found"}

	validate := assert.New(t)
	validate.False(err.IsInternalServerError())

	err = &ProbeError{Message: "Server returned 5XX Internal server error"}
	validate.True(err.IsInternalServerError())

	err = &ProbeError{Message: "Server returned 7XX Internal server error"}
	validate.False(err.IsInternalServerError())
}
