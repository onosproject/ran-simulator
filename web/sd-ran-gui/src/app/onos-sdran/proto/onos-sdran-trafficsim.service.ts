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
import {TrafficClient} from './github.com/OpenNetworkingFoundation/gmap-ran/api/trafficsim/trafficsimServiceClientPb';
import {Observable, Subscriber} from 'rxjs';
import {Tower} from './github.com/OpenNetworkingFoundation/gmap-ran/api/types/types_pb';
import {ListTowersRequest} from './github.com/OpenNetworkingFoundation/gmap-ran/api/trafficsim/trafficsim_pb';
import * as grpcWeb from 'grpc-web';

@Injectable()
export class OnosSdranTrafficsimService {

    trafficClient: TrafficClient;

    constructor(@Inject('trafficSimUrl') private trafficSimUrl: string) {
        this.trafficClient = new TrafficClient(trafficSimUrl);

        console.log('Config Admin Url ', trafficSimUrl);
    }

    requestListTowers(): Observable<Tower> {
        console.log('Calling on', this.trafficSimUrl, ' to list towers');
        const stream = this.trafficClient.listTowers(new ListTowersRequest());

        const listTowersObs = new Observable<Tower>((observer: Subscriber<Tower>) => {
            stream.on('data', (tower: Tower) => {
                observer.next(tower);
            });
            stream.on('error', (error: grpcWeb.Error) => {
                observer.error(error);
            });
            stream.on('end', () => {
                observer.complete();
            });
            stream.on('status', (status: grpcWeb.Status) => {
                console.log('ListTowersRequest status', status.code, status.details, status.metadata);
            });
            return stream.cancel();
        });
        return listTowersObs;
    }
}
