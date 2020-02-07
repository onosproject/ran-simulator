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
import {AfterViewInit, Component, OnDestroy, OnInit, ViewChild, ViewChildren} from '@angular/core';
import {GoogleMap, MapInfoWindow, MapMarker} from '@angular/google-maps';
import {Subscription} from 'rxjs';
import {OnosSdranTrafficsimService} from '../proto/onos-sdran-trafficsim.service';
import {Point, Route, Tower, Ue} from '../proto/github.com/onosproject/ran-simulator/api/types/types_pb';
import {ListUesResponse, Type} from '../proto/github.com/onosproject/ran-simulator/api/trafficsim/trafficsim_pb';

export const CAR_ICON = 'M29.395,0H17.636c-3.117,0-5.643,3.467-5.643,6.584v34.804c0,3.116,2.526,5.644,5.643,5.644h11.759 ' +
    'c3.116,0,5.644-2.527,5.644-5.644V6.584C35.037,3.467,32.511,0,29.395,0z M34.05,14.188v11.665l-2.729,0.351v-4.806L34.05,14.188z ' +
    'M32.618,10.773c-1.016,3.9-2.219,8.51-2.219,8.51H16.631l-2.222-8.51C14.41,10.773,23.293,7.755,32.618,10.773z M15.741,21.713 ' +
    'v4.492l-2.73-0.349V14.502L15.741,21.713z M13.011,37.938V27.579l2.73,0.343v8.196L13.011,37.938z M14.568,40.882l2.218-3.336 ' +
    'h13.771l2.219,3.336H14.568z M31.321,35.805v-7.872l2.729-0.355v10.048L31.321,35.805';
export const LINE_SYMBOL = {
    path: 'M 0,-1 0,1',
    strokeOpacity: 1,
    scale: 2
};

export interface Car {
    num: number;
    line: google.maps.Polyline;
    tower: google.maps.Marker;
    tower1: google.maps.Marker;
    tower2: google.maps.Marker;
    delta: number;
    route: google.maps.Polyline;
}

@Component({
    selector: 'onos-mapview',
    templateUrl: './mapview.component.html',
    styleUrls: ['./mapview.component.css']
})
export class MapviewComponent implements OnInit, AfterViewInit, OnDestroy {
    @ViewChild('map', {static: false}) googleMap: GoogleMap;
    @ViewChild(MapInfoWindow, {static: false}) infoWindow: MapInfoWindow;
    infoContent: string[];
    showRoutes = true;
    showMap = false;
    showPower = false;
    zoom = 12.0;
    center: google.maps.LatLng;
    towers: google.maps.Marker[] = [];
    routes: Map<string, google.maps.Polyline>;
    carMap: Map<string, google.maps.Marker>;
    carLineMap: Map<string, google.maps.Polyline>;
    powerCircleMap: Map<string, google.maps.Circle>;
    towerSub: Subscription;
    routesSub: Subscription;
    uesSub: Subscription;

    constructor(
        private trafficSimService: OnosSdranTrafficsimService
    ) {
        this.routes = new Map<string, google.maps.Polyline>();
        this.carMap = new Map<string, google.maps.Marker>();
        this.carLineMap = new Map<string, google.maps.Polyline>();
        this.powerCircleMap = new Map<string, google.maps.Circle>();
    }

    ngOnInit() {
        this.trafficSimService.requestGetMapLayout().subscribe((mapLayout) => {
            this.center = new google.maps.LatLng(mapLayout.getCenter().getLat(), mapLayout.getCenter().getLng());
            this.zoom = mapLayout.getZoom();
            this.showRoutes = mapLayout.getShowroutes();
            this.showMap = !mapLayout.getFade();
        });

        this.towerSub = this.trafficSimService.requestListTowers().subscribe((tower) => {
            this.initTower(tower, this.zoom);
        }, error => {
            console.error('Tower', error);
        }, () => {
            console.log(this.towers.length, 'towers retrieved');
        });

    }

    ngAfterViewInit(): void {
        // Create a Custom Map type to display by default - makes it easier to see cars
        const bwMapStyle = new google.maps.StyledMapType([
            {elementType: 'all', stylers: [{lightness: 70}]}
        ]);

        this.googleMap._googleMap.mapTypes.set('custom', bwMapStyle);
        this.googleMap._googleMap.setMapTypeId('custom');

        // Get the list of routes - we're doing this here because we need to wait until `googleMap` object is populated
        this.routesSub = this.trafficSimService.requestListRoutes().subscribe((resp) => {
            if (resp.getType() === Type.NONE || resp.getType() === Type.ADDED) {
                this.initRoute(resp.getRoute());
            } else if (resp.getType() === Type.UPDATED) {
                this.updateRoute(resp.getRoute());
            } else if (resp.getType() === Type.REMOVED) {
                this.routes.delete(resp.getRoute().getName());
            } else {
                console.warn('Unhandled Route response type', resp.getType(), 'for', resp.getRoute().getName());
            }
        }, error => {
            console.error('Tower', error);
        }, () => {
            console.log(this.routes.size, 'routes retrieved');
        });

        this.uesSub = this.trafficSimService.requestListUes().subscribe((resp: ListUesResponse) => {
            if (resp.getType() === Type.NONE || resp.getType() === Type.ADDED) {
                this.initCar(resp.getUe());
            } else if (resp.getType() === Type.UPDATED) {
                this.updateCar(resp.getUe());
            } else if (resp.getType() === Type.REMOVED) {
                this.carMap.delete(resp.getUe().getName());
                this.carLineMap.delete(resp.getUe().getName());
            } else {
                console.warn('Unhandled Ue response type', resp.getType(), 'for', resp.getUe().getName());
            }
        }, error => {
            console.error('Tower', error);
        }, () => {
            console.log(this.routes.size, 'routes retrieved');
        });
    }

    ngOnDestroy(): void {
        if (this.towerSub !== undefined) {
            this.towerSub.unsubscribe();
        }
        if (this.routesSub !== undefined) {
            this.routesSub.unsubscribe();
        }
        if (this.uesSub !== undefined) {
            this.uesSub.unsubscribe();
        }
    }

    updateRoutes(update: boolean) {
        this.routes.forEach((r) => {
            r.setVisible(update);
        });
    }

    updateMap(update: boolean) {
        this.googleMap._googleMap.setMapTypeId(update ? 'roadmap' : 'custom');
        this.googleMap._googleMap.setOptions({disableDefaultUI: !update} as google.maps.MapOptions);
    }

    updatePower(update: boolean) {
        this.powerCircleMap.forEach((pc) => {
            pc.setVisible(update);
        });
    }

    openTowerInfo(tower: MapMarker) {
        this.infoContent = [tower.getTitle(), 'Lat: ' + tower.getPosition().lat(), 'Lng: ' + tower.getPosition().lng()];
        this.infoWindow.open(tower);
    }

    private initTower(tower: Tower, zoom: number): void {
        const pos = {lat: tower.getLocation().getLat(), lng: tower.getLocation().getLng()};
        const marker = new google.maps.Marker();
        marker.setPosition(pos);
        marker.setTitle(tower.getName());
        marker.setOptions({
                icon: {
                    path: google.maps.SymbolPath.FORWARD_CLOSED_ARROW,
                    scale: zoom * .25,
                    strokeColor: 'blue',
                    strokeWeight: 1,
                    fillColor: tower.getColor(),
                    fillOpacity: 1,
                }
            }
        );
        this.towers.push(marker);

        const powerCircle = new google.maps.Circle({
            center: pos,
            radius: 600,
            fillOpacity: 0,
            strokeColor: tower.getColor(),
            strokeWeight: 0.5,
            strokeOpacity: 1,
            icons: [{
                icon: LINE_SYMBOL,
                offset: '0',
                repeat: '10px',
            }],
        } as google.maps.CircleOptions);
        powerCircle.setMap(this.googleMap._googleMap);
        powerCircle.setVisible(this.showPower);
        this.powerCircleMap.set(tower.getName(), powerCircle);
    }

    private initRoute(route: Route): void {
        const path: google.maps.LatLng[] = [];
        route.getWaypointsList().forEach((point: Point) => {
            const latLng = new google.maps.LatLng(point.getLat(), point.getLng());
            path.push(latLng);
        });
        const polyline = new google.maps.Polyline({
            visible: this.showRoutes,
            strokeWeight: 1,
            strokeOpacity: 0,
            strokeColor: route.getColor(),
            icons: [{
                icon: LINE_SYMBOL,
                offset: '0',
                repeat: '10px',
            }],
        } as google.maps.PolylineOptions);
        polyline.setPath(path);
        polyline.setMap(this.googleMap._googleMap);
        this.routes.set(route.getName(), polyline);
    }

    private updateRoute(route: Route): void {
        console.log('Updated route', route.getName());
        const path: google.maps.LatLng[] = [];
        route.getWaypointsList().forEach((point: Point) => {
            const latLng = new google.maps.LatLng(point.getLat(), point.getLng());
            path.push(latLng);
        });
        this.routes.get(route.getName()).setPath(path);
    }

    private initCar(car: Ue): void {
        const nearestTower = this.towers.find((t) => t.getTitle() === car.getTower());

        const carMarker = new google.maps.Marker({
            icon: {
                path: CAR_ICON,
                scale: this.zoom * .025,
                fillColor: (nearestTower.getIcon() as google.maps.ReadonlySymbol).fillColor,
                anchor: new google.maps.Point(25, 25),
                fillOpacity: 1,
                rotation: 0,
                strokeWeight: 1
            }
        });
        carMarker.setLabel(car.getName());
        carMarker.setTitle(car.getName());
        carMarker.setPosition(
            new google.maps.LatLng(
                car.getPosition().getLat(),
                car.getPosition().getLng()));
        carMarker.setMap(this.googleMap._googleMap);
        this.carMap.set(car.getName(), carMarker);

        // Now need a line from the car to the tower
        const carPolyline = new google.maps.Polyline({
            strokeColor: (nearestTower.getIcon() as google.maps.ReadonlySymbol).fillColor,
            strokeWeight: 1
        } as google.maps.PolylineOptions);
        carPolyline.setPath([carMarker.getPosition(), nearestTower.getPosition()]);
        carPolyline.setMap(this.googleMap._googleMap);
        this.carLineMap.set(car.getName(), carPolyline);
    }

    private updateCar(car: Ue): void {
        const newPos = new google.maps.LatLng(car.getPosition().getLat(), car.getPosition().getLng());
        if (this.carMap.get(car.getName()) !== undefined) {
            this.carMap.get(car.getName()).setPosition(newPos);
            this.carMap.get(car.getName()).get('icon').rotation = car.getRotation();
            const icon: google.maps.Symbol = this.carMap.get(car.getName()).get('icon');
            icon.rotation = 90 - car.getRotation();
            this.carMap.get(car.getName()).setIcon(icon);
        } else {
            console.warn('Car', car.getName(), 'not found in "carMap"');
        }
        const tower = this.towers.find((t) => t.getTitle() === car.getTower());
        if (this.carLineMap.get(car.getName()) !== undefined && tower !== undefined) {
            this.carLineMap.get(car.getName()).setPath([newPos, tower.getPosition()]);
        } else {
            console.warn('Car', car.getName(), 'not found in "carLineMap"');
        }
    }
}
