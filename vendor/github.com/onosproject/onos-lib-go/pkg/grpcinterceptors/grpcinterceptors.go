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

	"github.com/onosproject/onos-lib-go/pkg/auth"

	"github.com/onosproject/onos-lib-go/pkg/logging"

	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
)

const (
	// ContextMetadataTokenKey metadata token key
	ContextMetadataTokenKey = "bearer"
)

var log = logging.GetLogger("interceptors")

// AuthenticationInterceptor an interceptor for authentication
func AuthenticationInterceptor(ctx context.Context) (context.Context, error) {
	log.Info("Authenticating the user")

	// Extract token from metadata in the context
	tokenString, err := grpc_auth.AuthFromMD(ctx, ContextMetadataTokenKey)
	if err != nil {
		return nil, err
	}

	// Authenticate the jwt token
	jwtAuth := new(auth.JwtAuthenticator)
	_, err = jwtAuth.ParseAndValidate(tokenString)
	if err != nil {
		return ctx, err
	}

	return ctx, nil

}
