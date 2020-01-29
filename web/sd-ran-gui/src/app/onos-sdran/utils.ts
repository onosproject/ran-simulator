
export class Utils {
    static getRandomColor(): string {
        const letters = '0123456789ABCDEF';
        let color = '#';
        for (let i = 0; i < 6; i++) {
            color += letters[Math.floor(Math.random() * 16)];
        }
        return color;
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
