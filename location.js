
var locations = []

function initLocations() {
    //Create random location
    for (let i = 0; i < config.NUM_LOCS; i++) {
        locations.push(randomLatLng(config.LOC))
    }
}

function getStartingLocation(carNum) {
            return locations[carNum % locations.length]
}

// Get a random location
function getRandomLocation() {
    return locations[getRandomIntInclusive(0, locations.length - 1)]
}
