/*
 * Copyright 2020-present Open Networking Foundation
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
/**
 * @fileoverview gRPC-Web generated client stub for ran.trafficsim
 * @enhanceable
 * @public
 */



import * as grpcWeb from 'grpc-web';

import * as github_com_OpenNetworkingFoundation_gmap$ran_api_types_types_pb from '../../../../../github.com/OpenNetworkingFoundation/gmap-ran/api/types/types_pb';

import {
  ListRoutesRequest,
  ListRoutesResponse,
  ListTowersRequest,
  ListUesRequest,
  ListUesResponse,
  MapLayoutRequest} from './trafficsim_pb';

export class TrafficClient {
  client_: grpcWeb.AbstractClientBase;
  hostname_: string;
  credentials_: null | { [index: string]: string; };
  options_: null | { [index: string]: string; };

  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: string; }) {
    if (!options) options = {};
    if (!credentials) credentials = {};
    options['format'] = 'text';

    this.client_ = new grpcWeb.GrpcWebClientBase(options);
    this.hostname_ = hostname;
    this.credentials_ = credentials;
    this.options_ = options;
  }

  methodInfoGetMapLayout = new grpcWeb.AbstractClientBase.MethodInfo(
    github_com_OpenNetworkingFoundation_gmap$ran_api_types_types_pb.MapLayout,
    (request: MapLayoutRequest) => {
      return request.serializeBinary();
    },
    github_com_OpenNetworkingFoundation_gmap$ran_api_types_types_pb.MapLayout.deserializeBinary
  );

  getMapLayout(
    request: MapLayoutRequest,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.Error,
               response: github_com_OpenNetworkingFoundation_gmap$ran_api_types_types_pb.MapLayout) => void) {
    return this.client_.rpcCall(
      this.hostname_ +
        '/ran.trafficsim.Traffic/GetMapLayout',
      request,
      metadata || {},
      this.methodInfoGetMapLayout,
      callback);
  }

  methodInfoListTowers = new grpcWeb.AbstractClientBase.MethodInfo(
    github_com_OpenNetworkingFoundation_gmap$ran_api_types_types_pb.Tower,
    (request: ListTowersRequest) => {
      return request.serializeBinary();
    },
    github_com_OpenNetworkingFoundation_gmap$ran_api_types_types_pb.Tower.deserializeBinary
  );

  listTowers(
    request: ListTowersRequest,
    metadata?: grpcWeb.Metadata) {
    return this.client_.serverStreaming(
      this.hostname_ +
        '/ran.trafficsim.Traffic/ListTowers',
      request,
      metadata || {},
      this.methodInfoListTowers);
  }

  methodInfoListRoutes = new grpcWeb.AbstractClientBase.MethodInfo(
    ListRoutesResponse,
    (request: ListRoutesRequest) => {
      return request.serializeBinary();
    },
    ListRoutesResponse.deserializeBinary
  );

  listRoutes(
    request: ListRoutesRequest,
    metadata?: grpcWeb.Metadata) {
    return this.client_.serverStreaming(
      this.hostname_ +
        '/ran.trafficsim.Traffic/ListRoutes',
      request,
      metadata || {},
      this.methodInfoListRoutes);
  }

  methodInfoListUes = new grpcWeb.AbstractClientBase.MethodInfo(
    ListUesResponse,
    (request: ListUesRequest) => {
      return request.serializeBinary();
    },
    ListUesResponse.deserializeBinary
  );

  listUes(
    request: ListUesRequest,
    metadata?: grpcWeb.Metadata) {
    return this.client_.serverStreaming(
      this.hostname_ +
        '/ran.trafficsim.Traffic/ListUes',
      request,
      metadata || {},
      this.methodInfoListUes);
  }

}

