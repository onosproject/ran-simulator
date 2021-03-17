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

package grpcinterceptors

import (
	"context"
	"fmt"
	"github.com/grpc-ecosystem/go-grpc-middleware/util/metautils"
	"github.com/onosproject/onos-lib-go/pkg/auth"
	"strings"

	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
)

const (
	// ContextMetadataTokenKey metadata token key
	ContextMetadataTokenKey = "bearer"
)

// AuthenticationInterceptor an interceptor for authentication
func AuthenticationInterceptor(ctx context.Context) (context.Context, error) {
	// Extract token from metadata in the context
	tokenString, err := grpc_auth.AuthFromMD(ctx, ContextMetadataTokenKey)
	if err != nil {
		return nil, err
	}

	// Authenticate the jwt token
	jwtAuth := new(auth.JwtAuthenticator)
	authClaims, err := jwtAuth.ParseAndValidate(tokenString)
	if err != nil {
		return ctx, err
	}

	niceMd := metautils.ExtractIncoming(ctx)
	niceMd.Del("authorization")
	if name, ok := authClaims["name"]; ok {
		niceMd.Set("name", name.(string))
	}
	if email, ok := authClaims["email"]; ok {
		niceMd.Set("email", email.(string))
	}
	if aud, ok := authClaims["aud"]; ok {
		niceMd.Set("aud", aud.(string))
	}
	if exp, ok := authClaims["exp"]; ok {
		niceMd.Set("exp", fmt.Sprintf("%s", exp))
	}
	if iat, ok := authClaims["iat"]; ok {
		niceMd.Set("iat", fmt.Sprintf("%s", iat))
	}
	if iss, ok := authClaims["iss"]; ok {
		niceMd.Set("iss", iss.(string))
	}
	if sub, ok := authClaims["sub"]; ok {
		niceMd.Set("sub", sub.(string))
	}
	if atHash, ok := authClaims["at_hash"]; ok {
		niceMd.Set("at_hash", atHash.(string))
	}

	groupsIf, ok := authClaims["groups"].([]interface{})
	if ok {
		groups := make([]string, 0)
		for _, g := range groupsIf {
			groups = append(groups, g.(string))
		}
		niceMd.Set("groups", strings.Join(groups, ";"))
	}
	return niceMd.ToIncoming(ctx), nil
}
