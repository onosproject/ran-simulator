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
import {TestBed, async} from '@angular/core/testing';
import {OnosComponent} from './onos.component';
import {RouterTestingModule} from '@angular/router/testing';
import {FnService, LogService} from 'gui2-fw-lib';
import {ActivatedRoute, Params} from '@angular/router';
import {from, of} from 'rxjs';
import {HttpClient} from '@angular/common/http';

class MockActivatedRoute extends ActivatedRoute {
    constructor(params: Params) {
        super();
        this.queryParams = of(params);
    }
}

class MockHttpClient {
    get() {
        return from(['{"id":"app","icon":"nav_apps","cat":"PLATFORM","label":"Applications"},' +
        '{"id":"settings","icon":"nav_settings","cat":"PLATFORM","label":"Settings"}']);
    }

    subscribe() {}
}

describe('OnosComponent', () => {
    let fs: FnService;
    let ar: MockActivatedRoute;
    let windowMock: Window;

    let logServiceSpy: jasmine.SpyObj<LogService>;

    beforeEach(async(() => {
        const logSpy = jasmine.createSpyObj('LogService', ['info', 'debug', 'warn', 'error']);
        ar = new MockActivatedRoute({ debug: 'txrx' });

        windowMock = {
            location: {
                hostname: 'foo',
                host: 'foo',
                port: '80',
                protocol: 'http',
                search: { debug: 'true' },
                href: 'ws://foo:123/onos/ui/websock/path',
                absUrl: 'ws://foo:123/onos/ui/websock/path'
            } as any
        } as any;
        fs = new FnService(ar, logSpy, windowMock);

        TestBed.configureTestingModule({
            imports: [
                RouterTestingModule
            ],
            declarations: [
                OnosComponent
            ],
            providers: [
                { provide: 'Window', useValue: windowMock },
                { provide: LogService, useValue: logSpy },
                { provide: HttpClient, useClass: MockHttpClient },
            ]
        }).compileComponents();
        logServiceSpy = TestBed.get(LogService);
    }));

    it('should create the app', () => {
        const fixture = TestBed.createComponent(OnosComponent);
        const app = fixture.componentInstance;
        expect(app).toBeTruthy();
    });

    it(`should have as title 'sd-ran-gui'`, () => {
        const fixture = TestBed.createComponent(OnosComponent);
        const app = fixture.componentInstance;
        expect(app.title).toEqual('sd-ran-gui');
    });
});
