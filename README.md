<div style="display:flex;justify-content:center;align-item:center;">
    <div>
        <img src="./appicon.png" style="height:300px;width:300px">
    </div>
</div>
<div style="display:flex;justify-content:center;">
    <p style="font-size:1.2em;font-weight:bold">Novel Maker - A text editor to make the EPUB</p>
</div>
<div style="display:flex;justify-content: center;width:100%;flex-direction:row">
    <a href="./README.md" style="margin-right: 2em">English</a>
    <a href="./README_zh.md" style="margin-right: 2em">简体中文</a>
</div>

## Introduce

This is a rich text editor built with Wails. It can export the content to the epub file.

This project is based on the [wails](https://wails.io/) and more information about the project settings can be found
here: https://wails.io/docs/reference/project-config.

## Feature
- Rich Text Editor with Modern Interface
- Support the markdown
- Can export to the EPUB file
- Multi-Language Support

## Future Roadmap
### v1
- Mathematical Formulas
- Custom EPUB Style
- Preview Mode

### v2
- EPUB v3 Support
- EPUB Import 

## Requiremnet
### Running
|Environment|Version|
|---|---|
|System|Window 10 or 11|
|Webview2|130.x or upper|
|GPU|Whatever can open the browser|
|CPU|Whatever can run the browser|
|Memory|The more the best|

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