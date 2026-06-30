<script setup>
import { X } from '@lucide/vue'
import { useEditorStore } from '@/stores/editor'
import { useWorkspaceStore } from '@/stores/workspace'

const editor = useEditorStore()
const workspace = useWorkspaceStore()

function handleKeydown(e) {
  if ((e.ctrlKey || e.metaKey) && e.key === 's') {
    e.preventDefault()
    editor.save()
  }
}

function close() {
  workspace.setActiveFile(null)
  editor.close()
}
</script>

<template>
  <div class="file-editor">
    <!-- Tab bar -->
    <div class="tab-bar">
      <div class="tab">
        <span v-if="editor.isDirty" class="dirty-dot" title="Alterações não salvas">●</span>
        <span class="tab-name">{{ editor.fileName }}</span>
        <button class="tab-close" title="Fechar (sem salvar)" @click="close">
          <X :size="12" />
        </button>
      </div>
    </div>

    <!-- Editor body -->
    <div class="editor-body">
      <textarea
        v-model="editor.content"
        class="editor-textarea"
        spellcheck="false"
        autocomplete="off"
        autocorrect="off"
        autocapitalize="off"
        @keydown="handleKeydown"
      />
    </div>
  </div>
</template>

<style scoped>
.file-editor {
  display: flex;
  flex-direction: column;
  height: 100%;
  background: var(--color-bg);
}

/* ── Tab bar ── */
.tab-bar {
  height: 35px;
  background: #2d2d2d;
  border-bottom: 1px solid var(--color-border);
  display: flex;
  align-items: stretch;
  flex-shrink: 0;
}

.tab {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 0 8px 0 14px;
  border-right: 1px solid var(--color-border);
  border-top: 2px solid var(--color-accent);
  background: var(--color-bg);
  font-size: 13px;
  color: var(--color-text);
  max-width: 200px;
}

.dirty-dot {
  font-size: 14px;
  color: #e8c07d;
  flex-shrink: 0;
  line-height: 1;
}

.tab-name {
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  flex: 1;
  min-width: 0;
}

.tab-close {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 18px;
  height: 18px;
  border-radius: 3px;
  color: var(--color-text-muted);
  flex-shrink: 0;
  opacity: 0;
  transition: opacity 0.1s, background 0.1s, color 0.1s;
}

.tab:hover .tab-close {
  opacity: 1;
}

.tab-close:hover {
  background: rgba(255, 255, 255, 0.12);
  color: var(--color-text);
}

/* ── Editor body ── */
.editor-body {
  flex: 1;
  overflow: hidden;
}

.editor-textarea {
  width: 100%;
  height: 100%;
  padding: 20px 24px;
  background: var(--color-bg);
  color: var(--color-text);
  font-family: 'Cascadia Code', 'Fira Code', 'Consolas', 'Menlo', monospace;
  font-size: 13px;
  line-height: 1.65;
  border: none;
  outline: none;
  resize: none;
  box-sizing: border-box;
  tab-size: 4;
}
</style>
