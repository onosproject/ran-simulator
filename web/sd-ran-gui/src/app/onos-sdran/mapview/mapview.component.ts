import {
    AfterViewInit,
    Component,
    OnDestroy,
    OnInit,
    ViewChild, ViewChildren
} from '@angular/core';
import {GoogleMap, MapInfoWindow, MapMarker} from '@angular/google-maps';
import {Utils} from '../utils';
import {Observable, Subscriber, Subscription} from 'rxjs';
import {OnosSdranTrafficsimService} from '../proto/onos-sdran-trafficsim.service';
import {trafficSimUrl} from '../../../environments/environment';

export const LOC = {lat: 52.5200, lng: 13.4050} as google.maps.LatLngLiteral; // Ich bin ein Berliner
export const NUM_CARS = 10;
export const NUM_LOCS = 10;
export const NUM_TOWER_ROWS = 3;
export const NUM_TOWER_COLS = 3;
export const NUM_DELTAS = 10;
export const UPDATE_DELAY = 100;
export const CAR_ICON = 'M29.395,0H17.636c-3.117,0-5.643,3.467-5.643,6.584v34.804c0,3.116,2.526,5.644,5.643,5.644h11.759 ' +
    'c3.116,0,5.644-2.527,5.644-5.644V6.584C35.037,3.467,32.511,0,29.395,0z M34.05,14.188v11.665l-2.729,0.351v-4.806L34.05,14.188z ' +
    'M32.618,10.773c-1.016,3.9-2.219,8.51-2.219,8.51H16.631l-2.222-8.51C14.41,10.773,23.293,7.755,32.618,10.773z M15.741,21.713 ' +
    'v4.492l-2.73-0.349V14.502L15.741,21.713z M13.011,37.938V27.579l2.73,0.343v8.196L13.011,37.938z M14.568,40.882l2.218-3.336 ' +
    'h13.771l2.219,3.336H14.568z M31.321,35.805v-7.872l2.729-0.355v10.048L31.321,35.805';

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
    selector: 'app-mapview',
    templateUrl: './mapview.component.html',
    styleUrls: ['./mapview.component.css']
})
export class MapviewComponent implements OnInit, AfterViewInit, OnDestroy {
    @ViewChild('map', {static: false}) googleMap: GoogleMap;
    @ViewChild(MapInfoWindow, { static: false }) infoWindow: MapInfoWindow;
    @ViewChildren('carElem') carElems: MapMarker[];
    infoContent: string[];
    showRoutes = true;
    showMap = false;
    zoom = 12;
    center = LOC;
    towers: google.maps.Marker[] = [];
    locations: google.maps.LatLng[] = [];
    cars: Car[] = [];
    count = 0;
    towerSub: Subscription;

    constructor(
        private directionsService: google.maps.DirectionsService,
        private trafficSimService: OnosSdranTrafficsimService
    ) {
    }

    ngOnInit() {
        this.initTowers(NUM_TOWER_ROWS, NUM_TOWER_COLS, this.zoom);
        this.initLocations(NUM_LOCS);
        this.initCars(NUM_CARS);
    }

    ngAfterViewInit(): void {
        // Create a Custom Map type to display by default - makes it easier to see cars
        const bwMapStyle = new google.maps.StyledMapType([
            {elementType: 'all', stylers: [{lightness: 70}]}
        ]);

        this.googleMap._googleMap.mapTypes.set('custom', bwMapStyle);
        this.googleMap._googleMap.setMapTypeId('custom');

        // Attach the car's line to the map - this can only be done after the map is loaded
        this.cars.forEach((c: Car, i: number) => {
            c.line.setMap(this.googleMap._googleMap);
            let carMarker: MapMarker;
            this.carElems.forEach((m) => {
                if (m.getTitle() === 'Car' + c.num) {
                    carMarker = m;
                }
            });

            carMarker._marker.setPosition(this.getStartingLocation(i));
            carMarker._marker.setOptions({
                icon: {
                    path: CAR_ICON,
                    scale: this.zoom * .025,
                    fillColor: undefined,
                    anchor: new google.maps.Point(25, 25),
                    fillOpacity: 1,
                    strokeWeight: 1
                }
            });
            this.attachCar(c, carMarker);

            this.retrieveRoute(carMarker.getPosition(), this.getRandomLocation(carMarker.getPosition())).subscribe(
                (pos: google.maps.LatLng) => {
                    c.route.getPath().push(pos);
                },
                (err) => {
                    console.error('Could not get a route for car', c.num);
                    c.route.getPath().clear();
                    c.route.setMap(null);
                },
                () => {
                    const numSteps = c.route.getPath().getLength();
                    console.log('Route', c.num, 'has', numSteps, 'steps');
                    c.route.setMap(this.googleMap._googleMap);
                    let position = 1;
                    const timer = setInterval(() => {
                        this.moveCar(c, carMarker, position);
                        position++;
                        if (position === numSteps) {
                            clearInterval(timer);
                        }
                    }, 1000);
                }
            );
        });
    }

    ngOnDestroy(): void {
        this.cars.forEach((c: Car) => {
            c.line.setMap(null);
        });
        this.cars.length = 0;
        this.locations.length = 0;
        this.towers.length = 0;
    }

    updateRoutes(update: boolean) {
        this.cars.forEach((c: Car) => {
            c.route.setVisible(update);
        });

        if (update) {
            console.log('Connecting to', trafficSimUrl);
            this.towerSub = this.trafficSimService.requestListTowers().subscribe((tower) => {
                console.log('Tower', tower);
            }, error => {
                console.error('Tower', error);
            });
        } else {
            if (this.towerSub !== undefined) {
                this.towerSub.unsubscribe();
            }
        }
    }

    updateMap(update: boolean) {
        this.googleMap._googleMap.setMapTypeId(update ? 'roadmap' : 'custom');
        this.googleMap._googleMap.setOptions({disableDefaultUI: !update} as google.maps.MapOptions);
    }

    openCarInfo(marker: MapMarker, car: Car) {
        this.infoContent = ['Car' + car.num, car.tower.getTitle()];
        this.infoWindow.open(marker);
    }

    openTowerInfo(tower: MapMarker) {
        this.infoContent = [tower.getTitle(), 'Lat: ' + tower.getPosition().lat(), 'Lng: ' + tower.getPosition().lng()];
        this.infoWindow.open(tower);
    }

    private initTowers(rows: number, cols: number, zoom: number): void {
        const topLeft = {lat: LOC.lat + 0.02 * rows / 2, lng: LOC.lng - 0.03 * cols / 2};
        let towerNum = 0;

        for (let r = 0; r < rows; r++) {
            for (let c = 0; c < cols; c++, towerNum++) {
                const color = Utils.getRandomColor();
                const pos = {lat: topLeft.lat - 0.03 * r, lng: topLeft.lng + 0.05 * c};
                const marker = new google.maps.Marker();
                marker.setPosition(pos);
                marker.setTitle('Tower' + towerNum);
                marker.setOptions({
                    icon: {
                        path: google.maps.SymbolPath.FORWARD_CLOSED_ARROW,
                        scale: zoom * .25,
                        strokeColor: 'blue',
                        strokeWeight: 1,
                        fillColor: color,
                        fillOpacity: 1,
                    }}
                );

                this.towers.push(marker);
            }
        }
    }

    private initLocations(numLocations: number): void {
        for (let i = 0; i < numLocations; i++) {
            const loc = Utils.randomLatLng(LOC);
            this.locations.push(loc);
        }
    }

    private initCars(numCars: number): void {
        for (let i = 0; i < numCars; i++) {
            const car = {num: i} as Car;
            car.line = new google.maps.Polyline({
                strokeWeight: 1
            });
            const lineSymbol = {
                path: 'M 0,-1 0,1',
                strokeOpacity: 1,
                scale: 2
            };
            car.route = new google.maps.Polyline({
                visible: this.showRoutes,
                strokeWeight: 1,
                strokeOpacity: 0,
                icons: [{
                    icon: lineSymbol,
                    offset: '0',
                    repeat: '10px'
                }],
            });

            this.cars.push(car);
        }
    }

    private getStartingLocation(carNum: number): google.maps.LatLng {
        return this.locations[carNum % this.locations.length];
    }

    // Get a random location, whose position is other than 'exclude'
    private getRandomLocation(exclude?: google.maps.LatLng): google.maps.LatLng {
        // tslint:disable-next-line:prefer-for-of
        for (let i = 0; i < this.locations.length; i++) {
            const idx = Utils.getRandomIntInclusive(0, this.locations.length - 1);
            if (this.locations[idx].lat() === exclude.lat() && this.locations[idx].lng() === exclude.lng()) {
                continue;
            }
            return this.locations[idx];
        }

        // If we get to here it means there's only 1 location - take it
        return this.locations[0];
    }

    private attachCar(car: Car, carMarker: MapMarker): void {
        car.tower = this.findClosestTower(carMarker);
        car.line.setOptions({
            path: [carMarker.getPosition(), car.tower.getPosition()],
            strokeColor: (car.tower.getIcon() as google.maps.ReadonlySymbol).fillColor
        } as google.maps.PolylineOptions);
        const towerColor = (car.tower.getIcon() as google.maps.ReadonlySymbol).fillColor;
        const carIcon = carMarker._marker.getIcon() as google.maps.Symbol;
        carIcon.fillColor = towerColor;
        carMarker._marker.setIcon(carIcon);
        car.route.set('strokeColor', towerColor);
    }

    private findClosestTower(carMarker: MapMarker): google.maps.Marker {
        let closestTower = this.towers[0];
        let closestDistance = Utils.distanceTo(carMarker.getPosition(), closestTower.getPosition());
        let closestTowerIdx = 0;
        this.towers.forEach((t, i) => {
            const distance = Utils.distanceTo(carMarker.getPosition(), t.getPosition());
            if (distance < closestDistance) {
                closestDistance = distance;
                closestTower = t;
                closestTowerIdx = i;
            }
        });
        return closestTower;
    }

    // Move car in a delta increment along path towards pos(ition)
    private moveCar(car: Car, carMarker: MapMarker, pos: number, delta?: number): void {
        const currLat = carMarker.getPosition().lat();
        const currLng = carMarker.getPosition().lng();
        const latIncr = (car.route.getPath().getAt(pos).lat() - currLat);
        const lngIncr = (car.route.getPath().getAt(pos).lng() - currLng);
        const newLat = currLat + latIncr;
        const newLng = currLng + lngIncr;

        const radian = Math.atan2(lngIncr, latIncr);
        carMarker._marker.set('rotation', Utils.radians_to_degrees(radian));
        carMarker._marker.set('scale', this.zoom * .025);
        const newPos = new google.maps.LatLng(newLat, newLng);
        carMarker._marker.setPosition(new google.maps.LatLng(newLat, newLng));

        car.line.setPath([car.tower.getPosition(), newPos]);

        // if (delta !== NUM_DELTAS) {
        //     setTimeout(this.moveCar, UPDATE_DELAY, car, path, pos, ++delta);
        // } else {
        //     this.attachCar(car);
        //     if (pos === path.length - 1) {
        //         // Reached destination - Start a new journey
        //         this.moveCarToRoute(
        //             car,
        //             car.marker.getPosition(),
        //             this.getRandomLocation());
        //     } else {
        //         // Move car to next position in path
        //         this.moveCar(car, path, pos + 1, 0);
        //     }
        // }
    }

    private retrieveRoute(start: google.maps.LatLng, end: google.maps.LatLng): Observable<google.maps.LatLng> {
        const request = {
            origin: start,
            destination: end,
            travelMode: 'DRIVING'
        } as google.maps.DirectionsRequest;
        const routeObs = new Observable<google.maps.LatLng>((observer: Subscriber<google.maps.LatLng>) => {
            this.directionsService.route(request, (result, status) => {
                if (status !== 'OK') {
                    observer.error(new Error('Error getting directions' + status));
                }
                result.routes[0].overview_path.forEach((pos: google.maps.LatLng) => {
                    observer.next(pos);
                });
                observer.complete();
            });
        });
        return routeObs;
    }
}
