/* 
@Author: Hughie
@CreateTime: 2025-3-27
@LastEditors: Hughie
@LastEditTime: 2025-3-28
@Description: the text align extension of tiptap
*/

import Paragraph from '@tiptap/extension-paragraph'

const TextAlign = Paragraph.extend({
    priority: 1000,
    addAttributes() {
        return {
            ...this.parent?.(),
            textAlign: {
                default: "left",
            }
        }
    },
    addOptions() {
        return {
            ...this.parent?.()
        }
    },
    renderHTML({ HTMLAttributes }) {
        return ['p', { style: `text-align: ${HTMLAttributes.textAlign};` }, 0];
    },
    parseHTML() {
        return [{
            tag: 'p',
            getAttrs: node => {
                const textAlign = node.style.textAlign;
                return { textAlign };
            },
        }]
    },
    addCommands() {
        return {
            setTextAlign: align => ({ commands }) => {
                return commands.setNode('paragraph', {'textAlign': align});
            }
        }
    }
})

export default TextAlign

export {
    TextAlign
}