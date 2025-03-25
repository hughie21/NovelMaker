import { Node } from '@tiptap/core';

const RubyExtension = Node.create({
    name: 'ruby',
    inline: true,         
    group: 'inline', 
    selectable: true,
    draggable: false,
    addOptions() {
        return {
        HTMLAttributes: {},
        }
    },
    addAttributes() {
        return {
        aboveText: {
            default: null,
            parseHTML: element => element.querySelector('rt')?.textContent,
            renderHTML: attributes => ({ 'data-above': attributes.aboveText }),
        },
        belowText: {
            default: null,
            parseHTML: element => element.childNodes[0]?.textContent,
            renderHTML: attributes => ({ 'data-below': attributes.belowText }),
        },
        }
    },
    parseHTML() {
        return [{
        tag: 'ruby',
        getAttrs: node => ({
            aboveText: node.querySelector('rt')?.textContent,
            belowText: node.childNodes[0]?.textContent,
        }),
        }]
    },
    renderHTML({ node, HTMLAttributes }) {
        return ['ruby', HTMLAttributes,
        node.attrs.belowText,
        ['rt', {}, node.attrs.aboveText],
        ]
    },
    addCommands() {
        return {
        setRuby: (aboveText, belowText) => ({ commands }) => {
            return commands.insertContent({
            type: this.name,
            attrs: { aboveText, belowText },
            })
        },
        }
    },
})

export default RubyExtension;
export { RubyExtension }
