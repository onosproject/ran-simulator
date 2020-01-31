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
var map

function initMap() {

    map = new google.maps.Map(document.getElementById('map'), {
        center: config.LOC,
        zoom: 12
    })

    initTowers()

    initLocations()

    initCars()

    directionsService = new google.maps.DirectionsService()

    //Move cars
    moveIt(0)

    function moveIt(i) {
        moveCarToRoute(
            cars[i],
            cars[i].marker.getPosition(),
            getRandomLocation())

        if (i == cars.length - 1) {
            return
        } else {
            setTimeout(moveIt, 1000, ++i)
        }
    }

    function moveCarToRoute(car, start, end) {
        var request = {
            origin: start,
            destination: end,
            travelMode: 'DRIVING'
        }
        directionsService.route(request, function(result, status) {
            if (status == 'OK') {
                var path = result.routes[0].overview_path
                car.marker.setPosition(path[0])
                moveCar(car, path, 0, 0)
            }
        })
    }

    function moveCar(car, path, pos, delta) {

        // Move car in a delta increment along path towards pos(ition)
        incrMoveCar(car, path, pos)

        attachCar(car)

        if (delta != config.NUM_DELTAS) {
            setTimeout(moveCar, config.UPDATE_DELAY, car, path, pos, ++delta)
        } else {
            if (pos == path.length - 1) {
                // Reached destination - Start a new journey
                moveCarToRoute(
                    car,
                    car.marker.getPosition(),
                    getRandomLocation())
            } else {
                // Move car to next position in path
                moveCar(car, path, pos + 1, 0)
            }
        }
    }
}
