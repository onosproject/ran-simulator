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
function Car(i) {
    this.num = i
    this.marker = null
    this.line = null
    this.tower = {}
    this.delta = 0
}
var cars = []

function initCars() {
    for (let i = 0; i < config.NUM_CARS; i++) {
        var car = new Car(i)
        car.marker = new google.maps.Marker({
            map: map,
            position: getStartingLocation(i),
            title: "Car" + car.num,
            icon: {
                path: config.CAR_ICON,
                scale: map.zoom * .025,
                fillColor: undefined,
                anchor: new google.maps.Point(25, 25),
                fillOpacity: 1,
                strokeWeight: 1
            }
        })
        car.line = new google.maps.Polyline({
            map: map,
            strokeWeight: 1})

        attachCar(car)

        cars.push(car)
    }
}

function incrMoveCar(car, path, pos) {
    var incr = {}
    incr.lat = (path[pos].lat() - car.marker.getPosition().lat()) / config.NUM_DELTAS
    incr.lng = (path[pos].lng() - car.marker.getPosition().lng()) / config.NUM_DELTAS
    let radian = Math.atan2(incr.lng, incr.lat)

    car.marker.getIcon().rotation = radians_to_degrees(radian)
    car.marker.getIcon().scale = map.zoom * .025
    car.marker.setIcon(car.marker.getIcon())
    car.marker.setPosition({
        lat: car.marker.getPosition().lat() + incr.lat,
        lng: car.marker.getPosition().lng() + incr.lng
    })
}

function attachCar(car) {
    car.tower = findClosestTower(car)
    car.line.setOptions({
        path: [car.marker.getPosition(), car.tower.serving.marker.getPosition()],
        strokeColor: car.tower.serving.marker.getIcon().strokeColor
    })
    car.marker.getIcon().fillColor = car.tower.serving.marker.getIcon().strokeColor
}
