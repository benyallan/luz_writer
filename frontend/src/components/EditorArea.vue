<script setup lang="ts">
import {watch, onBeforeUnmount} from 'vue'
import {useEditor, EditorContent} from '@tiptap/vue-3'
import {buildExtensions} from '../editor/extensions'
import EditorToolbar from './EditorToolbar.vue'
import {useEditorStore} from '../stores/editor'
import {useUiStore} from '../stores/ui'

const editorStore = useEditorStore()
const ui = useUiStore()

const editor = useEditor({
  extensions: buildExtensions(),
  content: {type: 'doc', content: [{type: 'paragraph'}]},
  editable: false,
  onUpdate: () => {
    editorStore.scheduleSave(() => editor.value!.getJSON())
  },
})

watch(
  () => editorStore.activeChapterId,
  () => {
    if (!editor.value) return
    if (editorStore.activeChapterId && editorStore.initialContent) {
      editor.value.commands.setContent(editorStore.initialContent, false)
      editor.value.setEditable(true)
    } else {
      editor.value.commands.clearContent(false)
      editor.value.setEditable(false)
    }
  },
)

onBeforeUnmount(() => {
  editor.value?.destroy()
})
</script>

<template>
  <main class="editor-area">
    <div v-if="editor && editorStore.activeChapterId" class="editor-area__toolbar-row">
      <EditorToolbar :editor="editor" :role="editorStore.activeChapterMeta?.role ?? null" />
      <button
        class="editor-area__page-settings"
        type="button"
        title="Configurações desta página"
        @click="ui.openPageOverrides(editorStore.activeChapterId!)"
      >
        ⚙
      </button>
    </div>
    <div class="editor-area__scroll">
      <div class="editor-area__flow">
        <p v-if="!editorStore.activeChapterId" class="editor-area__placeholder">
          Abra ou crie um capítulo para começar a escrever.
        </p>
        <EditorContent v-else :editor="editor" />
      </div>
    </div>
  </main>
</template>

<style scoped>
.editor-area {
  flex: 1 1 auto;
  min-width: 0;
  display: flex;
  flex-direction: column;
  background: var(--luz-bg-editor);
}

.editor-area__toolbar-row {
  display: flex;
  align-items: center;
}

.editor-area__toolbar-row > :first-child {
  flex: 1;
  min-width: 0;
}

.editor-area__page-settings {
  border: none;
  border-bottom: 1px solid var(--luz-border);
  background: var(--luz-bg-sidebar);
  color: var(--luz-fg-muted);
  align-self: stretch;
  padding: 0 12px;
  cursor: pointer;
}

.editor-area__page-settings:hover {
  color: var(--luz-fg);
  background: var(--luz-bg-hover);
}

.editor-area__scroll {
  flex: 1 1 auto;
  overflow-y: auto;
  display: flex;
  justify-content: center;
}

.editor-area__flow {
  width: 100%;
  max-width: 800px;
  padding: 48px 24px;
}

.editor-area__placeholder {
  color: var(--luz-fg-muted);
  text-align: center;
  margin-top: 20vh;
}
</style>
