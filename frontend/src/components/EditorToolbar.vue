<script setup lang="ts">
import {computed} from 'vue'
import type {Editor} from '@tiptap/vue-3'
import {
  DropdownMenuRoot,
  DropdownMenuTrigger,
  DropdownMenuPortal,
  DropdownMenuContent,
  DropdownMenuItem,
} from 'reka-ui'
import {LUZ_LANG_OPTIONS} from '../editor/marks/luzLang'
import {useWorkspaceStore} from '../stores/workspace'
import {usePluginsStore} from '../stores/plugins'
import type {model} from '../../wailsjs/go/models'

const props = defineProps<{
  editor: Editor
  role: string | null
}>()

const workspace = useWorkspaceStore()
const pluginsStore = usePluginsStore()

// Toolbar sensível ao role (seção 5.5): páginas especiais escondem títulos
// estruturais; dedicatória/epígrafe escondem também listas e imagens.
const headingsVisible = computed(() => props.role === 'chapter' || props.role === 'appendix' || props.role === null)
const listsAndImagesVisible = computed(() => props.role !== 'dedication' && props.role !== 'epigraph')

function setHeading(name: 'luzChapter' | 'luzSection' | 'luzSubsection') {
  props.editor.chain().focus().setNode(name, {numbered: true, includeInToc: true}).run()
}

function setParagraph() {
  props.editor.chain().focus().setNode('paragraph').run()
}

function setAlign(align: 'justify' | 'left' | 'center' | 'right') {
  props.editor.chain().focus().updateAttributes('paragraph', {align}).run()
}

async function insertImage() {
  const App = await import('../../wailsjs/go/main/App')
  const src = await App.ImportImage()
  if (!src) return
  props.editor.chain().focus().insertContent({type: 'luzImage', attrs: {src, caption: '', width: 80}}).run()
}

function insertFootnote() {
  props.editor.chain().focus().insertContent({type: 'luzFootnote', attrs: {number: 1, text: ''}}).run()
}

function insertSoftHyphen() {
  props.editor.chain().focus().insertContent({type: 'luzSoftHyphen'}).run()
}

// O Reka UI devolve o foco ao botão-gatilho do dropdown ao fechá-lo, depois
// do @select disparar — adiar com setTimeout(0) garante que o .focus() do
// Tiptap rode por último e vença essa corrida, preservando o cursor no
// editor para o usuário continuar digitando.
function afterDropdownClose(fn: () => void) {
  setTimeout(fn, 0)
}

function insertVariable(name: string) {
  afterDropdownClose(() => {
    props.editor.chain().focus().insertContent({type: 'luzVariable', attrs: {name}}).run()
  })
}

function setLang(lang: string) {
  afterDropdownClose(() => {
    props.editor.chain().focus().setMark('luzLang', {lang}).run()
  })
}

function clearLang() {
  afterDropdownClose(() => {
    props.editor.chain().focus().unsetMark('luzLang').run()
  })
}

function setCustomStyle(styleId: string) {
  afterDropdownClose(() => {
    props.editor.chain().focus().setMark('luzCustomStyle', {styleId}).run()
  })
}

function clearCustomStyle() {
  afterDropdownClose(() => {
    props.editor.chain().focus().unsetMark('luzCustomStyle').run()
  })
}

const variables = computed<model.Variable[]>(() => workspace.project?.variables ?? [])

// Gating de toolbar (Etapa 4, seção 8.3): idioma, hífen sugerido e estilos só
// aparecem com os respectivos plugins ativos.
const languagesEnabled = computed(() => pluginsStore.isPluginEnabled('languages'))
const hyphenationEnabled = computed(() => pluginsStore.isPluginEnabled('hyphenation'))
const customStylesEnabled = computed(() => pluginsStore.isPluginEnabled('customStyles'))
</script>

<template>
  <div class="editor-toolbar">
    <button type="button" class="toolbar-btn" :class="{active: editor.isActive('bold')}" title="Negrito" @click="editor.chain().focus().toggleBold().run()">B</button>
    <button type="button" class="toolbar-btn" :class="{active: editor.isActive('italic')}" title="Itálico" @click="editor.chain().focus().toggleItalic().run()">I</button>
    <button type="button" class="toolbar-btn" :class="{active: editor.isActive('underline')}" title="Sublinhado" @click="editor.chain().focus().toggleUnderline().run()">U</button>
    <button type="button" class="toolbar-btn" :class="{active: editor.isActive('strike')}" title="Tachado" @click="editor.chain().focus().toggleStrike().run()">S</button>
    <button type="button" class="toolbar-btn" :class="{active: editor.isActive('luzInlineQuote')}" title="Citação inline" @click="editor.chain().focus().toggleMark('luzInlineQuote').run()">"</button>

    <span class="toolbar-divider" />

    <template v-if="headingsVisible">
      <button type="button" class="toolbar-btn" :class="{active: editor.isActive('luzChapter')}" title="Capítulo" @click="setHeading('luzChapter')">Cap</button>
      <button type="button" class="toolbar-btn" :class="{active: editor.isActive('luzSection')}" title="Seção" @click="setHeading('luzSection')">Seç</button>
      <button type="button" class="toolbar-btn" :class="{active: editor.isActive('luzSubsection')}" title="Subseção" @click="setHeading('luzSubsection')">Sub</button>
      <button type="button" class="toolbar-btn" title="Parágrafo" @click="setParagraph">¶</button>
      <span class="toolbar-divider" />
    </template>

    <button type="button" class="toolbar-btn" title="Justificado" @click="setAlign('justify')">≡</button>
    <button type="button" class="toolbar-btn" title="Esquerda" @click="setAlign('left')">⇤</button>
    <button type="button" class="toolbar-btn" title="Centralizado" @click="setAlign('center')">↔</button>
    <button type="button" class="toolbar-btn" title="Direita" @click="setAlign('right')">⇥</button>

    <span class="toolbar-divider" />

    <template v-if="listsAndImagesVisible">
      <button type="button" class="toolbar-btn" :class="{active: editor.isActive('bulletList')}" title="Lista" @click="editor.chain().focus().toggleBulletList().run()">•</button>
      <button type="button" class="toolbar-btn" :class="{active: editor.isActive('orderedList')}" title="Lista numerada" @click="editor.chain().focus().toggleOrderedList().run()">1.</button>
      <button type="button" class="toolbar-btn" :class="{active: editor.isActive('blockquote')}" title="Citação em bloco" @click="editor.chain().focus().toggleBlockquote().run()">❝</button>
      <button type="button" class="toolbar-btn" title="Imagem" @click="insertImage">🖼</button>
      <span class="toolbar-divider" />
    </template>

    <button type="button" class="toolbar-btn" title="Nota de rodapé" @click="insertFootnote">[n]</button>
    <button v-if="hyphenationEnabled" type="button" class="toolbar-btn" title="Hífen sugerido" @click="insertSoftHyphen">‧</button>

    <DropdownMenuRoot v-if="languagesEnabled">
      <DropdownMenuTrigger class="toolbar-btn" title="Idioma">🌐</DropdownMenuTrigger>
      <DropdownMenuPortal>
        <DropdownMenuContent class="toolbar-menu">
          <DropdownMenuItem v-for="lang in LUZ_LANG_OPTIONS" :key="lang.value" class="toolbar-menu__item" @select="setLang(lang.value)">
            {{ lang.label }}
          </DropdownMenuItem>
          <DropdownMenuItem class="toolbar-menu__item" @select="clearLang">Remover marcação</DropdownMenuItem>
        </DropdownMenuContent>
      </DropdownMenuPortal>
    </DropdownMenuRoot>

    <DropdownMenuRoot v-if="customStylesEnabled && pluginsStore.styles.length">
      <DropdownMenuTrigger class="toolbar-btn" title="Estilos">Aa</DropdownMenuTrigger>
      <DropdownMenuPortal>
        <DropdownMenuContent class="toolbar-menu">
          <DropdownMenuItem v-for="s in pluginsStore.styles" :key="s.id" class="toolbar-menu__item" @select="setCustomStyle(s.id)">
            {{ s.name }}
          </DropdownMenuItem>
          <DropdownMenuItem class="toolbar-menu__item" @select="clearCustomStyle">Remover estilo</DropdownMenuItem>
        </DropdownMenuContent>
      </DropdownMenuPortal>
    </DropdownMenuRoot>

    <DropdownMenuRoot v-if="variables.length">
      <DropdownMenuTrigger class="toolbar-btn" title="Inserir variável">{{ '{x}' }}</DropdownMenuTrigger>
      <DropdownMenuPortal>
        <DropdownMenuContent class="toolbar-menu">
          <DropdownMenuItem v-for="v in variables" :key="v.name" class="toolbar-menu__item" @select="insertVariable(v.name)">
            {{ v.name }}
          </DropdownMenuItem>
        </DropdownMenuContent>
      </DropdownMenuPortal>
    </DropdownMenuRoot>
  </div>
</template>

<style scoped>
.editor-toolbar {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 2px;
  padding: 6px 12px;
  border-bottom: 1px solid var(--luz-border);
  background: var(--luz-bg-sidebar);
}

.toolbar-btn {
  border: none;
  background: transparent;
  color: var(--luz-fg-muted);
  border-radius: 4px;
  min-width: 26px;
  height: 26px;
  padding: 0 6px;
  cursor: pointer;
  font-size: 0.8rem;
}

.toolbar-btn:hover {
  background: var(--luz-bg-hover);
  color: var(--luz-fg);
}

.toolbar-btn.active {
  background: var(--luz-bg-hover);
  color: var(--luz-fg);
  font-weight: 700;
}

.toolbar-divider {
  width: 1px;
  height: 18px;
  background: var(--luz-border);
  margin: 0 4px;
}

.toolbar-menu {
  background: var(--luz-bg-editor);
  border: 1px solid var(--luz-border);
  border-radius: 6px;
  padding: 4px;
  min-width: 160px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

.toolbar-menu__item {
  padding: 6px 8px;
  border-radius: 4px;
  cursor: pointer;
  font-size: 0.82rem;
  outline: none;
}

.toolbar-menu__item[data-highlighted] {
  background: var(--luz-bg-hover);
}
</style>
