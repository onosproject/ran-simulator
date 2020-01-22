
function Car(i) {
    this.num = i;
    this.marker = null;
    this.line = null;
    this.tower = null;
    this.tower1 = null;
    this.tower2 = null;
    this.delta = 0;
}
var cars = [];

var map;

function initMap() {

    map = new google.maps.Map(document.getElementById('map'), {
        center: config.LOC,
        zoom: 12
    });

    initTowers();

    initLocations();

    //Create cars
    for (let i = 0; i < config.NUM_CARS; i++) {
        var car = new Car(i);
        car.marker = new google.maps.Marker({
            map: map,
            //position: config.LOC,
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
        });
        cars.push(car);
    }

    directionsService = new google.maps.DirectionsService();

    //Move cars
    moveIt(0);

    function moveIt(i) {
        moveCarToRoute(
            cars[i],
            cars[i].marker.getPosition(),
            getRandomLocation());

        if (i == cars.length - 1) {
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
        moveMarker(car, path, 0, 0);
    }

    function moveMarker(car, path, pos, delta) {
        var incr = {}
        incr.lat = (path[pos].lat() - car.marker.getPosition().lat()) / config.NUM_DELTAS;
        incr.lng = (path[pos].lng() - car.marker.getPosition().lng()) / config.NUM_DELTAS;
        let radian = Math.atan2(incr.lng, incr.lat);

        car.marker.getIcon().rotation = radians_to_degrees(radian);
        car.marker.getIcon().scale = map.zoom * .025;
        car.marker.setIcon(car.marker.getIcon())
        car.marker.setPosition({
            lat: car.marker.getPosition().lat() + incr.lat,
            lng: car.marker.getPosition().lng() + incr.lng
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
        if (delta != config.NUM_DELTAS) {
            setTimeout(moveMarker, config.UPDATE_DELAY, car, path, pos, ++delta);
        } else {
            if (pos == path.length - 1) {
                // Reached destination
                moveCarToRoute(
                    car,
                    car.marker.getPosition(),
                    getRandomLocation());
                car.line.setMap(undefined);
                car.line.setPath([]);
                return;
            } else {
                // Move car to next position in path
                moveMarker(car, path, pos + 1, 0);
            }
        }
    }
}
