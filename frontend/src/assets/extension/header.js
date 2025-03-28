/* 
@Author: Hughie
@CreateTime: 2025-3-27
@LastEditors: Hughie
@LastEditTime: 2025-3-28
@Description: the heading extension of tiptap
*/

import Heading from "@tiptap/extension-heading"
import { textblockTypeInputRule } from "@tiptap/vue-3"

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
                    return commands.setNode('paragraph', {"textAlign": "left"})
                }
                if (!this.options.levels.includes(level)){
                    return false;
                }
                let uuid = new Date().getTime().toString(36) + Math.random().toString(36).slice(2,9);
                return commands.setNode(this.name, {level: level, id: `tap-heading-${uuid}`});
            },
            toggleHeader: level => ({ commands }) => {
                if(level == 0) {
                    return commands.toggleNode('heading' ,'paragraph', {"textAlign": "left"})
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

export {
    CustomHeading
}

export default CustomHeading