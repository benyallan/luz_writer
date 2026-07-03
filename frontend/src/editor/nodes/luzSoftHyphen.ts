import {Node, mergeAttributes} from '@tiptap/core'

// Ponto de quebra de palavra sugerido pelo autor (seção 6). Puramente
// decorativo no editor (glifo ‧); vira \- na exportação.
export const LuzSoftHyphen = Node.create({
  name: 'luzSoftHyphen',
  group: 'inline',
  inline: true,
  atom: true,
  selectable: false,

  parseHTML() {
    return [{tag: 'span[data-luz-node="luzSoftHyphen"]'}]
  },

  renderHTML({HTMLAttributes}) {
    return ['span', mergeAttributes(HTMLAttributes, {'data-luz-node': 'luzSoftHyphen', class: 'luz-soft-hyphen'}), '‧']
  },
})
