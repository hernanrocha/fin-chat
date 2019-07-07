package stooq

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStooqHandler(t *testing.T) {
	res, err := StooqHandler("AAPL")
	assert.Nil(t, err)
	assert.Contains(t, res, "per share")
}
