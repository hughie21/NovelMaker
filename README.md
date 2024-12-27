<p align="center" style="text-align: center">
    <img src="./appicon.png" style="width:40%">
</p>
<p align="center">
    A rich text editor for the EPUB.
    <br/>
    <br/>
    <a href="https://github.com/hughie21/NovelMaker/blob/master/LICENSE">
        <img alt="GitHub" src="https://img.shields.io/github/license/hughie21/NovelMaker"/>
    </a> 
    <a href="https://app.fossa.com/projects/git%2Bgithub.com%2Fhughie21%2FNovelMaker?ref=badge_shield&issueType=license" alt="FOSSA Status">
        <img src="https://app.fossa.com/api/projects/git%2Bgithub.com%2Fhughie21%2FNovelMaker.svg?type=shield&issueType=license"/>
    </a>
    <a href="https://app.fossa.com/projects/git%2Bgithub.com%2Fhughie21%2FNovelMaker?ref=badge_shield&issueType=security" alt="FOSSA Status">
        <img src="https://app.fossa.com/api/projects/git%2Bgithub.com%2Fhughie21%2FNovelMaker.svg?type=shield&issueType=security"/>
    </a>
    <a href="https://github.com/hughie21/NovelMaker/releases">
        <img alt="GitHub release" src="https://img.shields.io/github/release/hughie21/NovelMaker"/>
    </a>
    <a href="https://github.com/hughie21/NovelMaker/issues">
        <img alt="GitHub issues" src="https://img.shields.io/github/issues/hughie21/NovelMaker"/>
    </a>
</p>

<div align="center">

[English](README.md) | [简体中文](README_zh.md)

</div>

## Introduce

This is a rich text editor built with Wails. It can export the content to the epub file.

This project is based on the [wails](https://wails.io/) and more information about the project settings can be found
here: https://wails.io/docs/reference/project-config.

## Feature
- Rich Text Editor with Modern Interface
- Support the markdown grammer
- Can edite epub files
- Multi-Language Support

## Requiremnet
### Running
|Environment|Version|
|---|---|
|OS|Window 10 or 11|
|GPU|Whatever can open the browser|
|CPU|Whatever can run the browser|
|Memory|16G or more|

### Develop
|Environment|Version|
|---|---|
|Nodejs|16.15.1 or upper|
|Wails|2.8.1 or upper|
|Golang|1.20.3 or upper|

## Live Development

To run in the develop environment, you should install all the nodejs requirements as well as the go package.

```bash
# npm
npm install --dependencies

# go
go mod tidy
go mod download
```

After all the requirements have been installed, then run `wails dev`.

## Building

To build a redistributable, production mode package, use `wails build`.

You will find your executable file in the @build/bin.

If you want to build the installation package, use `wails build -nsis`.

For more detail in building see [here](https://wails.io/docs/next/reference/cli).