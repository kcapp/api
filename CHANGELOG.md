# Changelog

## [1.2.0] - 2020-10-10
#### Feature
- Support for new game types `Tic-Tac-Toe`, `Bermuda Triangle`, and `420`
- Global inidicator for Offices

#### Changed
- Each statistics type contain `office_id`

## [1.1.0] - 2020-07-18
#### Feature
- Support for new game types `Around The World`, `Around The Clock` and `Shanghai`
- Added match statistics for `9 Dart Shootout`
- New convenience endpoints`/statistics/<type_id>/<from>/<to>`
- More statistics to global statistics endpoint
- New endpoints for getting player statistics and player history

#### Changed
- Writing of `9 Dart Shootout` score to database

#### Fixed
- Correctly calculate legs played and won for different statistics
- Fixed calculating of matches and legs played and won for shootout
- Synchronized `AddVisit` function to prevent multiple entries of same score
- Graceful handing of matches with venue id 0

## [1.0.0] - 2020-05-03
#### Feature
- Intial version of API for [kcapp-frontend](https://github.com/kcapp/frontend)

[1.2.0]: https://github.com/kcapp/api/compare/v1.1.0...v1.2.0
[1.1.0]: https://github.com/kcapp/api/compare/v1.0.0...v1.1.0
[1.0.0]: https://github.com/kcapp/api/releases/tag/v1.0.0
