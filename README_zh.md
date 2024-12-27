<p align="center" style="text-align: center">
    <img src="./appicon.png" style="width:40%">
</p>
<p align="center">
    一款EPUB富文本编辑器
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

## 介绍
这是一款基于Wails框架的富文本编辑器，允许用户导出EPUB文件。

这个项目使用了[wails](https://wails.io/)，更多关于项目配置的详情信息在https://wails.io/docs/reference/project-config。

## 功能
- 现代化界面的编辑器
- 支持Markdown语法
- 可以编辑EPUB文件
- 多语言支持

## 配置需求
### 运行
|环境|版本|
|---|---|
|系统|Window 10或更高|
|GPU|能打开浏览器就行|
|CPU|能打开浏览器就行|
|内存|16G或以上|

### 开发
|环境|版本|
|---|---|
|Nodejs|16.15.1或更高|
|Wails|2.8.1或更高|
|Golang|1.20.3或更高|

## 本地开发

要在开发环境中运行，您首先应该安装所有nodejs拓展以及go模块。

```bash
# npm
npm install --dependencies

# go
go mod tidy
go mod download
```

等待所有依赖安装完成后，运行`wails dev`。

## 打包

要构建分发的可执行文件，运行`wails build`。

您将会在@build/bin中找到该文件。

如果您想要构建安装文件，运行`wails build -nsis`.

更多构建命令请查看[这里](https://wails.io/docs/next/reference/cli).
