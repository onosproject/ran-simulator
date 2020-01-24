
export class Utils {
    static getRandomColor(): string {
        const letters = '0123456789ABCDEF';
        let color = '#';
        for (let i = 0; i < 6; i++) {
            color += letters[Math.floor(Math.random() * 16)];
        }
        return color;
    }

    /**
     * Generates a random latlng value in 1000 meter radius of loc
     */
    static randomLatLng(loc): google.maps.LatLng {
        // var r = 10000 / 111300 // = 100 meters
        const r = 5000 / 111300; // = 100 meters
        const y0 = loc.lat;
        const x0 = loc.lng;
        const u = Math.random();
        const v = Math.random();
        const w = r * Math.sqrt(u);
        const t = 2 * Math.PI * v;
        const x = w * Math.cos(t);
        const y1 = w * Math.sin(t);
        const x1 = x / Math.cos(y0);

        const newY = Utils.round(y0 + y1, 6);
        const newX = Utils.round(x0 + x1, 6);

        return new google.maps.LatLng(newY, newX);
    }

    /**
     * Rounds number to decimals
     */
    private static round(value: number, decimals: number): number {
        const intValue = value * Math.pow(10, decimals);
        return Number(Math.round(intValue) / Math.pow(10, decimals));
    }

    static radians_to_degrees(radians: number): number {
        return radians * (180 / Math.PI);
    }

    static distanceTo(pos1: google.maps.LatLng, pos2: google.maps.LatLng): number {
        const dX: number = pos1.lat() - pos2.lat();
        const dY: number = pos1.lng() - pos2.lng();
        return Math.sqrt(Math.pow(dX, 2) + Math.pow(dY, 2));
    }

    static getRandomIntInclusive(min: number, max: number): number {
        min = Math.ceil(min);
        max = Math.floor(max);
        return Math.floor(Math.random() * (max - min + 1)) + min; // The maximum is inclusive and the minimum is inclusive
    }
}
