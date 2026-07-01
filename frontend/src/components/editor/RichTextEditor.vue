<script setup>
import { watch } from 'vue'
import { useEditor, EditorContent } from '@tiptap/vue-3'
import StarterKit from '@tiptap/starter-kit'
import Typography from '@tiptap/extension-typography'
import Placeholder from '@tiptap/extension-placeholder'
import {
  Bold, Italic, Strikethrough,
  Heading1, Heading2, Heading3,
  List, ListOrdered, Quote,
} from '@lucide/vue'
import { useEditorStore } from '@/stores/editor'

const store = useEditorStore()

const EMPTY_DOC = { type: 'doc', content: [{ type: 'paragraph' }] }

function parseContent(raw) {
  if (!raw?.trim()) return EMPTY_DOC
  try { return JSON.parse(raw) } catch { return EMPTY_DOC }
}

function normalizeEmpty(editor) {
  if (!store.content.trim()) {
    const canonical = JSON.stringify(editor.getJSON())
    store.content = canonical
    store.savedContent = canonical
  }
}

const tiptap = useEditor({
  extensions: [
    StarterKit,
    Typography,
    Placeholder.configure({ placeholder: 'Comece a escrever...' }),
  ],
  content: parseContent(store.content),
  onCreate({ editor }) {
    normalizeEmpty(editor)
  },
  onUpdate({ editor }) {
    store.content = JSON.stringify(editor.getJSON())
  },
})

// Reload Tiptap when a different .luztxt file is opened
watch(() => store.filePath, () => {
  tiptap.value?.commands.setContent(parseContent(store.content), false)
  normalizeEmpty(tiptap.value)
})
</script>

<template>
  <div class="rich-editor">
    <!-- Toolbar -->
    <div class="toolbar">
      <div class="toolbar-group">
        <button
          class="tool-btn"
          :class="{ active: tiptap?.isActive('bold') }"
          title="Negrito (Ctrl+B)"
          @click="tiptap?.chain().focus().toggleBold().run()"
        ><Bold :size="15" /></button>
        <button
          class="tool-btn"
          :class="{ active: tiptap?.isActive('italic') }"
          title="Itálico (Ctrl+I)"
          @click="tiptap?.chain().focus().toggleItalic().run()"
        ><Italic :size="15" /></button>
        <button
          class="tool-btn"
          :class="{ active: tiptap?.isActive('strike') }"
          title="Tachado"
          @click="tiptap?.chain().focus().toggleStrike().run()"
        ><Strikethrough :size="15" /></button>
      </div>

      <div class="toolbar-sep" />

      <div class="toolbar-group">
        <button
          class="tool-btn"
          :class="{ active: tiptap?.isActive('heading', { level: 1 }) }"
          title="Título 1"
          @click="tiptap?.chain().focus().toggleHeading({ level: 1 }).run()"
        ><Heading1 :size="15" /></button>
        <button
          class="tool-btn"
          :class="{ active: tiptap?.isActive('heading', { level: 2 }) }"
          title="Título 2"
          @click="tiptap?.chain().focus().toggleHeading({ level: 2 }).run()"
        ><Heading2 :size="15" /></button>
        <button
          class="tool-btn"
          :class="{ active: tiptap?.isActive('heading', { level: 3 }) }"
          title="Título 3"
          @click="tiptap?.chain().focus().toggleHeading({ level: 3 }).run()"
        ><Heading3 :size="15" /></button>
      </div>

      <div class="toolbar-sep" />

      <div class="toolbar-group">
        <button
          class="tool-btn"
          :class="{ active: tiptap?.isActive('bulletList') }"
          title="Lista com marcadores"
          @click="tiptap?.chain().focus().toggleBulletList().run()"
        ><List :size="15" /></button>
        <button
          class="tool-btn"
          :class="{ active: tiptap?.isActive('orderedList') }"
          title="Lista numerada"
          @click="tiptap?.chain().focus().toggleOrderedList().run()"
        ><ListOrdered :size="15" /></button>
        <button
          class="tool-btn"
          :class="{ active: tiptap?.isActive('blockquote') }"
          title="Citação"
          @click="tiptap?.chain().focus().toggleBlockquote().run()"
        ><Quote :size="15" /></button>
      </div>
    </div>

    <!-- Content area -->
    <div class="content-area">
      <EditorContent :editor="tiptap" class="paper-content" />
    </div>
  </div>
</template>

<style scoped>
.rich-editor {
  display: flex;
  flex-direction: column;
  height: 100%;
}

/* ── Toolbar ── */
.toolbar {
  display: flex;
  align-items: center;
  gap: 2px;
  padding: 4px 8px;
  background: #2d2d2d;
  border-bottom: 1px solid var(--color-border);
  flex-shrink: 0;
}

.toolbar-group {
  display: flex;
  gap: 1px;
}

.toolbar-sep {
  width: 1px;
  height: 18px;
  background: var(--color-border);
  margin: 0 4px;
}

.tool-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 28px;
  height: 28px;
  border-radius: 3px;
  color: var(--color-text-muted);
  transition: background 0.1s, color 0.1s;
}

.tool-btn:hover {
  background: rgba(255, 255, 255, 0.08);
  color: var(--color-text);
}

.tool-btn.active {
  background: rgba(14, 99, 156, 0.35);
  color: #4fc3f7;
}

/* ── Content area ── */
.content-area {
  flex: 1;
  overflow-y: auto;
  background: #faf9f7;
  padding: 40px 48px;
}
</style>

<!-- Global: styles that apply inside the ProseMirror content area -->
<style>
.paper-content .ProseMirror {
  outline: none;
  font-family: Georgia, 'Times New Roman', 'Palatino Linotype', serif;
  font-size: 17px;
  line-height: 1.85;
  color: #1a1a1a;
  width: 100%;
  min-height: calc(100vh - 120px);
}

.paper-content .ProseMirror > * + * {
  margin-top: 0.75em;
}

.paper-content .ProseMirror h1 {
  font-size: 2em;
  line-height: 1.2;
  margin-top: 1.4em;
  margin-bottom: 0.4em;
  font-weight: 700;
}
.paper-content .ProseMirror h2 {
  font-size: 1.45em;
  line-height: 1.3;
  margin-top: 1.2em;
  margin-bottom: 0.3em;
  font-weight: 600;
}
.paper-content .ProseMirror h3 {
  font-size: 1.15em;
  line-height: 1.4;
  margin-top: 1em;
  margin-bottom: 0.2em;
  font-weight: 600;
}

.paper-content .ProseMirror p { margin: 0; }
.paper-content .ProseMirror p + p { text-indent: 1.5em; }

.paper-content .ProseMirror ul,
.paper-content .ProseMirror ol {
  padding-left: 1.6em;
}

.paper-content .ProseMirror blockquote {
  border-left: 3px solid #c0b090;
  padding-left: 1em;
  color: #555;
  font-style: italic;
  margin: 1em 0;
}

.paper-content .ProseMirror strong { font-weight: 700; }
.paper-content .ProseMirror em     { font-style: italic; }
.paper-content .ProseMirror s      { text-decoration: line-through; }

.paper-content .ProseMirror hr {
  border: none;
  border-top: 1px solid #d0c8bc;
  margin: 2em 0;
}

/* Placeholder */
.paper-content .ProseMirror p.is-editor-empty:first-child::before {
  content: attr(data-placeholder);
  float: left;
  color: #b0a898;
  pointer-events: none;
  height: 0;
}
</style>
