package stooq

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStooqHandler(t *testing.T) {
	res, err := Handle(StooqRequest{ Code: "AAPL" })
	assert.Nil(t, err)
	assert.Contains(t, res.Result, "per share")
}
