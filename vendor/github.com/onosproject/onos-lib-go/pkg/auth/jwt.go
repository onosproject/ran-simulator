// Copyright 2020-present Open Networking Foundation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package auth

import (
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	ecoidc "github.com/ericchiang/oidc"
	"github.com/onosproject/onos-lib-go/pkg/logging"
	"gopkg.in/square/go-jose.v2"

	"google.golang.org/grpc/status"

	"google.golang.org/grpc/codes"

	"github.com/dgrijalva/jwt-go"
)

var log = logging.GetLogger("jwt")

const (
	// SharedSecretKey shared secret key for signing a token
	SharedSecretKey = "SHARED_SECRET_KEY"
	// OIDCServerURL - will be accessed as Environment variable
	OIDCServerURL = "OIDC_SERVER_URL"
	// OpenidConfiguration is the discovery point on the OIDC server
	OpenidConfiguration = ".well-known/openid-configuration"
	// HS prefix for HS family algorithms
	HS = "HS"
	// RS prefix for RS family algorithms
	RS = "RS"
)

// JwtAuthenticator jwt authenticator
type JwtAuthenticator struct {
	publicKeys map[string][]byte
}

// ParseToken parse token and Ensure that the JWT conforms to the structure of a JWT.
func (j *JwtAuthenticator) parseToken(tokenString string) (*jwt.Token, jwt.MapClaims, error) {
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		// HS256, HS384, or HS512
		if strings.HasPrefix(token.Method.Alg(), HS) {
			key := os.Getenv(SharedSecretKey)
			return []byte(key), nil
			// RS256, RS384, or RS512
		} else if strings.HasPrefix(token.Method.Alg(), RS) {
			keyID, ok := token.Header["kid"]
			if !ok {
				return nil, status.Errorf(codes.Unauthenticated, "token header not found 'kid' (key ID)")
			}
			keyIDStr := keyID.(string)
			publicKey, ok := j.publicKeys[keyIDStr]
			if !ok {
				// Keys may have been refreshed on the server
				// Fetch them again and try once more before failing
				if err := j.refreshJwksKeys(); err != nil {
					return nil, status.Errorf(codes.Unauthenticated, "unable to refresh keys from ID provider %s", err)
				}
				// try again after refresh
				if publicKey, ok = j.publicKeys[keyIDStr]; !ok {
					return nil, status.Errorf(codes.Unauthenticated, "token has obsolete key ID %s", keyID)
				}
			}
			rsaPublicKey, err := jwt.ParseRSAPublicKeyFromPEM(publicKey)
			if err != nil {
				return nil, status.Errorf(codes.Unauthenticated, err.Error())
			}
			return rsaPublicKey, nil
		}
		return nil, status.Errorf(codes.Unauthenticated, "unknown signing algorithm: %s", token.Method.Alg())
	})

	return token, claims, err

}

// ParseAndValidate parse a jwt string token and validate it
func (j *JwtAuthenticator) ParseAndValidate(tokenString string) (jwt.MapClaims, error) {
	token, claims, err := j.parseToken(tokenString)
	if err != nil {
		log.Warnf("Error parsing token: %s", tokenString)
		log.Warnf("Error %s", err.Error())
		return nil, err
	}

	// Check the token is valid
	if !token.Valid {
		return nil, status.Errorf(codes.Unauthenticated, "token is not valid %v", token)
	}

	return claims, nil
}

// Connect back to the OpenIDConnect server to retrieve the keys
// They are rotated every 6 hours by default - we keep the keys in a cache
// It's a 2 step process
// 1) connect to $OIDCServerURL/.well-known/openid-configuration and retrieve the JSON payload
// 2) lookup the "keys" parameter and get keys from $OIDCServerURL/keys
// The keys are in a public key format and are converted to RSA Public Keys
func (j *JwtAuthenticator) refreshJwksKeys() error {
	oidcURL := os.Getenv(OIDCServerURL)

	client := new(http.Client)
	resOpenIDConfig, err := client.Get(fmt.Sprintf("%s/%s", oidcURL, OpenidConfiguration))
	if err != nil {
		return err
	}
	if resOpenIDConfig.Body != nil {
		defer resOpenIDConfig.Body.Close()
	}
	openIDConfigBody, readErr := ioutil.ReadAll(resOpenIDConfig.Body)
	if readErr != nil {
		return err
	}
	var openIDprovider ecoidc.Provider
	jsonErr := json.Unmarshal(openIDConfigBody, &openIDprovider)
	if jsonErr != nil {
		return err
	}
	resOpenIDKeys, err := client.Get(openIDprovider.JWKSURL)
	if err != nil {
		return err
	}
	if resOpenIDKeys.Body != nil {
		defer resOpenIDKeys.Body.Close()
	}
	bodyOpenIDKeys, readErr := ioutil.ReadAll(resOpenIDKeys.Body)
	if readErr != nil {
		return err
	}
	var jsonWebKeySet jose.JSONWebKeySet
	if err := json.Unmarshal(bodyOpenIDKeys, &jsonWebKeySet); err != nil {
		return err
	}

	if j.publicKeys == nil {
		j.publicKeys = make(map[string][]byte)
	}
	for _, key := range jsonWebKeySet.Keys {
		data, err := x509.MarshalPKIXPublicKey(key.Key)
		if err != nil {
			return err
		}
		block := pem.Block{
			Type:  "PUBLIC KEY",
			Bytes: data,
		}
		pemBytes := pem.EncodeToMemory(&block)
		j.publicKeys[key.KeyID] = pemBytes
	}

	return nil
}
