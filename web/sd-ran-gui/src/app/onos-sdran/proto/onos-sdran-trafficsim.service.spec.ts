import { TestBed } from '@angular/core/testing';

import { OnosSdranTrafficsimService } from './onos-sdran-trafficsim.service';

describe('OnosSdranTrafficsimService', () => {
  let service: OnosSdranTrafficsimService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(OnosSdranTrafficsimService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
