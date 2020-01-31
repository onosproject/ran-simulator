/*
 * Copyright 2018-present Open Networking Foundation
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
import {NgModule} from '@angular/core';
import {Routes, RouterModule, Route} from '@angular/router';

/**
 * The set of Routes in the application - can be chosen from nav menu or
 * elsewhere like tabular icon for flows etc
 */
const onosRoutes: Routes = [
    {
        path: 'mapview',
        pathMatch: 'full',
        loadChildren: () => import('./onos-sdran/onos-sdran.module').then(m => m.OnosSdranModule)
    } as Route,
    {
        path: '',
        redirectTo: 'mapview', // Default to Topology view
        pathMatch: 'full'
    } as Route
];

/**
 * ONOS GUI -- Main Routing Module - allows modules to be lazy loaded
 *
 * See https://angular.io/guide/lazy-loading-ngmodules
 * for the theory of operation
 */
@NgModule({
    imports: [
        RouterModule.forRoot(onosRoutes, {useHash: true})
    ],
    exports: [RouterModule],
    providers: []
})
export class OnosRoutingModule {
}
