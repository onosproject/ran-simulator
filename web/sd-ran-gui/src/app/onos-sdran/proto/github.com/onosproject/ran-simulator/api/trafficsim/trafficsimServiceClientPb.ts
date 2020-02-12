/**
 * @fileoverview gRPC-Web generated client stub for ran.trafficsim
 * @enhanceable
 * @public
 */

// GENERATED CODE -- DO NOT EDIT!


import * as grpcWeb from 'grpc-web';

import * as github_com_onosproject_ran$simulator_api_types_types_pb from '../../../../../github.com/onosproject/ran-simulator/api/types/types_pb';

import {
  ListRoutesRequest,
  ListRoutesResponse,
  ListTowersRequest,
  ListTowersResponse,
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
    github_com_onosproject_ran$simulator_api_types_types_pb.MapLayout,
    (request: MapLayoutRequest) => {
      return request.serializeBinary();
    },
    github_com_onosproject_ran$simulator_api_types_types_pb.MapLayout.deserializeBinary
  );

  getMapLayout(
    request: MapLayoutRequest,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.Error,
               response: github_com_onosproject_ran$simulator_api_types_types_pb.MapLayout) => void) {
    return this.client_.rpcCall(
      this.hostname_ +
        '/ran.trafficsim.Traffic/GetMapLayout',
      request,
      metadata || {},
      this.methodInfoGetMapLayout,
      callback);
  }

  methodInfoListTowers = new grpcWeb.AbstractClientBase.MethodInfo(
    ListTowersResponse,
    (request: ListTowersRequest) => {
      return request.serializeBinary();
    },
    ListTowersResponse.deserializeBinary
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

