<script setup>
import { ref, watch } from 'vue'
import {
  DialogRoot,
  DialogPortal,
  DialogOverlay,
  DialogContent,
  DialogTitle,
  DialogClose,
} from 'reka-ui'
import { FolderOpen, X } from '@lucide/vue'
import { GetHomeDir, OpenFolder, CreateProject } from '@wails/go/main/App'
import { useWorkspaceStore } from '@/stores/workspace'

const workspace = useWorkspaceStore()

const projectName = ref('')
const parentPath = ref('')
const errorMsg = ref('')
const isSubmitting = ref(false)

watch(
  () => workspace.isNewProjectDialogOpen,
  async (open) => {
    if (open && !parentPath.value) {
      parentPath.value = await GetHomeDir()
    }
    if (!open) {
      projectName.value = ''
      errorMsg.value = ''
      isSubmitting.value = false
    }
  }
)

async function chooseParent() {
  const path = await OpenFolder()
  if (path) parentPath.value = path
}

async function submit() {
  errorMsg.value = ''
  if (!projectName.value.trim()) {
    errorMsg.value = 'Digite um nome para o projeto.'
    return
  }
  isSubmitting.value = true
  try {
    const projectPath = await CreateProject(projectName.value.trim(), parentPath.value)
    await workspace.openAt(projectPath)
    workspace.isNewProjectDialogOpen = false
  } catch (e) {
    errorMsg.value = String(e)
  } finally {
    isSubmitting.value = false
  }
}
</script>

<template>
  <DialogRoot v-model:open="workspace.isNewProjectDialogOpen">
    <DialogPortal>
      <DialogOverlay class="dialog-overlay" />
      <DialogContent class="dialog-box" @open-auto-focus.prevent>
        <div class="dialog-header">
          <DialogTitle class="dialog-title">Novo Projeto</DialogTitle>
          <DialogClose class="dialog-close-btn">
            <X :size="16" />
          </DialogClose>
        </div>

        <form class="dialog-body" @submit.prevent="submit">
          <label class="field-label" for="project-name">Nome do projeto</label>
          <input
            id="project-name"
            v-model="projectName"
            class="field-input"
            type="text"
            placeholder="Meu Livro"
            autocomplete="off"
            autofocus
            @keydown.enter.prevent="submit"
          />

          <label class="field-label">Localização</label>
          <div class="location-row">
            <span class="location-path" :title="parentPath">{{ parentPath }}</span>
            <button type="button" class="location-btn" @click="chooseParent">
              <FolderOpen :size="14" />
              Alterar...
            </button>
          </div>

          <p v-if="errorMsg" class="error-msg">{{ errorMsg }}</p>

          <div class="dialog-footer">
            <DialogClose class="btn-secondary">Cancelar</DialogClose>
            <button
              type="submit"
              class="btn-primary"
              :disabled="isSubmitting"
            >
              {{ isSubmitting ? 'Criando...' : 'Criar projeto' }}
            </button>
          </div>
        </form>
      </DialogContent>
    </DialogPortal>
  </DialogRoot>
</template>

<style scoped>
.dialog-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px 20px 0;
}

.dialog-title {
  font-size: 14px;
  font-weight: 600;
  color: var(--color-text);
}

.dialog-close-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 24px;
  height: 24px;
  border-radius: 4px;
  color: var(--color-text-muted);
  transition: background 0.1s, color 0.1s;
}

.dialog-close-btn:hover {
  background: var(--color-hover);
  color: var(--color-text);
}

.dialog-body {
  padding: 16px 20px 20px;
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.field-label {
  font-size: 11px;
  font-weight: 600;
  color: var(--color-text-muted);
  text-transform: uppercase;
  letter-spacing: 0.06em;
  margin-top: 8px;
}

.field-label:first-child {
  margin-top: 0;
}

.field-input {
  width: 100%;
  background: var(--color-bg);
  border: 1px solid var(--color-border);
  border-radius: 2px;
  padding: 6px 10px;
  font-size: 13px;
  color: var(--color-text);
  font-family: inherit;
  outline: none;
  transition: border-color 0.15s;
}

.field-input:focus {
  border-color: var(--color-accent);
}

.location-row {
  display: flex;
  align-items: center;
  gap: 8px;
}

.location-path {
  flex: 1;
  font-size: 12px;
  color: var(--color-text-muted);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  direction: rtl;
  text-align: left;
}

.location-btn {
  display: flex;
  align-items: center;
  gap: 5px;
  padding: 4px 10px;
  font-size: 12px;
  color: var(--color-text);
  border: 1px solid var(--color-border);
  border-radius: 2px;
  flex-shrink: 0;
  transition: background 0.1s;
}

.location-btn:hover {
  background: var(--color-hover);
}

.error-msg {
  font-size: 12px;
  color: #f48771;
  margin-top: 4px;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 8px;
  margin-top: 16px;
}

.btn-secondary {
  padding: 5px 14px;
  font-size: 13px;
  color: var(--color-text);
  border: 1px solid var(--color-border);
  border-radius: 2px;
  transition: background 0.1s;
}

.btn-secondary:hover {
  background: var(--color-hover);
}

.btn-primary {
  padding: 5px 14px;
  font-size: 13px;
  background: var(--color-accent);
  color: #fff;
  border-radius: 2px;
  transition: background 0.1s;
}

.btn-primary:hover:not(:disabled) {
  background: #1177bb;
}

.btn-primary:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}
</style>

<style>
.dialog-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.5);
  z-index: 100;
  animation: fadeIn 0.15s ease-out;
}

.dialog-box {
  position: fixed;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  width: 440px;
  background: #252526;
  border: 1px solid var(--color-border);
  border-radius: 4px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.5);
  z-index: 101;
  animation: slideIn 0.15s ease-out;
}

@keyframes slideIn {
  from { opacity: 0; transform: translate(-50%, -52%); }
  to   { opacity: 1; transform: translate(-50%, -50%); }
}
</style>
