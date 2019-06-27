# `dist-zip-cnb`
The Cloud Foundry DistZip Buildpack is a Cloud Native Buildpack V3 that provides enables the running of DistZip style applications.

## Detection
The detection phase passes if:

* The build plan contains `jvm-application`

## Build
If the build plan contains

* `jvm-application`
  * Checks for the existence of a single non-Windows start script in `<APPLICATION_ROOT>/*/bin/*`
  * If found,
    * Contributes suitably configured process types

## License
This buildpack is released under version 2.0 of the [Apache License][a].

[a]: https://www.apache.org/licenses/LICENSE-2.0
