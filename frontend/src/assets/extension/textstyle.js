/*
@Author: Hughie
@CreateTime: 2024-9-21
@LastEditors: Hughie
@LastEditTime: 2025-3-28
@Description: the text style extension of tiptap
*/

import TextStyle from '@tiptap/extension-text-style'

const TextStyleExtends = TextStyle.extend({
    priority: 1000,
    addAttributes() {
        return {
            ...this.parent?.(),
            fontSize: {
                default: "1rem", // default value
            },
            fontColor: {},
            backgroundColor: {
                default: "transparent",
            },
            fontFamily: {
                default: "Arial",
            }
        }
    },
    addOptions() {
        return {
            ...this.parent?.()
        }
    },
    renderHTML({ HTMLAttributes }) {
        return ['span', { 
            style: `font-size: ${HTMLAttributes.fontSize};color: ${HTMLAttributes.fontColor};background: ${HTMLAttributes.backgroundColor};font-family: ${HTMLAttributes.fontFamily};`
        }, 0];
    },
    parseHTML() {
        return [{
            tag: 'span',
            getAttrs: node => {
                const fontSize = node.style.fontSize;
                const fontColor = node.style.color;
                const backgroundColor = node.style.background;
                const fontFamily = node.style.fontFamily;
                return { fontSize, fontColor, backgroundColor, fontFamily };
            },
        }]
    },
    addCommands() {
        return {
            setFontSize: size => ({ commands }) => {
                return commands.setMark(this.name, { fontSize: size });
            },
            setColor: color => ({ commands }) => {
                return commands.setMark(this.name, { fontColor: color });
            },
            setTextBackColor: color => ({ commands }) => {
                return commands.setMark(this.name, { backgroundColor: color });
            },
            setFontFamily: family => ({ commands }) => {
                return commands.setMark(this.name, { fontFamily: family });
            }
        }
    }
 })

// this is a extension for the search selection background
const SearchSelBackground = TextStyle.extend({
    priority: 1000,
    name: 'seach-select-background',
    addAttributes() {
        return {
            ...this.parent?.(),
            backgroundColor: {
                default: "rgba(164, 61, 61, 0.6)",
            },
            color: {
                default: "#ffffff"
            }
        }
    },
    addOptions() {
        return {
            ...this.parent?.()
        }
    },
    renderHTML({ HTMLAttributes }) {
        return ['span', { style: `background-color: ${HTMLAttributes.backgroundColor};color: ${HTMLAttributes.color}` }, 0];
    },
    addCommands() {
        return {
            setSearchSelection: (range) => ({ chain }) => {
                chain().setTextSelection(range).setMark(this.name).run();
                return;
            },
            unsetSearchSelection: () => ({ commands }) => {
                return commands.unsetMark(this.name)
            }
        }
    }
})

export {
    SearchSelBackground,
    TextStyleExtends,
}

export default {
    
}