import StarterKit from '@tiptap/starter-kit'
import Underline from '@tiptap/extension-underline'
import {LuzParagraph} from './nodes/paragraph'
import {luzChapter, luzSection, luzSubsection} from './nodes/luzHeading'
import {LuzImage} from './nodes/luzImage'
import {LuzFootnote} from './nodes/luzFootnote'
import {LuzSoftHyphen} from './nodes/luzSoftHyphen'
import {LuzVariable} from './nodes/luzVariable'
import {LuzInlineQuote} from './marks/luzInlineQuote'
import {LuzLang} from './marks/luzLang'
import {LuzCustomStyle} from './marks/luzCustomStyle'

// A interface descarta headings genéricos (H1/H2) e demais elementos fora
// do escopo do MVP (seção "Fora de escopo") em favor dos nós semânticos
// definidos na seção 6.
export function buildExtensions() {
  return [
    StarterKit.configure({
      heading: false,
      paragraph: false,
      codeBlock: false,
      code: false,
      horizontalRule: false,
    }),
    LuzParagraph,
    Underline,
    luzChapter,
    luzSection,
    luzSubsection,
    LuzImage,
    LuzFootnote,
    LuzSoftHyphen,
    LuzVariable,
    LuzInlineQuote,
    LuzLang,
    LuzCustomStyle,
  ]
}
