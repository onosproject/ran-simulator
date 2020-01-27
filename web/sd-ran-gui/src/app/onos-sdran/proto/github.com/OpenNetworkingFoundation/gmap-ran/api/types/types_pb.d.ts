import * as jspb from "google-protobuf"

export class Point extends jspb.Message {
  getLat(): number;
  setLat(value: number): void;

  getLng(): number;
  setLng(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Point.AsObject;
  static toObject(includeInstance: boolean, msg: Point): Point.AsObject;
  static serializeBinaryToWriter(message: Point, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Point;
  static deserializeBinaryFromReader(message: Point, reader: jspb.BinaryReader): Point;
}

export namespace Point {
  export type AsObject = {
    lat: number,
    lng: number,
  }
}

export class Route extends jspb.Message {
  getName(): string;
  setName(value: string): void;

  getWaypointsList(): Array<Point>;
  setWaypointsList(value: Array<Point>): void;
  clearWaypointsList(): void;
  addWaypoints(value?: Point, index?: number): Point;

  getLengthm(): number;
  setLengthm(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Route.AsObject;
  static toObject(includeInstance: boolean, msg: Route): Route.AsObject;
  static serializeBinaryToWriter(message: Route, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Route;
  static deserializeBinaryFromReader(message: Route, reader: jspb.BinaryReader): Route;
}

export namespace Route {
  export type AsObject = {
    name: string,
    waypointsList: Array<Point.AsObject>,
    lengthm: number,
  }
}

export class Ue extends jspb.Message {
  getName(): string;
  setName(value: string): void;

  getType(): string;
  setType(value: string): void;

  getColor(): string;
  setColor(value: string): void;

  getPosition(): Point | undefined;
  setPosition(value?: Point): void;
  hasPosition(): boolean;
  clearPosition(): void;

  getRotation(): number;
  setRotation(value: number): void;

  getRoute(): string;
  setRoute(value: string): void;

  getTower(): string;
  setTower(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Ue.AsObject;
  static toObject(includeInstance: boolean, msg: Ue): Ue.AsObject;
  static serializeBinaryToWriter(message: Ue, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Ue;
  static deserializeBinaryFromReader(message: Ue, reader: jspb.BinaryReader): Ue;
}

export namespace Ue {
  export type AsObject = {
    name: string,
    type: string,
    color: string,
    position?: Point.AsObject,
    rotation: number,
    route: string,
    tower: string,
  }
}

export class Tower extends jspb.Message {
  getName(): string;
  setName(value: string): void;

  getLocation(): Point | undefined;
  setLocation(value?: Point): void;
  hasLocation(): boolean;
  clearLocation(): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Tower.AsObject;
  static toObject(includeInstance: boolean, msg: Tower): Tower.AsObject;
  static serializeBinaryToWriter(message: Tower, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Tower;
  static deserializeBinaryFromReader(message: Tower, reader: jspb.BinaryReader): Tower;
}

export namespace Tower {
  export type AsObject = {
    name: string,
    location?: Point.AsObject,
  }
}

