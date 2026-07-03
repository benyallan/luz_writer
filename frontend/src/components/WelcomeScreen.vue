<script setup lang="ts">
import {ref} from 'vue'
import {
  DialogRoot,
  DialogPortal,
  DialogOverlay,
  DialogContent,
  DialogTitle,
  SelectRoot,
  SelectTrigger,
  SelectValue,
  SelectIcon,
  SelectPortal,
  SelectContent,
  SelectViewport,
  SelectItem,
  SelectItemText,
} from 'reka-ui'
import {useWorkspaceStore} from '../stores/workspace'

const workspace = useWorkspaceStore()

const createDialogOpen = ref(false)
const newPath = ref('')
const newTitle = ref('')
const newAuthor = ref('')
const newLanguage = ref('pt-BR')

const languages = [
  {value: 'pt-BR', label: 'Português (Brasil)'},
  {value: 'en', label: 'Inglês'},
  {value: 'es', label: 'Espanhol'},
  {value: 'fr', label: 'Francês'},
  {value: 'it', label: 'Italiano'},
  {value: 'de', label: 'Alemão'},
  {value: 'la', label: 'Latim'},
]

async function choosePath() {
  const picked = await workspace.pickDirectory()
  if (picked) newPath.value = picked
}

async function submitCreate() {
  if (!newPath.value || !newTitle.value) return
  const info = await workspace.createWorkspace(newPath.value, newTitle.value, newAuthor.value, newLanguage.value)
  if (info) createDialogOpen.value = false
}
</script>

<template>
  <div class="welcome">
    <div class="welcome__card">
      <h1 class="welcome__title">Luz Writer</h1>
      <p class="welcome__subtitle">Uma IDE para escritores.</p>

      <div class="welcome__actions">
        <button class="welcome__button welcome__button--primary" type="button" @click="workspace.openWorkspaceDialog()">
          Abrir Projeto Existente
        </button>
        <button class="welcome__button" type="button" @click="createDialogOpen = true">
          Criar Novo Projeto
        </button>
      </div>

      <p v-if="workspace.error" class="welcome__error">{{ workspace.error }}</p>
    </div>

    <DialogRoot v-model:open="createDialogOpen">
      <DialogPortal>
        <DialogOverlay class="dialog-overlay" />
        <DialogContent class="dialog-content">
          <DialogTitle class="dialog-title">Criar Novo Projeto</DialogTitle>

          <div class="form-field">
            <label>Pasta do projeto</label>
            <div class="form-field__row">
              <input class="form-input" type="text" readonly :value="newPath" placeholder="Nenhuma pasta selecionada" />
              <button class="welcome__button" type="button" @click="choosePath">Escolher...</button>
            </div>
          </div>

          <div class="form-field">
            <label for="new-title">Título</label>
            <input id="new-title" class="form-input" type="text" v-model="newTitle" placeholder="Meu Livro" />
          </div>

          <div class="form-field">
            <label for="new-author">Autor</label>
            <input id="new-author" class="form-input" type="text" v-model="newAuthor" placeholder="Nome do Autor" />
          </div>

          <div class="form-field">
            <label>Idioma</label>
            <SelectRoot v-model="newLanguage">
              <SelectTrigger class="select-trigger">
                <SelectValue />
                <SelectIcon>▾</SelectIcon>
              </SelectTrigger>
              <SelectPortal>
                <SelectContent class="select-content">
                  <SelectViewport>
                    <SelectItem v-for="lang in languages" :key="lang.value" :value="lang.value" class="select-item">
                      <SelectItemText>{{ lang.label }}</SelectItemText>
                    </SelectItem>
                  </SelectViewport>
                </SelectContent>
              </SelectPortal>
            </SelectRoot>
          </div>

          <div class="dialog-actions">
            <button class="welcome__button" type="button" @click="createDialogOpen = false">Cancelar</button>
            <button
              class="welcome__button welcome__button--primary"
              type="button"
              :disabled="!newPath || !newTitle || workspace.loading"
              @click="submitCreate"
            >
              Criar Projeto
            </button>
          </div>
        </DialogContent>
      </DialogPortal>
    </DialogRoot>
  </div>
</template>

<style scoped>
.welcome {
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--luz-bg-editor);
}

.welcome__card {
  text-align: center;
}

.welcome__title {
  margin: 0;
  font-size: 2rem;
}

.welcome__subtitle {
  color: var(--luz-fg-muted);
  margin-top: 4px;
}

.welcome__actions {
  margin-top: 24px;
  display: flex;
  flex-direction: column;
  gap: 8px;
  width: 260px;
}

.welcome__button {
  padding: 8px 16px;
  border-radius: 6px;
  border: 1px solid var(--luz-border);
  background: transparent;
  color: var(--luz-fg);
  cursor: pointer;
  font-size: 0.9rem;
}

.welcome__button:hover {
  background: var(--luz-bg-hover);
}

.welcome__button--primary {
  background: var(--luz-fg);
  color: var(--luz-bg-editor);
  border-color: var(--luz-fg);
}

.welcome__button:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.welcome__error {
  margin-top: 16px;
  color: #c0392b;
  font-size: 0.85rem;
  max-width: 360px;
}

.dialog-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.4);
}

.dialog-content {
  position: fixed;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  width: 360px;
  background: var(--luz-bg-editor);
  border: 1px solid var(--luz-border);
  border-radius: 8px;
  padding: 20px;
}

.dialog-title {
  margin: 0 0 16px;
  font-size: 1.1rem;
}

.dialog-actions {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
  gap: 8px;
}

.form-field {
  margin-bottom: 14px;
  display: flex;
  flex-direction: column;
  gap: 4px;
  text-align: left;
}

.form-field label {
  font-size: 0.78rem;
  color: var(--luz-fg-muted);
}

.form-field__row {
  display: flex;
  gap: 8px;
}

.form-input {
  padding: 6px 8px;
  border-radius: 6px;
  border: 1px solid var(--luz-border);
  background: transparent;
  color: var(--luz-fg);
  flex: 1;
  min-width: 0;
}

.select-trigger {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 6px 8px;
  border-radius: 6px;
  border: 1px solid var(--luz-border);
  background: transparent;
  color: var(--luz-fg);
  cursor: pointer;
}

.select-content {
  background: var(--luz-bg-editor);
  border: 1px solid var(--luz-border);
  border-radius: 6px;
  padding: 4px;
}

.select-item {
  padding: 6px 8px;
  border-radius: 4px;
  cursor: pointer;
  outline: none;
}

.select-item[data-highlighted] {
  background: var(--luz-bg-hover);
}
</style>
