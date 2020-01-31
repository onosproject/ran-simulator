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
var towers = []

function Tower() {
    this.marker = null
}

function initTowers() {
    var topLeft = {lat: config.LOC.lat + 0.02*config.NUM_TOWER_ROWS/2, lng: config.LOC.lng - 0.03*config.NUM_TOWER_COLS/2}
    var tower_num = 0

    for (let r = 0; r < config.NUM_TOWER_ROWS; r++) {
        for (let c = 0; c < config.NUM_TOWER_COLS; c++, tower_num++) {
            var tower = new Tower()
            let color = getRandomColor()
            let pos = {lat: topLeft.lat - 0.03*r, lng: topLeft.lng + 0.05*c}
            tower.marker = (new google.maps.Marker({
                map: map,
                position: pos,
                title: "Tower" + tower_num,
                icon: {
                    path: google.maps.SymbolPath.FORWARD_CLOSED_ARROW,
                    scale: map.zoom * .25,
                    strokeColor: color,
                    fillColor: color,
                    fillOpacity: 1,
                }
            }))
            towers.push(tower)
        }
    }

    return towers
}

function findClosestTower(car) {
    var serving
    var candidate1
    var candidate2
    var servingDist = Number.MAX_SAFE_INTEGER
    var candidate1Dist = Number.MAX_SAFE_INTEGER
    var candidate2Dist = Number.MAX_SAFE_INTEGER

    for (let i = 0; i < towers.length; i++) {
        var distance = distanceTo(car.marker.getPosition(), towers[i].marker.getPosition())
        if (distance < servingDist) {
            candidate2 = candidate1
            candidate2Dist = candidate1Dist
            candidate1 = serving
            candidate1Dist = servingDist
            serving = towers[i]
            servingDist = distance
        } else if (distance < candidate1Dist) {
            candidate2 = candidate1
            candidate2Dist = candidate1Dist
            candidate1 = towers[i]
            candidate1Dist = distance
        } else if (distance < candidate1Dist) {
            candidate2 = towers[i]
            candidate2Dist = distance
        }
    }
    return {
        serving: serving,
        candidate1: candidate1,
        candidate2: candidate2,
    }
}
