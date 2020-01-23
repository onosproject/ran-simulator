
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
    var closestTower = towers[0]
    var closestDistance = distanceTo(car.marker.getPosition(), towers[0].marker.getPosition())
    for (let i = 0; i < towers.length; i++) {
        var distance = distanceTo(car.marker.getPosition(), towers[i].marker.getPosition())
        if (distance < closestDistance) {
            closestDistance = distance
            closestTower = towers[i]
        }
    }
    return closestTower
}
