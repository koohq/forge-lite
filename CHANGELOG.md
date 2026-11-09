# Changelog

## [0.2.0](https://github.com/koohq/forge-lite/compare/v0.1.1...v0.2.0) (2026-11-09)


### Features

* add Auto Explore to keep descending until a stop point ([cb8fdf9](https://github.com/koohq/forge-lite/commit/cb8fdf9fb97305fe8a3e205366d067ad06f394e3))
* mark B100F of the endless dungeon as the main-game clear ([d997458](https://github.com/koohq/forge-lite/commit/d9974585c71a0df70510ce7ee937f4ddbdbd240b))


### Bug Fixes

* join English enemy names without doubled spaces ([6723470](https://github.com/koohq/forge-lite/commit/672347074118997fa03e9cde712b68d8479bf347))
* stop showing the defeat message when fleeing mid-battle ([c18389f](https://github.com/koohq/forge-lite/commit/c18389f3ccf51f813b043708f6637ebfc3c0ba82))
* treat old saves with a record past B100F as cleared ([8b2843a](https://github.com/koohq/forge-lite/commit/8b2843a6fb0c7adbdfd61573a58366cac24116ba))
* write Classic Fantasy weapon rank names in kanji ([0fd7507](https://github.com/koohq/forge-lite/commit/0fd75072ceececb94c567b25afda08b391ada944))

## [0.1.1](https://github.com/koohq/forge-lite/compare/v0.1.0...v0.1.1) (2026-09-15)


### Bug Fixes

* normalize version output by stripping the leading v ([#3](https://github.com/koohq/forge-lite/issues/3)) ([b9f23b6](https://github.com/koohq/forge-lite/commit/b9f23b6afb5a8a7ae5ae9b0c0c7707e5ecc402bf))

## 0.1.0 (2026-09-15)


### Features

* add --version flag reporting release or VCS build info ([9d8c439](https://github.com/koohq/forge-lite/commit/9d8c4392ab0bbf78ebc53013527d791325e52777))
* add dynamic theme vocabulary for hub display and menu ([80cec70](https://github.com/koohq/forge-lite/commit/80cec703c9347b03cb668bafe8e3c0f65e206738))
* add endless dungeon 'Infinite Abyss' and max HP training feature ([661ff4e](https://github.com/koohq/forge-lite/commit/661ff4ec057fcfc48452ef331264b1f683bc1fcc))
* add explanatory subtitles to hub menu items across all themes ([ca1b0b4](https://github.com/koohq/forge-lite/commit/ca1b0b4e660c91d8f1ebbe8fee4685551efd0b24))
* add full bilingual support for theme presets and refactor theme schema ([97b7910](https://github.com/koohq/forge-lite/commit/97b79100b04ffcccef8465698e0a0ac9a947a80b))
* add orb disassembly and synthesis features ([b3ff634](https://github.com/koohq/forge-lite/commit/b3ff63415d9d2c564457dff7d7029ef3f7097bbb))
* add QoL features and external theme configuration ([e068b83](https://github.com/koohq/forge-lite/commit/e068b836731500949d56e8f600769714aed2a82f))
* add save/load system, upper dungeon, and bulk upgrade ([a01525e](https://github.com/koohq/forge-lite/commit/a01525eb931c47579c950ea742373a95d85203a8))
* add target ranks progression, dynamic hub UI, and partner sync rework ([0f9a81e](https://github.com/koohq/forge-lite/commit/0f9a81ed41d3e5e9800a385fcc710286e656488f))
* implement system i18n and customizable theme presets ([b5e83ab](https://github.com/koohq/forge-lite/commit/b5e83ab0b5960893f5798acbacebca2d40c70055))
* port CLI roguelite game to Go 1.27.1 ([7e98fcc](https://github.com/koohq/forge-lite/commit/7e98fccde853c1dc331dcb2e31eb7a812e0bf27d))


### Bug Fixes

* adjust Partner Sync enemy prefixes with world-building cues ([0dceb0d](https://github.com/koohq/forge-lite/commit/0dceb0d1c97ce4c21c4c7c51d83eaedb5d55c98b))
* correct module path to canonical GitHub import path ([26b3415](https://github.com/koohq/forge-lite/commit/26b34154fa5d3a4e9d4f6a0c70b75c3acabb2538))
* polish custom theme example orbs, hubTitle, and enemy prefixes ([56e6b55](https://github.com/koohq/forge-lite/commit/56e6b5530fc38370429281ffcbfbc46570113ac6))
* polish theme labels, messages, and enemy name generation ([e6381a2](https://github.com/koohq/forge-lite/commit/e6381a2942689900b4875fc97e209e60716e02c6))
* refine theme wording and fix label duplication in partner sync ([ef393e0](https://github.com/koohq/forge-lite/commit/ef393e088458b7d14c950a2022b0fa031eee5a14))
* update Partner Sync menu 3 label to 連携スキルを設定 ([265887c](https://github.com/koohq/forge-lite/commit/265887c91682a52bc00cc3f2aa9c7f01a6217594))
