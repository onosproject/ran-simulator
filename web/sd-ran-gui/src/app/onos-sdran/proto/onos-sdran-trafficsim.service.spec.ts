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
import {TestBed} from '@angular/core/testing';

import {OnosSdranTrafficsimService} from './onos-sdran-trafficsim.service';

describe('OnosSdranTrafficsimService', () => {
    let service: OnosSdranTrafficsimService;

    beforeEach(() => {
        TestBed.configureTestingModule({
            providers: [
                {
                    provide: OnosSdranTrafficsimService,
                    useValue: new OnosSdranTrafficsimService('http://localhost:8080')
                }
            ]
        });
        service = TestBed.inject(OnosSdranTrafficsimService);
    });

    it('should be created', () => {
        expect(service).toBeTruthy();
    });
});
