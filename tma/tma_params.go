// Package tma provides functionality for handling Telegram Mini App parameters
// and user data. It includes structures and methods for parsing and accessing
// initialization data received from Telegram Mini Apps.
package tma

import (
	"strconv"
	"time"

	jsoniter "github.com/json-iterator/go"
)

// json is a package-level JSON encoder/decoder instance configured to be
// compatible with the standard library's JSON implementation while offering
// better performance. It's used for all JSON operations in this package.
var json = jsoniter.ConfigCompatibleWithStandardLibrary

// Params represents the initialization parameters received from a Telegram Mini App.
// These parameters contain information about the user, chat, and authentication.
//
// The struct tags define the field names for both JSON and MessagePack serialization,
// ensuring compatibility with different data formats.
type Params struct {
	// UserData contains the serialized user information in JSON format.
	// This field should be parsed using the User() method to access structured data.
	UserData     string    `json:"user" msgpack:"user"`

	// ChatInstance is a unique identifier for the chat session where the app was launched.
	// This helps distinguish between different chat instances.
	ChatInstance string    `json:"chat_instance" msgpack:"chat_instance"`

	// ChatType indicates the type of chat where the app was launched.
	// Common values include "private", "group", "channel", etc.
	ChatType     string    `json:"chat_type" msgpack:"chat_type"`

	// AuthDate represents the timestamp when the authentication occurred.
	// This is typically parsed from a Unix timestamp string.
	AuthDate     time.Time `json:"auth_date" msgpack:"auth_date"`

	// Hash is the verification hash used to validate the authenticity
	// of the received parameters. This should be verified before trusting the data.
	Hash         string    `json:"hash" msgpack:"hash"`
}

// User represents a Telegram user with all available information from the Mini Apps
// initialization. This includes basic profile information and chat-specific details.
//
// The struct tags ensure proper serialization to both JSON and MessagePack formats.
type User struct {
	// ID is the user's unique identifier in Telegram.
	ID                    int    `json:"id" msgpack:"id"`

	// FirstName is the user's first name as set in their Telegram profile.
	FirstName             string `json:"first_name" msgpack:"first_name"`

	// LastName is the user's last name (optional, may be empty).
	LastName              string `json:"last_name" msgpack:"last_name"`

	// UserName is the user's Telegram username in @handle format (optional).
	UserName              string `json:"username" msgpack:"username"`

	// PhotoURL is a link to the user's profile photo (optional).
	PhotoURL              string `json:"photo_url" msgpack:"photo_url"`

	// Language is the user's interface language code (e.g., "en", "ru").
	Language              string `json:"language_code" msgpack:"language_code"`

	// ChatType indicates the type of chat where the app was launched.
	// This mirrors the value from Params but is included here for convenience.
	ChatType              string `json:"chat_type" msgpack:"chat_type"`

	// ChatInstance is a unique identifier for the chat session.
	// This mirrors the value from Params but is included here for convenience.
	ChatInstance          string `json:"chat_instance" msgpack:"chat_instance"`

	// IsPremium indicates whether the user has an active Telegram Premium subscription.
	IsPremium             bool   `json:"is_premium" msgpack:"is_premium"`

	// AllowsWriteToPM indicates whether the bot is allowed to send
	// direct messages to this user.
	AllowsWriteToPM       bool   `json:"allows_write_to_pm" msgpack:"allows_write_to_pm"`

	// AddedToAttachmentMenu indicates whether the user has added
	// the bot to their attachment menu.
	AddedToAttachmentMenu bool   `json:"added_to_attachment_menu" msgpack:"added_to_attachment_menu"`
}

// User parses the UserData field and returns a structured User object.
//
// Returns:
//   - *User: A pointer to the parsed User structure
//   - error: Any error that occurred during JSON unmarshaling
//
// Example usage:
//
//	user, err := params.User()
//	if err != nil {
//	    // handle error
//	}
//	fmt.Printf("User ID: %d", user.ID)
func (p *Params) User() (*User, error) {
	var user User

	// Unmarshal the JSON-encoded user data into the User struct
	err := json.Unmarshal([]byte(p.UserData), &user)

	return &user, err
}

// set updates the specified field in the Params struct based on the provided key.
// It handles type conversion and parsing as needed for different field types.
//
// Parameters:
//   - key: The field name to set (must match exactly)
//   - value: The string value to set/parse
//
// Supported keys and value formats:
//   - "user": Sets the raw UserData string (should be valid JSON)
//   - "chat_instance": Sets the ChatInstance string directly
//   - "chat_type": Sets the ChatType string directly
//   - "auth_date": Parses the value as a Unix timestamp string (seconds since epoch)
//
// Note: This method silently ignores unsupported keys and parsing errors.
func (p *Params) set(key string, value string) {
	switch key {
	case "user":
		p.UserData = value
	case "chat_instance":
		p.ChatInstance = value
	case "chat_type":
		p.ChatType = value
	case "auth_date":
		// Attempt to parse the value as a Unix timestamp
		if timestamp, err := strconv.ParseInt(value, 10, 64); err == nil {
			p.AuthDate = time.Unix(timestamp, 0)
		}
		// Silently ignore parsing errors (AuthDate remains zero value)
	}
}