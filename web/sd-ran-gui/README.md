# SdRanGui

This project was generated with [Angular CLI] version 9.0.0-rc.7.

It is based on the same foundation at the `onos-gui` of the µONOS project
([see](https://docs.onosproject.org/onos-gui/docs/config-gui/)).

This foundation was just copied over from ONOS GUI, and can be discarded when the 2 are merged together.

The main application happens inside the `web/sd-ran-gui/src/app/onos-sdran` folder.
It uses the new [Google Maps Angular Component]

It uses [grpc-web] on the [Envoy Proxy] to communicate to the `trafficsim` backend

Use [docker-compose] to run the `sd-ran-gui`, `envoy-proxy` and `trafficsim` together. See [README.md](../../README.md)

The original JavaScript only application is in the `web/sd-ran-gui/originalJs` directory

To see how to install NodeJS and Angular CLI on your system follow the
[Development Prerequisites](https://docs.onosproject.org/onos-gui/docs/prerequisites/)
of the ONOS GUI project.

## Development server

Run `ng serve` for a dev server. Navigate to `http://localhost:4200/`. The app will automatically
reload if you change any of the source files.

## Code scaffolding

Run `ng generate component component-name` to generate a new component. You can also use
`ng generate directive|pipe|service|class|guard|interface|enum|module`.

## Build

Run `ng build` to build the project. The build artifacts will be stored in the `dist/` directory.
Use the `--prod` flag for a production build.

## Running unit tests

Run `ng test` to execute the unit tests via [Karma].

## Running end-to-end tests

Run `ng e2e` to execute the end-to-end tests via [Protractor]).

## Further help

To get more help on the Angular CLI use `ng help` or go check out the [Angular CLI README].

[Google Maps Angular Component]: https://medium.com/angular-in-depth/google-maps-is-now-an-angular-component-821ec61d2a0
[docker-compose]: https://docs.docker.com/compose/
[grpc-web]: https://github.com/grpc/grpc-web
[Protractor]: http://www.protractortest.org/
[Karma]: https://karma-runner.github.io
[Envoy Proxy]: https://www.envoyproxy.io/
[Angular CLI]: https://github.com/angular/angular-cli
[Angular CLI README]: https://github.com/angular/angular-cli/blob/master/README.md
