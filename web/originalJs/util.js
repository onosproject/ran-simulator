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
function getColor() {
    return "hsl(" + 360 * Math.random() + ',' + (100 * Math.random()) + '%,' + (80 * Math.random()) + '%)'
}

/**
* Rounds number to decimals
*/
function round(value, decimals) {
    return Number(Math.round(value + 'e' + decimals) + 'e-' + decimals);
}

function radians_to_degrees(radians) {
    var pi = Math.PI;
    return radians * (180 / pi);
}

/**
*Generates a random latlng value in 1000 meter radius of loc
*/
function randomLatLng(loc) {
    //var r = 10000 / 111300 // = 100 meters
    var r = 5000 / 111300 // = 100 meters
    var y0 = loc.lat;
    var x0 = loc.lng;
    var u = Math.random();
    var v = Math.random();
    var w = r * Math.sqrt(u);
    var t = 2 * Math.PI * v;
    var x = w * Math.cos(t);
    var y1 = w * Math.sin(t);
    var x1 = x / Math.cos(y0);

    var newY = round(y0 + y1, 6);
    var newX = round(x0 + x1, 6);

    return {
        lat: newY,
        lng: newX
    };
}

function distanceTo(pos1, pos2) {
    var dX = pos1.lat() - pos2.lat();
    var dY = pos1.lng() - pos2.lng();
    return Math.sqrt(Math.pow(dX, 2) + Math.pow(dY, 2));
}

function getRandomIntInclusive(min, max) {
  min = Math.ceil(min);
  max = Math.floor(max);
  return Math.floor(Math.random() * (max - min + 1)) + min; //The maximum is inclusive and the minimum is inclusive
}

function getRandomColor() {
  var letters = '0123456789ABCDEF';
  var color = '#';
  for (var i = 0; i < 6; i++) {
    color += letters[Math.floor(Math.random() * 16)];
  }
  return color;
}
