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
import { TestBed, async } from '@angular/core/testing';
import { OnosComponent } from './onos.component';

describe('OnosComponent', () => {
  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [
        OnosComponent
      ],
    }).compileComponents();
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

  it('should render title', () => {
    const fixture = TestBed.createComponent(OnosComponent);
    fixture.detectChanges();
    const compiled = fixture.nativeElement;
    expect(compiled.querySelector('.content span').textContent).toContain('sd-ran-gui app is running!');
  });
});
