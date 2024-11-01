<div style="display:flex;justify-content:center;align-item:center;">
    <div>
        <img src="./appicon.png" style="height:300px;width:300px">
    </div>
</div>
<div style="display:flex;justify-content:center;">
    <p style="font-size:1.2em;font-weight:bold">Novel Maker - 一款能够制作EPUB的编辑器</p>
</div>
<div style="display:flex;justify-content: center;width:100%;flex-direction:row">
    <a href="./README.md" style="margin-right: 2em">English</a>
    <a href="./README_zh.md" style="margin-right: 2em">简体中文</a>
</div>

## 介绍

这是一款基于Wails框架的富文本编辑器，允许用户导出EPUB文件。

这个项目使用了[wails](https://wails.io/)，更多关于项目配置的详情信息在https://wails.io/docs/reference/project-config。

## 功能
- 现代化界面的编辑器
- 支持Markdown语法
- 可以导出EPUB文件
- 多语言支持

## 未来方向
### v1
- 支持数学公式
- 自定义EPUB样式
- 预览模式

### v2
- EPUB v3支持
- EPUB导入

## 配置需求
### 运行
|环境|版本|
|---|---|
|System|Window 10或更高|
|Webview2|130.x或更高|
|GPU|能打开浏览器就行|
|CPU|能打开浏览器就行|
|Memory|越多越好|

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