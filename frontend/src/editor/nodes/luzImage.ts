import {Node, mergeAttributes} from '@tiptap/core'
import {VueNodeViewRenderer} from '@tiptap/vue-3'
import LuzImageView from '../nodeviews/LuzImageView.vue'

export const LuzImage = Node.create({
  name: 'luzImage',
  group: 'block',
  atom: true,
  draggable: true,

  addAttributes() {
    return {
      src: {default: null},
      caption: {default: ''},
      width: {default: 80},
    }
  },

  parseHTML() {
    return [{tag: 'div[data-luz-node="luzImage"]'}]
  },

  renderHTML({HTMLAttributes}) {
    return ['div', mergeAttributes(HTMLAttributes, {'data-luz-node': 'luzImage'})]
  },

  addNodeView() {
    return VueNodeViewRenderer(LuzImageView)
  },
})
