import {NgModule} from '@angular/core';
import {CommonModule} from '@angular/common';
import {MapviewComponent} from './mapview/mapview.component';
import {RouterModule} from '@angular/router';
import {GoogleMapsModule} from '@angular/google-maps';
import {FormsModule} from '@angular/forms';

@NgModule({
    declarations: [MapviewComponent],
    imports: [
        CommonModule,
        GoogleMapsModule,
        FormsModule,
        RouterModule.forChild([{path: '', component: MapviewComponent}]),
    ],
    providers: [
        google.maps.DirectionsService
    ]
})
export class OnosSdranModule {
}
