# Changelog

## [2.6.0] - 2023-06-30
#### Fixed
- Correctly set First 9 Avg. if leg was won in <9 darts
- Return statistics about all `9 Dart Shootout` legs on player statistics

## [2.5.0] - 2023-03-10
#### Feature
- Insert BotConfig on next Leg
- Support for Player Options

#### Changed
- Removed bots and busts from Global statistics

#### Fixed
- Correctly calculate Gotcha statistics
- Ensure x01 statistics only count x01 legs

## [2.4.0] - 2022-10-21
#### Feature
- Change venue when match warmup starts
- New endpoint `/tournament/<id>/matches/result` returning results for all matches in tournament
- Added more general statistics for each tournament to `/tournament/<id>/statistics`
- Added visit statistics for 60+, 100+, 140+, and 180
- Option to mark players as "Placeholder" to be used in tournaments before final players are known
* Probabilities and Odds for Draw of matches based on Binomial Regression
- New `Scam` match type

#### Changed
- Switched from string to time for datetimes to correctly support timezones

#### Fixed
- Correctly rematch for matches with Bots
- Global Fish-n-Chips statistics being off when matches are played without office
- Fixed an issue with three darts average being showed incorrectly for tournament matches before they finish

## [2.3.0] - 2022-03-06
#### Feature
- Support for match presets
- Endpoint for adding tournament groups

#### Changed
- Reverse order of players on rematch

#### Fixed
- Don't allow numbers to be stolen for Tic-Tac-Toe

## [2.2.0] - 2021-12-04
#### Feature
- Smartcard `UID` support for each player
- Support for `BO4-NDS` and `BO2-NDS` mode
- New endpoint for setting leg warmup started
- Endpoint for getting the next tournament match
- Return number of `marks` hit per visit for `Cricket`
- Configured `GitHub Actions`

#### Changed
- Set reverse starting order on Shootout tie breaker legs
- Handle draw for a lot of game types

#### Fixed
- Bug where multiple legs of `Knockout` and `Tic-Tac-Toe` did not work correctly
- Correctly calculate score if `Knockout` is won in 1 visit
- Calculation of `420` scores per visit
- Don't show `9 Dart Shootout` as checkout statistics for Tournament
- Correctly calculate `PPD` for `9 Dart Shootout` where more than 9 darts are thrown
- Misc code smells

## [2.1.0] - 2021-10-17
#### Feature
- Support for new game type `JDC Practice Routine`
- Support for new game type `Knockout`
- Returning `X01 Handicap` statistics`
- Support for modes with a different tie break game type
- Return `active` flag on players

#### Changed
- Set matches as abandoned when legs are cancelled
- Correctly handle draw of 9 Dart Shootout between two players

## [2.0.0] - 2021-09-19
#### Feature
- Support for new game type `Kill Bull`
- Support for new game type `Gotcha`
- New properties for players `board_stream_url` and `board_stream_css`
- Endpoints for loading Elo Changelog for a player
- Start next leg when previous is finished
- Endpoint for getting recent players at a given venue
- Endpoint for getting unfinished matches at a given venue

#### Changed
- Capped `Elo` at lower boundry of `400`
- Added Go module files

#### Fixed
- Misc code fixes

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

[2.6.0]: https://github.com/kcapp/api/compare/v2.5.0...v2.6.0
[2.5.0]: https://github.com/kcapp/api/compare/v2.4.0...v2.5.0
[2.4.0]: https://github.com/kcapp/api/compare/v2.3.0...v2.4.0
[2.3.0]: https://github.com/kcapp/api/compare/v2.2.0...v2.3.0
[2.2.0]: https://github.com/kcapp/api/compare/v2.1.0...v2.2.0
[2.1.0]: https://github.com/kcapp/api/compare/v2.0.0...v2.1.0
[2.0.0]: https://github.com/kcapp/api/compare/v1.2.0...v2.0.0
[1.2.0]: https://github.com/kcapp/api/compare/v1.1.0...v1.2.0
[1.1.0]: https://github.com/kcapp/api/compare/v1.0.0...v1.1.0
[1.0.0]: https://github.com/kcapp/api/releases/tag/v1.0.0
