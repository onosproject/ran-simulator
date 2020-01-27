import {NgModule} from '@angular/core';
import {CommonModule} from '@angular/common';
import {MapviewComponent} from './mapview/mapview.component';
import {RouterModule} from '@angular/router';
import {GoogleMapsModule} from '@angular/google-maps';
import {FormsModule} from '@angular/forms';
import {OnosSdranTrafficsimService} from './proto/onos-sdran-trafficsim.service';
import {trafficSimUrl} from '../../environments/environment';

@NgModule({
    declarations: [MapviewComponent],
    imports: [
        CommonModule,
        GoogleMapsModule,
        FormsModule,
        RouterModule.forChild([{path: '', component: MapviewComponent}]),
    ],
    providers: [
        {
            provide: google.maps.DirectionsService,
            useClass: google.maps.DirectionsService
        },
        {
            provide: OnosSdranTrafficsimService,
            useValue: new OnosSdranTrafficsimService(trafficSimUrl)
        }
    ]
})
export class OnosSdranModule {
}
