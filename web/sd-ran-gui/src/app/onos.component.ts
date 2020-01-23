import { Component } from '@angular/core';
import {FnService, IconService, KeysService, LogService} from 'gui2-fw-lib';
import {ConnectivityService} from './connectivity.service';

@Component({
  selector: 'onos-root',
  templateUrl: './onos.component.html',
  styleUrls: ['./onos.component.css']
})
export class OnosComponent {
  title = 'sd-ran-gui';

  constructor(
    protected fs: FnService,
    protected ks: KeysService,
    protected log: LogService,
    protected is: IconService,
    public connectivity: ConnectivityService
  ) {
    this.is.loadIconDef('bird');
    console.log('Constructed OnosComponent');
  }
}
