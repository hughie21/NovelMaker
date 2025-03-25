export default{
    nav: {
        file: "文件",
        edit: "编辑",
        view: "查看",
        setting: "设置",
        help: "帮助"
    },
    header: {
        file: "文件(F)",
        new: "新建",
        open: "打开",
        import: "导入",
        save: "保存",
        saveAs: "另存为",
        lookUp: "查找",
        replace: "替换",
        exit: "退出",
        edit: "编辑(E)",
        undo: "撤销",
        redo: "恢复",
        copy: "复制",
        cut: "剪切",
        paste: "粘贴",
        selectAll: "全选",
        help: "帮助(H)",
        about: "关于",
    },
    toolBar: {
        file: {
            new: "新增",
            open: "打开",
            import: "导入",
            save: "保存",
            saveAs: "另存为",
            fileInfo: "书籍信息",
            close: "关闭",
            lookUp: "查找",
            toggleCate: "切换目录"
        },
        edit: {
            setting: "设置",
            paste: "粘贴",
            media: "插入图片",
            row: "列",
            column: "行",
            insertTable: "插入表格",
            code: "输入要插入的编程语言",
            link: "超链接",
            abovePlaceHolder: "请输入顶部文字",
        },
        help: {
            help: "帮助",
            about: "关于"
        },
        tooltip: {
            copy: "复制",
            cut: "剪切",
            bold: "粗体",
            italic: "斜体",
            strike: "删除线",
            clean: "清除格式",
            color: "字体颜色",
            background: "背景颜色",
            table: "表格",
            code: "代码块",
            orderlist: "有序列表",
            unorderlist: "无序列表",
            textAlignLeft: "左对齐",
            textAlignCenter: "居中对齐",
            textAlignRight: "右对齐",
            textAlignJustify: "两端对齐",
            aboveText: "顶部文字"
        }
    },
    table: {
        header: "切换表头",
        merge: "合并单元格",
        split: "取消合并",
        insertRow: "插入行",
        insertColumn: "插入列",
        deleteRow: "删除行",
        deleteColumn: "删除列",
        deleteTable: "删除表格",
    },
    tour: {
        next: "下一步",
        previous: "上一步",
        step1: {
            title: "欢迎使用NovelMaker!",
            content: "你可以在此区域编辑你的书籍"
        },
        step2: {
            title: "新建",
            content: "你可以新建一个书籍"
        },
        step3: {
            title: "打开书籍",
            content: "你可以打开一个已有的书籍"
        },
        step4: {
            title: "保存书籍",
            content: "你可以保存你的书籍"
        },
        step5: {
            title: "修改书籍信息",
            content: "你可以在这里修改书籍的信息"
        },
        step6: {
            title: "编辑区域",
            content: "你可以在这里编辑字体的大小、颜色和样式等"
        },
        step7: {
            title: "插入图片",
            content: "你可以在这里插入图片"
        },
        step8: {
            title: "设置",
            content: "你可以在这里设置软件的一些参数"
        },
        step9: {
            title: "最后",
            content: "祝你使用愉快！"
        }
    },
    dialog: {
        confirm: "确认",
        cancel: "取消",
        search: {
            lookUp: "查找",
            replace: "替换",
            defaultResult: "无结果",
            caseSensitive: "区分大小写",
            wholeWord: "全字匹配",
            nextResult: "下一个",
            previousResult: "上一个",
            close: "关闭",
            replaceAll: "替换全部",
            replace: "替换"
        },
        setting: {
            general: "通用",
            title : "设置",
            language: "语言",
            windowSize: "窗口大小",
            resPort: "静态资源服务端口",
            theme: "主题",
            dark: "黑暗主题",
            light: "明亮主题",
            window: "Window 设置",
            maximised: "最大化",
            fullScreen: "全屏",
            normal: "默认",
            windowGPU: "GPU 加速",
            linux: "Linux 设置",
            linuxGPU: "GPU 策略",
            auto: "自动",
            never: "关闭",
            always: "常驻",
            autoSave: "自动保存",
            autoSaveTime: "自动保存间隔",
            epubSaving: "EPUB版面配置",
            textDirec: {
                title: "文本方向",
                ltr: "从左到右",
                rtl: "从右到左",
                auto: "自动"
            },
            layout: {
                title: "内容布局",
                reflow: "流式布局",
                prePaginated: "预分页"
            },
            flow: {
                title: "内容溢出",
                paginated: "动态分页",
                scrolledContinuous: "整体滚动",
                scrolledDoc: "分节滚动",
                auto: "自动"
            },
            spread: {
                title: "双页显示",
                none: "无",
                landscape: "横向",
                portrait: "纵向",
                both: "同时",
                auto: "自动"
            },
            orientation: {
                title: "页面方向",
                auto: "自动",
                landscape: "横向",
                portrait: "纵向"
            },
            proportions: {
                title: "展示比例",
                auto: "自动"
            }
        },
        info: {
            title: "书籍信息",
            name: "书名",
            cover: "封面",
            identifier: "标识码",
            publisher: "出版社",
            author: "作者（多个请用逗号分隔）",
            contributer: "贡献者（多个请用逗号分隔）",
            description: "简介",
            subject: "标签",
            language: "语言",
            requireMessage: "书名不能为空",
            illegalMessage: "输入内容包含非法字符，如 \" 或 \\\\",
        },
        media: {
            title: "媒体管理",
            upload: "上传本地图片",
            link: "网络图片",
            delete: "删除",
            innerTitle: "图片大小",
            insert: "插入",
            cancel: "取消",
            downloadImageTitle: "网络图片",
            downloadImagePromt: "从网络链接处获取图片",
            deleteConfirm: "是否确认删除该项?",
            deleteTitle: "确认删除",
            deleteSuccess: "删除成功",
            deleteFail: "删除失败",
            uploadSuccess: "上传成功",
            uploadFail: "上传失败",
        },
        help: {
            title: "帮助",
            version: "版本",
            date: "日期",
            copyright: "版权",
        },
        link: {
            title: "插入超链接",
            link: "链接",
            text: "请输入网址",
            insert: "插入",
            cancel: "取消",
            error: "输入的内容不符合当前前缀格式，请重新输入!",
            prefix: "链接前缀"
        }
    },
    message: {
        confirm: "确认",
        cancel: "取消",
        error: "错误",
        success: "成功",
        warning: "警告",
        info: "信息",
        notSave: "不保存",
        invalidLink: "链接格式错误",
        imageLoadError: "加载图像失败",
        downloadSuccess: "下载成功",
        downloadError: "下载失败",
        downloadTitle: "下载开始",
        downloadNotice: "在下载的过程中，请不要关闭该软件或者计算机",
        configLoadError: "加载配置文件失败",
        configSaveError: "配置文件保存失败",
        configSaveInfo: "部分配置需要重启后才能生效",
        saveWarning: "您有未保存的编辑，是否保存？",
        saveSuccess: "保存成功",
        saveError: "保存失败",
        openInfo: "请注意，打开非本软件生成EPUB文件可能会丢失部分特殊的格式。",
        openError: "无法打开指定文件",
        importSuccess: "导入成功",
        importError: "导入失败",
        loadingMessage: "解析中...",
        bookinfoSaveError: "请检查输入内容"
    }
}