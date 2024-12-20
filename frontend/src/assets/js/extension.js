/*
@Author: Hughie
@CreateTime: 2024-9-21
@LastEditors: Hughie
@LastEditTime: 2024-12-20
@Description: This is the extension for the tiptap.
*/

import { nodeInputRule, textblockTypeInputRule } from "@tiptap/vue-3"
import Image from "@tiptap/extension-image";
import Heading from "@tiptap/extension-heading"
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
        // this is a random id for the heading
        // but this id will not be used in the exported epub file
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


// @reference: https://github.com/bae-sh/tiptap-extension-resize-image/blob/main/lib/imageResize.ts:83-311
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
            style: {
                default: 'width: 300px; height: auto; cursor: pointer;',
                parseHTML: (element) => {
                    const width = element.getAttribute('width');
                    return width
                        ? `width: ${width}px; height: auto; cursor: pointer;`
                        : `${element.style.cssText}`;
                },
            }
        }
    },
    addOptions() {
        return {
            ...this.parent?.()
        }
    },
    addNodeView() {
        return ({ node, editor, getPos }) => {
            const { view, options: { editable }, } = editor;
            const { style } = node.attrs;
            const $wrapper = document.createElement('div');
            const $container = document.createElement('div');
            const $img = document.createElement('img');
            const iconStyle = 'width: 18px; height: 18px; cursor: pointer;';
            const dispatchNodeView = () => {
                if (typeof getPos === 'function') {
                    const newAttrs = Object.assign(Object.assign({}, node.attrs), { style: `${$img.style.cssText}` });
                    view.dispatch(view.state.tr.setNodeMarkup(getPos(), null, newAttrs));
                }
            }

            const createIcon = (svg) => {
                const $icon = document.createElement('i');
                $icon.innerHTML = svg
                return $icon;
            }
            const paintPositionContoller = () => {
                const $postionController = document.createElement('div');
                const $leftController = createIcon(`<svg t="1733712915480" class="icon" viewBox="0 0 1024 1024" version="1.1" xmlns="http://www.w3.org/2000/svg" p-id="8173" width="200" height="200" style="fill: #000 !important"><path d="M128 213.333333l768 0q17.664 0 30.165333 12.501333t12.501333 30.165333-12.501333 30.165333-30.165333 12.501333l-768 0q-17.664 0-30.165333-12.501333t-12.501333-30.165333 12.501333-30.165333 30.165333-12.501333zM128 725.333333l597.333333 0q17.664 0 30.165333 12.501333t12.501333 30.165333-12.501333 30.165333-30.165333 12.501333l-597.333333 0q-17.664 0-30.165333-12.501333t-12.501333-30.165333 12.501333-30.165333 30.165333-12.501333zM128 554.666667l768 0q17.664 0 30.165333 12.501333t12.501333 30.165333-12.501333 30.165333-30.165333 12.501333l-768 0q-17.664 0-30.165333-12.501333t-12.501333-30.165333 12.501333-30.165333 30.165333-12.501333zM128 384l597.333333 0q17.664 0 30.165333 12.501333t12.501333 30.165333-12.501333 30.165333-30.165333 12.501333l-597.333333 0q-17.664 0-30.165333-12.501333t-12.501333-30.165333 12.501333-30.165333 30.165333-12.501333z" fill="#444444" p-id="8174"></path></svg>`)
                const $centerController = createIcon(`<svg t="1733713030570" class="icon" viewBox="0 0 1024 1024" version="1.1" xmlns="http://www.w3.org/2000/svg" p-id="8489" width="200" height="200" style="fill: #000 !important"><path d="M128 213.333333l768 0q17.664 0 30.165333 12.501333t12.501333 30.165333-12.501333 30.165333-30.165333 12.501333l-768 0q-17.664 0-30.165333-12.501333t-12.501333-30.165333 12.501333-30.165333 30.165333-12.501333zM213.333333 725.333333l597.333333 0q17.664 0 30.165333 12.501333t12.501333 30.165333-12.501333 30.165333-30.165333 12.501333l-597.333333 0q-17.664 0-30.165333-12.501333t-12.501333-30.165333 12.501333-30.165333 30.165333-12.501333zM128 554.666667l768 0q17.664 0 30.165333 12.501333t12.501333 30.165333-12.501333 30.165333-30.165333 12.501333l-768 0q-17.664 0-30.165333-12.501333t-12.501333-30.165333 12.501333-30.165333 30.165333-12.501333zM213.333333 384l597.333333 0q17.664 0 30.165333 12.501333t12.501333 30.165333-12.501333 30.165333-30.165333 12.501333l-597.333333 0q-17.664 0-30.165333-12.501333t-12.501333-30.165333 12.501333-30.165333 30.165333-12.501333z" fill="#444444" p-id="8490"></path></svg>`)
                const $rightController = createIcon(`<svg t="1733712984344" class="icon" viewBox="0 0 1024 1024" version="1.1" xmlns="http://www.w3.org/2000/svg" p-id="8331" width="200" height="200" style="fill: #000 !important"><path d="M128 213.333333l768 0q17.664 0 30.165333 12.501333t12.501333 30.165333-12.501333 30.165333-30.165333 12.501333l-768 0q-17.664 0-30.165333-12.501333t-12.501333-30.165333 12.501333-30.165333 30.165333-12.501333zM298.666667 725.333333l597.333333 0q17.664 0 30.165333 12.501333t12.501333 30.165333-12.501333 30.165333-30.165333 12.501333l-597.333333 0q-17.664 0-30.165333-12.501333t-12.501333-30.165333 12.501333-30.165333 30.165333-12.501333zM128 554.666667l768 0q17.664 0 30.165333 12.501333t12.501333 30.165333-12.501333 30.165333-30.165333 12.501333l-768 0q-17.664 0-30.165333-12.501333t-12.501333-30.165333 12.501333-30.165333 30.165333-12.501333zM298.666667 384l597.333333 0q17.664 0 30.165333 12.501333t12.501333 30.165333-12.501333 30.165333-30.165333 12.501333l-597.333333 0q-17.664 0-30.165333-12.501333t-12.501333-30.165333 12.501333-30.165333 30.165333-12.501333z" fill="#444444" p-id="8332"></path></svg>`)

                const controllerMouseOver = (e) => {
                    e.target.style.opacity = 0.3;
                };
                const controllerMouseOut = (e) => {
                    e.target.style.opacity = 1;
                };

                $postionController.setAttribute('style', 'position: absolute; top: 5%; left: 50%; width: 100px; height: 25px; z-index: 999; background-color: rgba(255, 255, 255, 0.7); border-radius: 4px; border: 2px solid #6C6C6C; cursor: pointer; transform: translate(-50%, -50%); display: flex; justify-content: space-between; align-items: center; padding: 0 10px;');

                $leftController.setAttribute('style', iconStyle);
                $leftController.addEventListener('mouseover', controllerMouseOver);
                $leftController.addEventListener('mouseout', controllerMouseOut);

                $centerController.setAttribute('style', iconStyle);
                $centerController.addEventListener('mouseover', controllerMouseOver);
                $centerController.addEventListener('mouseout', controllerMouseOut);

                $rightController.setAttribute('style', iconStyle);
                $rightController.addEventListener('mouseover', controllerMouseOver);    
                $rightController.addEventListener('mouseout', controllerMouseOut);

                $leftController.addEventListener('click', () => {
                    $img.setAttribute('style', `${$img.style.cssText} margin: 0 auto 0 0;`);
                    dispatchNodeView();
                });

                $centerController.addEventListener('click', () => {
                    $img.setAttribute('style', `${$img.style.cssText} margin: 0 auto;`);
                    dispatchNodeView();
                });

                $rightController.addEventListener('click', () => {
                    $img.setAttribute('style', `${$img.style.cssText} margin: 0 0 0 auto;`);
                    dispatchNodeView();
                })

                $postionController.appendChild($leftController);
                $postionController.appendChild($centerController);
                $postionController.appendChild($rightController);
                $container.appendChild($postionController);
            }

            $wrapper.setAttribute('style', `display: flex; padding: 5px 0;`);
            $wrapper.appendChild($container);

            $container.setAttribute('style', `${style}`);
            $container.appendChild($img);

            Object.entries(node.attrs).forEach(([key, value]) => {
                if (value === undefined || value === null) {
                    return;
                }
                $img.setAttribute(key, value);
            })

            if (!editable) return { dom: $img };

            const isMobile = document.documentElement.clientWidth < 768;
            const dotPosition = isMobile ? '-8px' : '-4px';

            const dotsPosition = [
                `top: ${dotPosition}; left: ${dotPosition}; cursor: nwse-resize;`,
                `top: ${dotPosition}; right: ${dotPosition}; cursor: nesw-resize;`,
                `bottom: ${dotPosition}; left: ${dotPosition}; cursor: nesw-resize;`,
                `bottom: ${dotPosition}; right: ${dotPosition}; cursor: nwse-resize;`,
            ];

            let isResizing = false;
            let startX, startWidth;

            $container.addEventListener('click', (e) => {
                var _a;

                //remove remaining dots and position controller
                const isMobile = document.documentElement.clientWidth < 768;

                isMobile && ((_a = document.querySelector('.ProseMirror-focused')) === null || _a === void 0 ? void 0 : _a.blur());

                if ($container.childElementCount > 3) {
                    for (let i = 0; i < 5; i++) {
                        $container.removeChild($container.lastChild);
                    }
                }

                paintPositionContoller();

                $container.setAttribute('style', `position: relative; border: 1px dashed #6C6C6C; ${style} cursor: pointer;`);

                Array.from({ length: 4 }, (_, index) => {
                    const $dot = document.createElement('div');
                    $dot.setAttribute('style', `position: absolute; width: ${isMobile ? 16 : 9}px; height: ${isMobile ? 16 : 9}px; border: 1.5px solid #6C6C6C; border-radius: 50%; ${dotsPosition[index]}`);
                    $dot.addEventListener('mousedown', (e) => {
                        e.preventDefault();

                        isResizing = true;
                        startX = e.clientX;
                        startWidth = $container.offsetWidth;

                        const onMouseMove = (e) => {
                            if (!isResizing) return;

                            const deltaX = index % 2 === 0 ? -(e.clientX - startX) : e.clientX - startX;
                            const newWidth = startWidth + deltaX;
                            if(newWidth < 50 || newWidth > 500) return;
                            $container.style.width = newWidth + 'px';
                            $img.style.width = newWidth + 'px';
                        }

                        const onMouseUp = () => {
                            if (isResizing) {
                                isResizing = false;
                            }

                            dispatchNodeView();

                            document.removeEventListener('mousemove', onMouseMove);
                            document.removeEventListener('mouseup', onMouseUp);
                        }
                        document.addEventListener('mousemove', onMouseMove);
                        document.addEventListener('mouseup', onMouseUp);
                    });
                    $dot.addEventListener('touchstart', (e) => {
                        e.cancelable && e.preventDefault();
                        isResizing = true;
                        startX = e.touches[0].clientX;
                        startWidth = $container.offsetWidth;

                        const onTouchMove = (e) => {
                            if (!isResizing) return;

                            const deltaX = index % 2 === 0
                                ? -(e.touches[0].clientX - startX)
                                : e.touches[0].clientX - startX;
                            const newWidth = startWidth + deltaX;
                            $container.style.width = newWidth + 'px';
                            $img.style.width = newWidth + 'px';
                        };

                        const onTouchEnd = () => {
                            if (isResizing) {
                                isResizing = false;
                            }

                            dispatchNodeView();

                            document.removeEventListener('touchmove', onTouchMove);
                            document.removeEventListener('touchend', onTouchEnd);

                        };
                        document.addEventListener('touchmove', onTouchMove);
                        document.addEventListener('touchend', onTouchEnd);
                    }, { passive: false });

                    $container.appendChild($dot);
                })
            })

            document.addEventListener('click', (e) => {
                const $target = e.target;
                const isClickInside = $container.contains($target) || $target.style.cssText === iconStyle;

                if (!isClickInside) {
                    const containerStyle = $container.getAttribute('style');
                    const newStyle = containerStyle === null || containerStyle === void 0 ? void 0 : containerStyle.replace('border: 1px dashed #6C6C6C;', '');
                    $container.setAttribute('style', newStyle);

                    if ($container.childElementCount > 3) {
                        for (let i = 0; i < 5; i++) {
                            $container.removeChild($container.lastChild);
                        }

                    }
                }
            })
            return {
                dom: $wrapper,
            };
        }
    },
    renderHTML({ HTMLAttributes }) {
        return ['div', {style: 'display: flex; padding: 5px 0;'}, 
            ['img', {src: HTMLAttributes.src, style: `${HTMLAttributes.style}`}]
        ];
    },
    parseHTML() {
        return [{
            tag: 'div',
            getAttrs: (element) => {
                const img = element.querySelector('img');
                return img ? {
                    src: img.getAttribute('src'),
                    alt: img.getAttribute('alt'),
                    title: img.getAttribute('title'),
                    style: img.getAttribute('style'),
                } : false;
            }
        }]
    },
    addCommands(){
        return {
            InsertImage: (options) => ({ commands }) => {
                const { src, alt, title } = options
                const defaultStyle = "width: 300px"
                return commands
                .insertContent({
                    type: 'image',
                    attrs: {
                        src,
                        alt,
                        title,
                        style: defaultStyle,
                    }
                }
                )
            }
        }
    },
    addInputRules() {
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
                        title,
                        style: "width: 300px",
                    }
                }
            })
        ]
    }
})

export {
    CustomImage,
    CustomHeading,
    SearchSelBackground,
    TextStyleExtends,
}

export default {
    
}