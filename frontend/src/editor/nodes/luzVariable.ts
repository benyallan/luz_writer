import {Node, mergeAttributes} from '@tiptap/core'
import {VueNodeViewRenderer} from '@tiptap/vue-3'
import LuzVariableView from '../nodeviews/LuzVariableView.vue'

export const LuzVariable = Node.create({
  name: 'luzVariable',
  group: 'inline',
  inline: true,
  atom: true,

  addAttributes() {
    return {
      name: {default: ''},
    }
  },

  parseHTML() {
    return [{tag: 'span[data-luz-node="luzVariable"]'}]
  },

  renderHTML({HTMLAttributes}) {
    return ['span', mergeAttributes(HTMLAttributes, {'data-luz-node': 'luzVariable'})]
  },

  addNodeView() {
    return VueNodeViewRenderer(LuzVariableView)
  },
})
