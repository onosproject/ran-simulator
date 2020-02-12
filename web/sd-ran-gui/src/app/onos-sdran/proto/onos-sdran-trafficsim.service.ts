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

import {Inject, Injectable} from '@angular/core';
import {TrafficClient} from './github.com/onosproject/ran-simulator/api/trafficsim/trafficsimServiceClientPb';
import {Observable, Subscriber} from 'rxjs';
import {MapLayout, Tower} from './github.com/onosproject/ran-simulator/api/types/types_pb';
import {
    ListRoutesRequest, ListRoutesResponse,
    ListTowersRequest, ListTowersResponse, ListUesResponse, MapLayoutRequest
} from './github.com/onosproject/ran-simulator/api/trafficsim/trafficsim_pb';
import * as grpcWeb from 'grpc-web';

@Injectable()
export class OnosSdranTrafficsimService {

    trafficClient: TrafficClient;

    constructor(@Inject('trafficSimUrl') private trafficSimUrl: string) {
        this.trafficClient = new TrafficClient(trafficSimUrl);

        console.log('TrafficSim grpc-web Url ', trafficSimUrl);
    }

    requestGetMapLayout(): Observable<MapLayout> {
        const getMapLayoutObs = new Observable<MapLayout>( (observer: Subscriber<MapLayout>) => {
            const call = this.trafficClient.getMapLayout(new MapLayoutRequest(), {}, ((err, response) => {
                if (err) {
                    observer.error(err);
                } else {
                    observer.next(response);
                }
                call.on('error', (error: grpcWeb.Error) => {
                    observer.error(error);
                });
                call.on('end', () => {
                    observer.complete();
                });
            }));
        });
        return getMapLayoutObs;
    }

    requestListTowers(): Observable<ListTowersResponse> {
        const req = new ListTowersRequest();
        req.setSubscribe(true);
        const stream = this.trafficClient.listTowers(req, {});

        const listTowersObs = new Observable<ListTowersResponse>((observer: Subscriber<ListTowersResponse>) => {
            stream.on('data', (tower: ListTowersResponse) => {
                observer.next(tower);
            });
            stream.on('error', (error: grpcWeb.Error) => {
                observer.error(error);
            });
            stream.on('end', () => {
                observer.complete();
            });
            // stream.on('status', (status: grpcWeb.Status) => {
            //     console.log('ListTowersRequest status', status.code, status.details, status.metadata);
            // });
            return () => stream.cancel();
        });
        return listTowersObs;
    }

    requestListRoutes(): Observable<ListRoutesResponse> {
        const routeReq = new ListRoutesRequest();
        routeReq.setSubscribe(true);
        routeReq.setWithoutreplay(false);
        const stream = this.trafficClient.listRoutes(routeReq, {});

        const listRoutesObs = new Observable<ListRoutesResponse>((observer: Subscriber<ListRoutesResponse>) => {
            stream.on('data', (resp: ListRoutesResponse) => {
                observer.next(resp);
            });
            stream.on('error', (error: grpcWeb.Error) => {
                observer.error(error);
            });
            stream.on('end', () => {
                observer.complete();
            });
            return () => stream.cancel();
        });

        return listRoutesObs;
    }

    requestListUes(): Observable<ListUesResponse> {
        const routeReq = new ListRoutesRequest();
        routeReq.setSubscribe(true);
        routeReq.setWithoutreplay(false);
        const stream = this.trafficClient.listUes(routeReq, {});

        const listUesObs = new Observable<ListUesResponse>((observer: Subscriber<ListUesResponse>) => {
            stream.on('data', (resp: ListUesResponse) => {
                observer.next(resp);
            });
            stream.on('error', (error: grpcWeb.Error) => {
                observer.error(error);
            });
            stream.on('end', () => {
                observer.complete();
            });
            return () => stream.cancel();
        });

        return listUesObs;
    }
}
