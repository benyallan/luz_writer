import Paragraph from '@tiptap/extension-paragraph'

export type ParagraphAlign = 'justify' | 'left' | 'center' | 'right'

// Estende o parágrafo padrão do Tiptap com o atributo align (seção 6,
// tabela de nós de bloco). "justify" é o único caso em que o visual do
// editor reflete literalmente a exportação (seção 6.0) — por isso não
// precisa de estilo inline, só os demais alinhamentos.
export const LuzParagraph = Paragraph.extend({
  addAttributes() {
    return {
      align: {
        default: 'justify' as ParagraphAlign,
        parseHTML: element => (element.getAttribute('data-align') as ParagraphAlign) || 'justify',
        renderHTML: attributes => {
          if (!attributes.align || attributes.align === 'justify') return {}
          return {
            'data-align': attributes.align,
            style: `text-align: ${attributes.align}`,
          }
        },
      },
    }
  },
})
