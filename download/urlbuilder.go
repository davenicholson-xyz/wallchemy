package download

import (
	"net/url"
	"strconv"
	"strings"
)

type URLBuilder struct {
	baseURL string
	values  url.Values
}

func NewURL(baseURL string) *URLBuilder {
	return &URLBuilder{
		baseURL: baseURL,
		values:  make(url.Values),
	}
}

// Add string parameter (only if value isn't empty)
func (u *URLBuilder) AddString(key, value string) {
	if value != "" {
		u.values.Add(key, value)
	}
}

// Add int parameter (only if value isn't 0)
func (u *URLBuilder) AddInt(key string, value int) {
	if value != 0 {
		u.values.Add(key, strconv.Itoa(value))
	}
}

// Add bool parameter (only if true)
func (u *URLBuilder) AddBool(key string, value bool) {
	if value {
		u.values.Add(key, "true")
	}
}

// SetString sets or updates a string parameter (only if value isn't empty)
func (u *URLBuilder) SetString(key, value string) {
	if value != "" {
		u.values.Set(key, value) // Uses Set instead of Add to replace existing
	} else {
		u.values.Del(key) // Remove if empty value
	}
}

// SetInt sets or updates an int parameter (only if value isn't 0)
func (u *URLBuilder) SetInt(key string, value int) {
	if value != 0 {
		u.values.Set(key, strconv.Itoa(value)) // Uses Set instead of Add
	} else {
		u.values.Del(key) // Remove if zero value
	}
}

// SetBool sets or updates a bool parameter (only if true)
func (u *URLBuilder) SetBool(key string, value bool) {
	if value {
		u.values.Set(key, "true") // Uses Set instead of Add
	} else {
		u.values.Del(key) // Remove if false
	}

}

// GetString returns the first value associated with the given key.
// If the key doesn't exist, it returns an empty string.
func (u *URLBuilder) GetString(key string) string {
	return u.values.Get(key)
}

// GetInt returns the integer value associated with the given key.
// If the ke     y doesn't exist or the value can't be parsed, it returns 0.
func (u *URLBuilder) GetInt(key string) int {
	val := u.values.Get(key)
	if val == "" {
		return 0
	}
	intVal, err := strconv.Atoi(val)
	if err != nil {
		return 0
	}
	return intVal
}

// GetBool returns the boolean value associated with the given key.
// Returns true only if the value is "true" (case-insensitive).
func (u *URLBuilder) GetBool(key string) bool {
	val := u.values.Get(key)
	return strings.EqualFold(val, "true")
}

// GetAll returns all values associated with the given key.
func (u *URLBuilder) GetAll(key string) []string {
	return u.values[key]
}

// Has returns true if the key exists in the parameters.
func (u *URLBuilder) Has(key string) bool {
	_, exists := u.values[key]
	return exists
}

// Build the final URL string
func (u *URLBuilder) Build() string {
	if len(u.values) == 0 {
		return u.baseURL
	}
	return u.baseURL + "?" + u.values.Encode()
}

func (u *URLBuilder) BuildWithout(key string) string {
	newValues := make(url.Values)
	for k, v := range u.values {
		if k != key {
			newValues[k] = v
		}
	}

	if len(newValues) == 0 {
		return u.baseURL
	}
	return u.baseURL + "?" + newValues.Encode()
}

func (u *URLBuilder) Without(key string) *URLBuilder {
	newBuilder := NewURL(u.baseURL)

	for k, v := range u.values {
		if k != key {
			newBuilder.values[k] = v
		}
	}

	return newBuilder
}

func (u *URLBuilder) Clone() *URLBuilder {
	clone := NewURL(u.baseURL)

	for k, v := range u.values {
		clone.values[k] = append([]string{}, v...)
	}

	return clone
}
