

function Tower() {
    this.marker = null;
}
var towers = [];

function Car(i) {
    this.num = i;
    this.marker = null;
    this.line = null;
    this.tower = null;
    this.tower1 = null;
    this.tower2 = null;
    let x = locations[i % locations.length];
    this.pos = {}
    this.pos['lat'] = x.lat;
    this.pos['lng'] = x.lng;
    this.delta = 0;
}
var cars = [];

var map;
var locations = [];

function initMap() {

    map = new google.maps.Map(document.getElementById('map'), {
        center: config.LOC,
        zoom: 11
    });

    //Create random location
    for (let i = 0; i < config.NUM_LOCS; i++) {
        locations.push(randomLatLng(config.LOC));
    }

    //Create cars
    for (let i = 0; i < locations.length; i++) {
        var car = new Car(i);
        car.marker = new google.maps.Marker({
            map: map,
            //position: config.LOC,
            position: car.pos,
            title: "Car" + car.num,
            icon: {
                path: config.CAR_ICON,
                scale: map.zoom * .03,
                fillColor: undefined,
                anchor: new google.maps.Point(25, 25),
                fillOpacity: 1,
                strokeWeight: 1
            }
        });
        cars.push(car);
    }

    //Create towers
    for (let i = 0; i < config.NUM_TOWERS; i++) {
        var tower = new Tower();
        let color = '#' + (Math.random() * 0xFFFFFF << 0).toString(16);
        tower.marker = (new google.maps.Marker({
            map: map,
            position: randomLatLng(config.LOC),
            title: "Tower",
            icon: {
                path: google.maps.SymbolPath.FORWARD_CLOSED_ARROW,
                scale: 4,
                strokeColor: color,
                fillColor: color,
                fillOpacity: 1,
            }
        }));
        towers.push(tower);
    }

    directionsService = new google.maps.DirectionsService();

    //Move cars
    moveIt(0);

    function moveIt(i) {
        moveCarToRoute(
            cars[i],
            locations[locations.length - i - 1],
            locations[i]);

        if (i == config.NUM_CARS - 1) {
            return;
        } else {
            setTimeout(moveIt, 1000, ++i)
        }
    }

    function moveCarToRoute(car, start, end) {
        var request = {
            origin: start,
            destination: end,
            travelMode: 'DRIVING'
        };
        directionsService.route(request, function(result, status) {
            if (status == 'OK') {
                var path = result.routes[0].overview_path;
                car.marker.setPosition(path[0]);
                moveCar(car, path);
            }
        });
    }

    function moveCar(car, path) {
        car.tower = findClosestTower(car);
        car.line = new google.maps.Polyline({
            map: map,
            path: [car.marker.getPosition(), car.tower.marker.getPosition()],
            strokeWeight: 1})
        transition(car, path, 0);
    }

    function transition(car, path, pos) {
        car.delta = 0;
        car.pos.lat = (path[pos].lat() - car.marker.getPosition().lat()) / config.NUM_DELTAS;
        car.pos.lng = (path[pos].lng() - car.marker.getPosition().lng()) / config.NUM_DELTAS;
        let radian = Math.atan2(car.pos.lng, car.pos.lat);

        car.marker.getIcon().rotation = radians_to_degrees(radian);
        car.marker.getIcon().scale = map.zoom * .035;
        car.marker.setIcon(car.marker.getIcon())

        moveMarker(car, path, pos);
    }

    function moveMarker(car, path, pos) {
        car.marker.setPosition({
            lat: car.marker.getPosition().lat() + car.pos.lat,
            lng: car.marker.getPosition().lng() + car.pos.lng
        });
        car.tower = findClosestTower(car);
        car.line.setOptions({
            map: map,
            path: [car.marker.getPosition(), car.tower.marker.getPosition()],
            strokeWeight: 1,
            strokeColor: car.tower.marker.getIcon().strokeColor
        })
        car.marker.getIcon().fillColor = car.tower.marker.getIcon().strokeColor;
        car.marker.setIcon(car.marker.getIcon());
        if (car.delta != config.NUM_DELTAS) {
            car.delta++;
            setTimeout(moveMarker, config.UPDATE_DELAY, car, path, pos);
        } else {
            if (pos == path.length - 1) {
                // Reached destination
                moveCarToRoute(
                    car,
                    car.marker.getPosition(),
                    locations[getRandomIntInclusive(0, locations.length - 1)]);
                car.line.setMap(undefined);
                car.line.setPath([]);
                return;
            } else {
                // Move car to next position in path
                transition(car, path, pos + 1);
            }
        }
    }

  function findClosestTower(car) {
    var closestTower = towers[0];
      var closestDistance = distanceTo(car.marker.getPosition(), towers[0].marker.getPosition());
    for (let i = 0; i < towers.length; i++) {
      var distance = distanceTo(car.marker.getPosition(), towers[i].marker.getPosition());
      if (distance < closestDistance) {
        closestDistance = distance;
        closestTower = towers[i];
      }
    }
    return closestTower;
  }
}
