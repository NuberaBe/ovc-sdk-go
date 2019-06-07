package ovc

import (
	"fmt"
	"io/ioutil"
	"log"
	"testing"
	"time"

	jwtLib "github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
)

const (
	jwtPubKey  = "./rsc_test/jwt_pub.pem"
	jwtPrivKey = "./rsc_test/jwt_key.pem"
)

func init() {
	b, err := ioutil.ReadFile(jwtPubKey)
	if err != nil {
		log.Fatal(err)
	}
	err = SetJWTPublicKey(string(b))
	if err != nil {
		log.Fatal(err)
	}
}

func TestNewJWT(t *testing.T) {
	tokenStr, err := createJWT(t, time.Hour*24, "", nil)
	assert.NoError(t, err)
	assert.NotNil(t, tokenStr)

	// Not supported idProvider
	j, err := NewJWT(tokenStr, "FOO")
	assert.Error(t, err)
	assert.Nil(t, j)

	// Valid but not refreshable
	j, err = NewJWT(tokenStr, "IYO")
	assert.NoError(t, err)
	assert.NotNil(t, j)
	assert.False(t, j.refreshable)

	res, err := j.Get()
	assert.NoError(t, err)
	assert.Equal(t, tokenStr, res)
}

func TestRefreshable(t *testing.T) {
	// Valid and refreshable token
	createClaims := make(map[string]string)
	createClaims["refresh_token"] = "foobar"
	tokenStr, err := createJWT(t, time.Hour*24, "", createClaims)
	assert.NoError(t, err)

	j, err := NewJWT(tokenStr, "IYO")
	assert.NoError(t, err)
	assert.NotNil(t, j)

	assert.True(t, j.refreshable)
}

func TestGetClaim(t *testing.T) {
	addClaims := make(map[string]string)
	addClaims["hello"] = "world"
	addClaims["foo"] = "bar"

	tokenStr, err := createJWT(t, time.Hour*24, "", addClaims)
	assert.NoError(t, err)
	assert.NotNil(t, tokenStr)

	j, err := NewJWT(tokenStr, "IYO")
	assert.NoError(t, err)
	assert.NotNil(t, j)

	for claim, val := range addClaims {
		resVal, err := j.Claim(claim)
		assert.NoError(t, err)
		assert.Equal(t, val, resVal.(string))
	}

	missingClaims := []string{"lorem", "ipsum"}
	for _, claim := range missingClaims {
		resVal, err := j.Claim(claim)
		assert.Nil(t, resVal)
		assert.Error(t, err)
		assert.Equal(t, ErrClaimNotPresent, err)
	}
}

func TestIsExpired(t *testing.T) {
	// Expired an hour ago
	tokenStr, err := createJWT(t, -time.Hour, "", nil)
	assert.NoError(t, err)
	assert.NotNil(t, tokenStr)

	token, err := parseJWT(tokenStr)
	assert.NoError(t, err)
	assert.NotNil(t, token)

	res := isExpired(token)
	assert.True(t, res)

	// Expired now
	tokenStr, err = createJWT(t, 0, "", nil)
	assert.NoError(t, err)
	assert.NotNil(t, tokenStr)

	token, err = parseJWT(tokenStr)
	assert.NoError(t, err)
	assert.NotNil(t, token)

	res = isExpired(token)
	assert.True(t, res)

	// Expired a second less than expiration buffer
	tokenStr, err = createJWT(t, expirationBuffer-time.Second, "", nil)
	assert.NoError(t, err)
	assert.NotNil(t, tokenStr)

	token, err = parseJWT(tokenStr)
	assert.NoError(t, err)
	assert.NotNil(t, token)

	res = isExpired(token)
	assert.True(t, res)

	// Expires 5 seconds after expiration buffer
	tokenStr, err = createJWT(t, expirationBuffer+5*time.Second, "", nil)
	assert.NoError(t, err)
	assert.NotNil(t, tokenStr)

	token, err = parseJWT(tokenStr)
	assert.NoError(t, err)
	assert.NotNil(t, token)

	res = isExpired(token)
	assert.False(t, res)
}

func TestRefresh(t *testing.T) {
	// Check if refresh func is triggered

	createClaims := make(map[string]string)
	createClaims["refresh_token"] = "foobar"
	tokenStr, err := createJWT(t, 0, "", createClaims)
	assert.NoError(t, err)
	assert.NotNil(t, tokenStr)

	j, err := NewJWT(tokenStr, "IYO")
	assert.NoError(t, err)
	assert.NotNil(t, j)
	assert.True(t, j.refreshable)

	// Refresh errored
	j.refreshFunc = func(token string) (string, error) {
		return "", fmt.Errorf("An error occurred")
	}
	res, err := j.Get()
	assert.Error(t, err)
	assert.Equal(t, tokenStr, res, "The original token should still be returned when an error occurred")

	// Invalid token returned
	j.refreshFunc = func(token string) (string, error) {
		return "Foo.Bar.Token", nil
	}
	res, err = j.Get()
	assert.Error(t, err)
	assert.Equal(t, tokenStr, res, "The original token should still be returned when an error occurred")

	// Valid token returned
	refreshedToken, err := createJWT(t, expirationBuffer+5*time.Second, "", nil)
	assert.NoError(t, err)
	j.refreshFunc = func(token string) (string, error) {
		return refreshedToken, nil
	}
	res, err = j.Get()
	assert.NoError(t, err)
	assert.Equal(t, refreshedToken, res, "Token should now be the one returned from the refresh func")
}

// CreateJWT generates a JWT
func createJWT(t *testing.T, timeValid time.Duration, scopes string, additionalClaims map[string]string) (string, error) {
	b, err := ioutil.ReadFile(jwtPrivKey)
	assert.NoError(t, err)

	key, err := jwtLib.ParseECPrivateKeyFromPEM(b)
	assert.NoError(t, err)

	claims := jwtLib.MapClaims{
		"exp":   time.Now().Add(timeValid).Unix(),
		"scope": scopes,
	}

	for k, v := range additionalClaims {
		claims[k] = v
	}

	token := jwtLib.NewWithClaims(jwtLib.SigningMethodES384, claims)

	return token.SignedString(key)
}
