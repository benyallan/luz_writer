import {Node, mergeAttributes} from '@tiptap/core'
import {Plugin, PluginKey} from '@tiptap/pm/state'
import {VueNodeViewRenderer} from '@tiptap/vue-3'
import LuzFootnoteView from '../nodeviews/LuzFootnoteView.vue'

export const LuzFootnote = Node.create({
  name: 'luzFootnote',
  group: 'inline',
  inline: true,
  atom: true,

  addAttributes() {
    return {
      number: {default: 1},
      text: {default: ''},
    }
  },

  parseHTML() {
    return [{tag: 'span[data-luz-node="luzFootnote"]'}]
  },

  renderHTML({HTMLAttributes}) {
    return ['span', mergeAttributes(HTMLAttributes, {'data-luz-node': 'luzFootnote'})]
  },

  addNodeView() {
    return VueNodeViewRenderer(LuzFootnoteView)
  },

  // Renumera as notas em ordem de documento a cada transação — o atributo
  // number é apenas de exibição no editor (o \footnote{} do LaTeX se
  // autonumera), mas precisa refletir a ordem atual para o marcador [n].
  addProseMirrorPlugins() {
    const typeName = this.name
    return [
      new Plugin({
        key: new PluginKey('luzFootnoteRenumber'),
        appendTransaction(transactions, _oldState, newState) {
          if (!transactions.some(tr => tr.docChanged)) return null

          let n = 0
          let changed = false
          const tr = newState.tr
          newState.doc.descendants((node, pos) => {
            if (node.type.name !== typeName) return
            n += 1
            if (node.attrs.number !== n) {
              tr.setNodeMarkup(pos, undefined, {...node.attrs, number: n})
              changed = true
            }
          })
          return changed ? tr : null
        },
      }),
    ]
  },
})
