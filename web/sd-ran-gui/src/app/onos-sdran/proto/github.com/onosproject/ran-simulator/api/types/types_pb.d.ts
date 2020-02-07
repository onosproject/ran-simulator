// Code generated by GENERATOR. DO NOT EDIT.

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

  getColor(): string;
  setColor(value: string): void;

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
    color: string,
  }
}

export class Ue extends jspb.Message {
  getName(): string;
  setName(value: string): void;

  getType(): string;
  setType(value: string): void;

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

  getTower2(): string;
  setTower2(value: string): void;

  getTower3(): string;
  setTower3(value: string): void;

  getCrnti(): string;
  setCrnti(value: string): void;

  getTowerdist(): number;
  setTowerdist(value: number): void;

  getTower2dist(): number;
  setTower2dist(value: number): void;

  getTower3dist(): number;
  setTower3dist(value: number): void;

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
    position?: Point.AsObject,
    rotation: number,
    route: string,
    tower: string,
    tower2: string,
    tower3: string,
    crnti: string,
    towerdist: number,
    tower2dist: number,
    tower3dist: number,
  }
}

export class TowersParams extends jspb.Message {
  getTowerrows(): number;
  setTowerrows(value: number): void;

  getTowercols(): number;
  setTowercols(value: number): void;

  getTowerspacingvert(): number;
  setTowerspacingvert(value: number): void;

  getTowerspacinghoriz(): number;
  setTowerspacinghoriz(value: number): void;

  getMaxues(): number;
  setMaxues(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): TowersParams.AsObject;
  static toObject(includeInstance: boolean, msg: TowersParams): TowersParams.AsObject;
  static serializeBinaryToWriter(message: TowersParams, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): TowersParams;
  static deserializeBinaryFromReader(message: TowersParams, reader: jspb.BinaryReader): TowersParams;
}

export namespace TowersParams {
  export type AsObject = {
    towerrows: number,
    towercols: number,
    towerspacingvert: number,
    towerspacinghoriz: number,
    maxues: number,
  }
}

export class Tower extends jspb.Message {
  getName(): string;
  setName(value: string): void;

  getLocation(): Point | undefined;
  setLocation(value?: Point): void;
  hasLocation(): boolean;
  clearLocation(): void;

  getColor(): string;
  setColor(value: string): void;

  getPlmnid(): string;
  setPlmnid(value: string): void;

  getEcid(): string;
  setEcid(value: string): void;

  getMaxues(): number;
  setMaxues(value: number): void;

  getNeighborsList(): Array<string>;
  setNeighborsList(value: Array<string>): void;
  clearNeighborsList(): void;
  addNeighbors(value: string, index?: number): void;

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
    color: string,
    plmnid: string,
    ecid: string,
    maxues: number,
    neighborsList: Array<string>,
  }
}

export class MapLayout extends jspb.Message {
  getCenter(): Point | undefined;
  setCenter(value?: Point): void;
  hasCenter(): boolean;
  clearCenter(): void;

  getZoom(): number;
  setZoom(value: number): void;

  getFade(): boolean;
  setFade(value: boolean): void;

  getShowroutes(): boolean;
  setShowroutes(value: boolean): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): MapLayout.AsObject;
  static toObject(includeInstance: boolean, msg: MapLayout): MapLayout.AsObject;
  static serializeBinaryToWriter(message: MapLayout, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): MapLayout;
  static deserializeBinaryFromReader(message: MapLayout, reader: jspb.BinaryReader): MapLayout;
}

export namespace MapLayout {
  export type AsObject = {
    center?: Point.AsObject,
    zoom: number,
    fade: boolean,
    showroutes: boolean,
  }
}

