import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';
import {ConsoleLoggerService, Gui2FwLibModule, LogService} from 'gui2-fw-lib';

import { OnosComponent } from './onos.component';
import {OnosRoutingModule} from './onos-routing.module';
import {NavComponent} from './nav/nav.component';
import {BrowserAnimationsModule} from '@angular/platform-browser/animations';
import {CommonModule} from '@angular/common';
import {HttpClientModule} from '@angular/common/http';

@NgModule({
  declarations: [
    OnosComponent, NavComponent
  ],
  imports: [
    BrowserModule,
    BrowserAnimationsModule,
    CommonModule,
    HttpClientModule,
    Gui2FwLibModule,
    OnosRoutingModule
  ],
  providers: [
    {provide: 'Window', useValue: window},
    {provide: LogService, useClass: ConsoleLoggerService},
  ],
  bootstrap: [OnosComponent]
})
export class OnosModule { }
