package queue

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestActivationMessage_JSON(t *testing.T) {
	orig := ActivationMessage{
		Email:          "user@example.com",
		FirstName:      "John",
		ActivationLink: "http://example.com/activate/abc123",
	}
	data, err := json.Marshal(orig)
	require.NoError(t, err)

	var decoded ActivationMessage
	require.NoError(t, json.Unmarshal(data, &decoded))
	assert.Equal(t, orig, decoded)
}

func TestActivationMessage_JSONKeys(t *testing.T) {
	msg := ActivationMessage{Email: "a@b.com", FirstName: "X", ActivationLink: "http://link"}
	data, _ := json.Marshal(msg)

	var raw map[string]string
	require.NoError(t, json.Unmarshal(data, &raw))
	assert.Contains(t, raw, "email")
	assert.Contains(t, raw, "first_name")
	assert.Contains(t, raw, "activation_link")
}

func TestPasswordResetMessage_JSON(t *testing.T) {
	orig := PasswordResetMessage{
		Email:     "user@example.com",
		FirstName: "Jane",
		ResetLink: "http://example.com/reset/xyz",
	}
	data, err := json.Marshal(orig)
	require.NoError(t, err)

	var decoded PasswordResetMessage
	require.NoError(t, json.Unmarshal(data, &decoded))
	assert.Equal(t, orig, decoded)
}

func TestPasswordResetMessage_JSONKeys(t *testing.T) {
	msg := PasswordResetMessage{Email: "a@b.com", FirstName: "X", ResetLink: "http://link"}
	data, _ := json.Marshal(msg)

	var raw map[string]string
	require.NoError(t, json.Unmarshal(data, &raw))
	assert.Contains(t, raw, "email")
	assert.Contains(t, raw, "first_name")
	assert.Contains(t, raw, "reset_link")
}

func TestPasswordConfirmationMessage_JSON(t *testing.T) {
	orig := PasswordConfirmationMessage{
		Email:     "user@example.com",
		FirstName: "Bob",
	}
	data, err := json.Marshal(orig)
	require.NoError(t, err)

	var decoded PasswordConfirmationMessage
	require.NoError(t, json.Unmarshal(data, &decoded))
	assert.Equal(t, orig, decoded)
}
