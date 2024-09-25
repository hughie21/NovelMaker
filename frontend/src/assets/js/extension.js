/*
@Author: Hughie
@CreateTime: 2024-9-21
@LastEditors: Hughie
@LastEditTime: 2024-09-25
@Description: This is the extension for the tiptap.
*/

import { nodeInputRule, textblockTypeInputRule } from "@tiptap/vue-3"
import Image from "@tiptap/extension-image";
import Heading from "@tiptap/extension-heading"
import { imageInfo } from "./utils";
import TextStyle from '@tiptap/extension-text-style'

const TextFontSize = TextStyle.extend({
    priority: 1000,
    name: 'font-size',
    addAttributes() {
        return {
            ...this.parent?.(),
            fontSize: {
                default: "14px", // default value
            },
        }
    },
    addOptions() {
        return {
            ...this.parent?.()
        }
    },
    renderHTML({ HTMLAttributes }) {
        return ['span', { style: `font-size: ${HTMLAttributes.fontSize};` }, 0];
    },
    parseHTML() {
        return [{
            tag: 'span',
            getAttrs: node => {
                const fontSize = node.style.fontSize;
                return { fontSize };
            },
        }]
    },
    addCommands() {
        return {
            setFontSize: size => ({ commands }) => {
                return commands.setMark(this.name, { fontSize: size });
            },
            unsetTextBackColor: () => ({ commands }) => {
                return commands.unsetMark(this.name)
            }
        }
    }
})

const TextBackground = TextStyle.extend({
    priority: 1000,
    name: 'background-color',
    addAttributes() {
        return {
            ...this.parent?.(),
            backgroundColor: {
                default: "transparent", // default value
            },
        }
    },
    addOptions() {
        return {
            ...this.parent?.()
        }
    },
    renderHTML({ HTMLAttributes }) {
        return ['span', { style: `background-color: ${HTMLAttributes.backgroundColor};` }, 0];
    },
    parseHTML() {
        return [{
            tag: 'span',
            getAttrs: node => {
                const backgroundColor = node.style.backgroundColor;
                return { backgroundColor };
            },
        }]
    },
    addCommands() {
        return {
            setTextBackColor: color => ({ commands }) => {
                return commands.setMark(this.name, { backgroundColor: color });
            },
            unsetTextBackColor: () => ({ commands }) => {
                return commands.unsetMark(this.name)
            }
        }
    }
})

const CustomHeading = Heading.extend({
    priority: 1000,
    name: 'custom-heading',
    addOptions(){
        return {
            ...this.parent?.()
        };
    },
    content: 'inline*',
    group: 'block',
    defining: true,
    addAttributes() {
        return {
            ...this.parent?.(),
        }
    },
    parseHTML(){
        return this.options.levels
            .map(function (level) { return ({
            tag: "h".concat(level),
            attrs: { level: level },
        }); });
    },
    renderHTML({ HTMLAttributes, node }) {
        const hasLevel = this.options.levels.includes(node.attrs.level);
        var level = hasLevel
            ? node.attrs.level
            : this.options.levels[0];
        let uuid = new Date().getTime().toString(36) + Math.random().toString(36).slice(2,9);
        HTMLAttributes.id = `tap-heading-${uuid}`;
        return [`h${level}`, HTMLAttributes, 0 ];
    },
    addCommands() {
        return {
            setHeader: level => ({ commands }) => {
                if(level == 0) {
                    return commands.setNode('paragraph', {})
                }
                if (!this.options.levels.includes(level)){
                    return false;
                }
                let uuid = new Date().getTime().toString(36) + Math.random().toString(36).slice(2,9);
                return commands.setNode(this.name, {level: level, id: `tap-heading-${uuid}`});
            },
            toggleHeader: level => ({ commands }) => {
                if(level == 0) {
                    return commands.toggleNode('heading' ,'paragraph', {})
                }
                if (!this.options.levels.includes(level)) {
                    return false
                }
                let uuid = new Date().getTime().toString(36) + Math.random().toString(36).slice(2,9);
                return commands.toggleNode(this.name, 'paragraph', {level: level, id: `tap-heading-${uuid}`})
            }
        }
    },
    addKeyboardShortcuts() {
        return {
            ...this.parent?.()
        }
    },
    addInputRules() {
        return this.options.levels.map(level => {
            return textblockTypeInputRule({
              find: new RegExp(`^(#{1,${level}})\\s$`),
              type: this.type,
              getAttributes: {
                level,
              },
            })
        })
    }
})

const CustomImage = Image.extend({
    draggable:false,
    addAttributes() {
        return {
            ...this.parent?.(),
            src: {
                default: null,
            },
            alt: {
                default: null,
            },
            title: {
                default: null,
            },
            zoom: {
                default: 100,
            },
            pos: {
                default: 'left',
            }
        }
    },
    addOptions() {
        return {
            ...this.parent?.()
        }
    },
    renderHTML({ HTMLAttributes }) {
        return ["span", {
            style: `display:flex;justify-content:${HTMLAttributes.pos};`,
        }  ,['img', {
            src: HTMLAttributes.src,
            alt: HTMLAttributes.alt,
            title: HTMLAttributes.title,
            style: `width:${HTMLAttributes.zoom}%;`,
        }]]
    },
    parseHTML() {
        return [{
            tag: 'span',
            getAttrs: node => {
                const imgNode = node.querySelector('img');
                if (!imgNode) return false;
                const src = imgNode.getAttribute('src');
                const alt = imgNode.getAttribute('alt');
                const title = imgNode.getAttribute('title') || '';
                const widthStyle = imgNode.style.width;
                const zoom = widthStyle ? parseInt(widthStyle.replace('%', ''), 10) : 100;
                const pos = node.style.justifyContent || 'left';
                return { src, alt, title, zoom, pos };
            },
        }]
    },
    addNodeView() {
        return node => {
            const imgNode = document.createElement('img');
            let { src, alt, title, zoom, pos } = node.HTMLAttributes;
            imgNode.src = src;
            imgNode.alt = alt;
            imgNode.title = title;
            imgNode.style.width = `${zoom}%`;
            imgNode.addEventListener("click", (e)=>{
                imageInfo.elem = e.target;
                imageInfo.zoom = parseInt(imgNode.style.height);
                imageInfo.postition = imgNode.parentNode.style.justifyContent;
            })

            let container = document.createElement('span');
            container.className = 'image-container';
            container.style.justifyContent = pos;
            container.draggable = false;
            container.appendChild(imgNode);

            this.contentDOM = imgNode;
            this.dom = container;

            let destroy = () => {
                container.parentNode.removeChild(container);
            };
            let update = (updateNode) => {
                if (updateNode.type.name !== this.name) {
                    return false;
                }
                this.contentDOM.style.width = `${updateNode.attrs.zoom}%`;
                this.dom.style.justifyContent = updateNode.attrs.pos;
                return true;
            }

            return {
                dom: container,
                contentDOM: imgNode,
                update: update,
                destroy: destroy
            };
        }
    },
    addCommands(){
        return {
            InsertImage: (options) => ({ commands }) => {
                const { src, alt, title, zoom, pos } = options
                return commands
                .insertContent({
                    type: 'image',
                    attrs: {
                        src,
                        alt,
                        title,
                        zoom,
                        pos
                    }
                }
                )
            }
        }
    },
    addInputRules(){
        const inputReg = /(?:^|\s)(!\[(.+|:?)]\((\S+)(?:(?:\s+)["'](\S+)["'])?\))$/
        return [
            nodeInputRule({
                find: inputReg,
                type: this.type,
                getAttributes: match => {
                    const [,, alt, src, title] = match
                    return {
                        src,
                        alt,
                        title
                    }
                }
            })
        ]
    }
})

export {
    CustomImage,
    CustomHeading,
    TextFontSize,
    TextBackground,
}

export default {
    
}