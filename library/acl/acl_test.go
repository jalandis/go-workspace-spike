package acl_test

import (
	"testing"

	"blab.com/library/acl"
	"github.com/stretchr/testify/assert"
)

func TestHasAccess(t *testing.T) {
	assert.Equal(t, false, acl.HasAccess("john"))
	assert.Equal(t, true, acl.HasAccess("not-john"))
}
