# ChangeLog
All notable changes to this project will be documented in this file.
Explanation of version number format: major.minor.patch.date_suffix
- Major version update when the program structure is changed
- Minor version update when the program function is changed
- Patch version update when the program bug is fixed
- Date: YYYYMMDD, the date of the version update

- Suffix:
    - Beta version is a stable test version
    - Alpha version is an early development version

<hr/>

## [v1.1.0.20250407] - 2025-04-07
### Added
- Add text alignment and allow math formulas to be entered and displayed within the editor [#a30ea21](https://github.com/hughie21/NovelMaker/commit/a30ea21527dd2ef6164fce067fc5ac4541527a77)
- Add support for note formatting (a small line of text at the top of the text) [#56c3983](https://github.com/hughie21/NovelMaker/commit/56c39835efc5ae69a64ae888d26b2f9ef39b0771)
- Add support for parsing mathematical formulas on the backend, add katex style and JavaScript files to epub files [#04c69db](https://github.com/hughie21/NovelMaker/commit/04c69dbec2d3f0ee77692e5760a0e8dc26c29382)

## [v1.0.0.20250314] - 2025-03-14
### Added
- Added catalogue node highlighting with scrolling [#129c135](https://github.com/hughie21/NovelMaker/commit/129c13589606a891771e05b3af3c02160565799b)
### Fix
- Fix the bug that some EPUB files cannot be parsed [#1920609](https://github.com/hughie21/NovelMaker/commit/19206092ef7bb17a2eb8b0e6930a1abc96e840d4)
- Update the vue-i18n version to 10.0.6 to aviod the potiential security issue [CVE-2025-27597](https://github.com/advisories/GHSA-p2ph-7g93-hw3m).

## [v1.0.0.20250104_beta] - 2025-01-04
### Fix [#07a7e7e](https://github.com/hughie21/NovelMaker/commit/07a7e7ee140e5df8f4032949e111a001019f8ad4)
- When a network error occurs while downloading an image, the error is not truncated causing the program to crash.
- When there are spaces in the epub file name, the file name is incomplete when openning the file directly.
- Directly open the file and the directory is not refreshed.
- Save the wrong image path when saving the epub.
- Update `vite` to 5.0.11 to fix the [CVE-2025-24010](https://github.com/advisories/GHSA-vg6x-rcgg-rjx6)

## [v1.0.0.20241227_beta] - 2024-12-27
### Added
- Add configuration items related to EPUB file layout control to the config file and corresponding options to the frontend settings. [#53f6c37](https://github.com/hughie21/NovelMaker/commit/53f6c37fedf880bde8b0f544594c44735dc861e1)
- Add several text format parsers to support more text formats like color and background color. [$9779319](https://github.com/hughie21/NovelMaker/commit/9779319ef27c42d4771ad55bde5630a71a8de294)
- Add a function that can save the file automatically at a set interval. [#31d3ad7](https://github.com/hughie21/NovelMaker/commit/31d3ad7d07a32c1d3f6e3ca7b44b0fed8e9e8d49)
- Add a new feature to open the file from the command line. [#189ca19](https://github.com/hughie21/NovelMaker/commit/189ca190be1431dae4fd7a46b3b43224be508445)
- Add a new feature for front-end error logging. [#cab5cfa](https://github.com/hughie21/NovelMaker/commit/cab5cfa45d9f791ac6e37d3583d53a85dbf45f98)
### Change
- Modify the logic of lookup.js to avoid time-consuming operations that could lead to program crashes. [#97bd667](https://github.com/hughie21/NovelMaker/commit/97bd667831a86f20e13ad2b2476965d66ac500a7)
- Modify the reader parsing logic to add parsing of EPUB package file profiles and tags. [#78b9a43](https://github.com/hughie21/NovelMaker/commit/78b9a4363d1d9125439df19f9fed6f0de5081201)
- Refactor the backend to make it more modular and easier to maintain, integrating backend functionality into a unified lib folder. [#71b02d1](https://github.com/hughie21/NovelMaker/commit/71b02d165f5fb3f7bec43d7ad8de8c4c1422e5f5)
- Modify the way ImageParser parses the attributes of the img tag for correct positioning. [#78b9a43](https://github.com/hughie21/NovelMaker/commit/78b9a4363d1d9125439df19f9fed6f0de5081201)
- Refactor the message structure in the logging module to a message interface. [#cab5cfa](https://github.com/hughie21/NovelMaker/commit/cab5cfa45d9f791ac6e37d3583d53a85dbf45f98)
### Fix
- Fix bugs in production environment that when opening epub directly the startup path of the programme is wrong. [#e8fd254](https://github.com/hughie21/NovelMaker/commit/e8fd254afb92b03991bc8d3504f0e6dc0b79d8a9)
- Fix potential bugs causing the program to get stuck when the number of texts to be replaced exceeds 100. [#97bd667](https://github.com/hughie21/NovelMaker/commit/97bd667831a86f20e13ad2b2476965d66ac500a7)
- Fix the issue of not being able to parse overlay styles correctly. [#12efa4c](https://github.com/hughie21/NovelMaker/commit/12efa4c8b95ef4b1926d60a4489b381e3a94ecbc)
- Fix extra spaces and line breaks in imported EPUB text. [#85e8397](https://github.com/hughie21/NovelMaker/commit/85e83972a6a8b666125038d47eefdcaa4ed19cc5)
- Fix the bug that the exported EPUB file's text.xhtml contains incorrect image paths. [#c847133](https://github.com/hughie21/NovelMaker/commit/c847133e236bec3e5951218dde8d26a0e2e28a05)
- Fix the bug where the header element's level cannot be parsed correctly. [#a97f175](https://github.com/hughie21/NovelMaker/commit/a97f175587b5f927d29c5587c2c7e22a3c34a6a4)

## [v1.0.0.20241101_alpha] - 2024-11-01
### Added
- Custom App Header: Implemented a custom header supporting multiple languages, custom display file names, and multi-theme colors.  [#d75989c](https://github.com/hughie21/NovelMaker/commit/d75989c5bc268620d3981f7477ce6d57b767d06d)
- Front-end Configuration Setting Improvement: Added an interface for modifying app initialization configurations to the front-end settings options. Also refactored the tab component and added a new tab.css file for styling. [#460cbb0](https://github.com/hughie21/NovelMaker/commit/460cbb00b71a7a6147ddb489b57f3bfa0b0346b7)
- Lookup and Replace Feature: Added functionality allowing users to search and replace terms within the editor, handling both case-sensitive and insensitive searches. [#95bb3cb](https://github.com/hughie21/NovelMaker/commit/95bb3cb448e7410f79ce8298671a9d2eb1668ce8)
### Change
- More Reasonable Logging: Refactored the logging package for better usability and introduced standardized log modules. Improved project.nsi file. [#223b523](https://github.com/hughie21/NovelMaker/commit/223b523f99c81af8d0bb903f957f625c045a5b01)
- Editor Replacement: Replaced the current editor with Tiptap editor to address performance issues such as lagging while typing. [#e69dcce](https://github.com/hughie21/NovelMaker/commit/e69dcce8763e2db09c046b49aed7fa27c6761221)
- Icon Replacement: Updated the software’s icon. [#ef60f3e](https://github.com/hughie21/NovelMaker/commit/ef60f3e0de9138f32c331cd1fbfe4f1dcb4e3dc6)
### Fix
- Image Display in Exported EPUB File: Resolved the issue where images weren’t displayed correctly in exported EPUB files; fixed the source attribute problem. [#d75989c](https://github.com/hughie21/NovelMaker/commit/d75989c5bc268620d3981f7477ce6d57b767d06d)
- Media Dialog Upload and Delete Issues: Fixed bugs preventing upload and deletion functionalities in the media dialog under production environments. [#223b523](https://github.com/hughie21/NovelMaker/commit/223b523f99c81af8d0bb903f957f625c045a5b01)
- Temp File Removal: Addressed the issue of temporary files remaining post-export by ensuring OS readers are closed properly during file export. [#0a60574](https://github.com/hughie21/NovelMaker/commit/0a60574a02b3168093927ed713e52e848e141624)
- Font Size Issue: Corrected the inability to change the font size of the editor's content. [#37a7847](https://github.com/hughie21/NovelMaker/commit/37a7847d4ab4942d36d1ef3403f72d90882761ae)
- File Reading Issue: Solved problems related to file reading in the file information dialog and special character handling in content when exporting to EPUB. [#e69dcce](https://github.com/hughie21/NovelMaker/commit/e69dcce8763e2db09c046b49aed7fa27c6761221)