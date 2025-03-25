export default{
    nav: {
        file: "File",
        edit: "Edit",
        view: "View",
        setting: "Setting",
        help: "Help"
    },
    header: {
        file: "File(F)",
        new: "New",
        open: "Open",
        import: "Import",
        save: "Save",
        saveAs: "Save As",
        lookUp: "Look up",
        replace: "Replace",
        exit: "Exit",
        edit: "Edit(E)",
        undo: "Undo",
        redo: "Redo",
        copy: "Copy",
        cut: "Cut",
        paste: "Paste",
        selectAll: "Select All",
        help: "Help(H)",
        about: "About",
    },
    toolBar: {
        file: {
            new: "New",
            open: "Open",
            import: "Import",
            save: "Save",
            saveAs: "Save As",
            fileInfo: "Book info",
            lookUp: "Look up",
            toggleCate: "Toggle Category"
        },
        edit: {
            setting: "Setting",
            paste: "Paste",
            media: "insert picture",
            row: "Row",
            column: "Column",
            insertTable: "Insert Table",
            code: "Enter the programming language",
            link: "Hyperlink",
            abovePlaceHolder: "Please input the text",
        },
        help: {
            help: "Help",
            about: "About"
        },
        tooltip: {
            copy: "Copy",
            cut: "Cut",
            bold: "Bold",
            italic: "Italic",
            strike: "Strikethrough",
            clean: "Clear Style",
            color: "Font Color",
            background: "Background Color",
            table: "Insert Table",
            code: "Insert Code",
            orderlist: "Order List",
            unorderlist: "Unorder List",
            textAlignLeft: "Text align left",
            textAlignCenter: "Text align center",
            textAlignRight: "Text align right",
            textAlignJustify: "Text align justify",
            aboveText: "Above Text"
        }
    },
    table: {
        header: "Toggle Header",
        merge: "Merge Cell",
        split: "Split Cell",
        insertRow: "Insert Row",
        insertColumn: "Insert Column",
        deleteRow: "Delete Row",
        deleteColumn: "Delete Column",
        deleteTable: "Delete Table",
    },
    tour: {
        next: "Next",
        previous: "Previous",
        step1: {
            title: "Welcome to the NovelMaker!",
            content: "Here you can write your own book!"
        },
        step2: {
            title: "New Books",
            content: "You can create a new book by clicking the new button"
        },
        step3: {
            title: "Open Books",
            content: "You can open a book by clicking the open button"
        },
        step4: {
            title: "Save Books",
            content: "You can save your book by clicking the save button"
        },
        step5: {
            title: "Modify the Book Metadata",
            content: "You can modify the book metadata by clicking the book info button"
        },
        step6: {
            title: "Edit Tab",
            content: "You can change the font size, font style, and font color in this tab"
        },
        step7: {
            title: "Insert Image",
            content: "You can uplaod or insert an image by clicking the image button"
        },
        step8: {
            title: "Setting",
            content: "You can change the setting by clicking the setting button"
        },
        step9: {
            title: "Here we go!",
            content: "Feel free to use the software!"
        }
    },
    dialog: {
        confirm: "Confirm",
        cancel: "Cancel",
        search: {
            lookUp: "Look up",
            replace: "Replace",
            defaultResult: "No result found",
            caseSensitive: "Case sensitive",
            wholeWord: "Whole word",
            nextResult: "Next",
            previousResult: "Previous",
            close: "Close",
            replaceAll: "Replace All",
            replace: "Replace"
        },
        setting: {
            general: "General",
            title : "Setting",
            language: "Language",
            resPort: "Resource Sevice Port",
            windowSize: "Window Size",
            maximised: "Maximised",
            fullScreen: "Full Screen",
            normal: "Default",
            theme: "Theme",
            dark: "Dark",
            light: "Light",
            window: "Window Config",
            windowGPU: "GPU Acceleration",
            linux: "Linux Config",
            linuxGPU: "GPU Policy",
            auto: "Auto",
            never: "Never",
            always: "Always",
            autoSave: "Auto Save",
            autoSaveTime: "Auto Save Time",
            epubSaving: "EPUB Layout Config",
            textDirec: {
                title: "Text Direction",
                ltr: "Left to Right",
                rtl: "Right to Left",
                auto: "Auto"
            },
            layout: {
                title: "Layout",
                reflow: "Reflowable",
                prePaginated: "Pre-paginated"
            },
            flow: {
                title: "Reflow Layout",
                paginated: "Paginated",
                scrolledContinuous: "Scrolled Continuous",
                scrolledDoc: "Scrolled Doc",
                auto: "Auto"
            },
            spread: {
                title: "Spread Layout",
                none: "None",
                landscape: "Landscape",
                portrait: "Portrait",
                both: "Both",
                auto: "Auto"
            },
            orientation: {
                title: "Orientation",
                auto: "Auto",
                landscape: "LandScape",
                portrait: "Portrait"
            },
            proportions: {
                title: "Display Range",
                auto: "Auto"
            }
        },
        info: {
            title: "Book info",
            name: "Book name",
            cover: "Cover",
            identifier: "Book identifier",
            publisher: "Publisher",
            author: "Authors (separate more than one with a comma)",
            contributer: "Contributors (separate more than one with a comma)",
            description: "Description",
            subject: "Subject",
            language: "Language",
            requireMessage: "The title of the book is required",
            illegalMessage: "The input contains illegal characters, like \" or \\\\",
        },
        media: {
            title: "Media manager",
            upload: "Upload local image",
            delete: "Delete",
            link: "Image link",
            innerTitle: "Size of the Picture",
            insert: "Insert",
            cancel: "Cancel",
            downloadImageTitle: "Web Image",
            downloadImagePromt: "Download Image from link",
            deleteConfirm: "Are you sure to delete this media?",
            deleteTitle: "Delete Confirmation",
            deleteSuccess: "Delete success",
            deleteFail: "Delete failed",
            uploadSuccess: "Upload success",
            uploadFail: "Upload failed",
        },
        help: {
            title: "Help",
            version: "Version",
            date: "Date",
            copyright: "Copyright",
        },
        link: {
            title: "Insert the hyperlink",
            link: "Link",
            text: "Pls input the URL",
            insert: "Insert",
            cancel: "Cancel",
            error: "The input does not match the current prefix format, please re-enter!",
            prefix: "Link Prefix"
        }
    },
    message: {
        confirm: "Confirm",
        cancel: "Cancel",
        error: "Error",
        success: "Success",
        warning: "Warning",
        notSave: "Do not save",
        info: "Info",
        invalidLink: "Invalid link format",
        imageLoadError: "Fail to load the image",
        downloadSuccess: "Download success",
        downloadError: "Download error",
        downloadTitle: "Download Start",
        downloadNotice: "Downloading the file, please do not close the software or turn off the computer.",
        configLoadError: "Fail to load config file",
        configSaveError: "Fail to save config",
        configSaveInfo: "Part of the configuration need to restart the software to take effect",
        saveWarning: "The edited file have not been saved, do you want to save it?",
        openInfo: "Some special formatting may be lost when opening EPUB files not generated by this software.",
        openError: "Failed to open the file",
        saveSuccess: "Save success",
        saveError: "Save error",
        importSuccess: "Import success",
        importError: "Import error",
        loadingMessage: "Loading...",
        bookinfoSaveError: "Please check the input"
    }
}