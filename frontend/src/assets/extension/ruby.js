/* 
@Author: Hughie
@CreateTime: 2025-3-28
@LastEditors: Hughie
@LastEditTime: 2025-3-28
@Description: the above text extension of tiptap
*/

import { Node } from '@tiptap/core';

const RubyExtension = Node.create({
  name: 'ruby',
  inline: true,
  group: 'inline',
  selectable: true,
  draggable: false,
  content: 'text*',
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
    }
  },
  parseHTML() {
    return [{
      tag: 'ruby',
      getAttrs: node => ({
        aboveText: node.querySelector('rt')?.textContent,
      }),
    }]
  },
  renderHTML({ node, HTMLAttributes }) {
    return [
      'ruby', HTMLAttributes,
      ['rb', 0],
      ['rt', {}, node.attrs.aboveText],
    ]
  },
  addCommands() {
    return {
      setRuby: (aboveText, belowText) => ({ commands }) => {
        return commands.insertContent({
          type: this.name,
          attrs: { aboveText },
          content: [{ type: 'text', text: belowText }],
        })
      },
    }
  },
  addNodeView() {
    return ({ node }) => {
      const rubyBlock = document.createElement('ruby');
      const belowWrapper = document.createElement('rb');
      const aboveText = document.createElement('rt');

      belowWrapper.style.display = 'inline-block';
      belowWrapper.style.minWidth = '1ch'; 

      aboveText.textContent = node.attrs.aboveText || '';
      aboveText.contentEditable = false;

      rubyBlock.appendChild(belowWrapper);
      rubyBlock.appendChild(aboveText);

      return {
        dom: rubyBlock,
        contentDOM: belowWrapper,
        update(newNode) {
          if (newNode.type !== node.type) {
            return false;
          }
          if (newNode.attrs.aboveText !== node.attrs.aboveText) {
            aboveText.textContent = newNode.attrs.aboveText;
          }
          node = newNode;
          return true;
        },
        stopEvent: () => false,
      };
    }
  }
});

export default RubyExtension;
export { RubyExtension };
