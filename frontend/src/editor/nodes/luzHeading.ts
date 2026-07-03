import {Node, mergeAttributes} from '@tiptap/core'
import {VueNodeViewRenderer} from '@tiptap/vue-3'
import HeadingNodeView from '../nodeviews/HeadingNodeView.vue'

// Fábrica compartilhada por luzChapter/luzSection/luzSubsection (seção 6) —
// os três têm o mesmo contrato de atributos e comportamento, diferindo só no
// nome/nível.
function createLuzHeading(name: string) {
  return Node.create({
    name,
    group: 'block',
    content: 'inline*',
    defining: true,

    addAttributes() {
      return {
        numbered: {default: true},
        includeInToc: {default: true},
      }
    },

    parseHTML() {
      return [{tag: `div[data-luz-node="${name}"]`}]
    },

    renderHTML({HTMLAttributes}) {
      return ['div', mergeAttributes(HTMLAttributes, {'data-luz-node': name}), 0]
    },

    addNodeView() {
      return VueNodeViewRenderer(HeadingNodeView)
    },
  })
}

export const luzChapter = createLuzHeading('luzChapter')
export const luzSection = createLuzHeading('luzSection')
export const luzSubsection = createLuzHeading('luzSubsection')
