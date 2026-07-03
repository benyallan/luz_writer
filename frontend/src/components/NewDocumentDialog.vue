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

const open = defineModel<boolean>('open', {default: false})
const workspace = useWorkspaceStore()

const title = ref('')
const role = ref('chapter')

// Rótulos da seção 5.5.
const roles = [
  {value: 'chapter', label: 'Capítulo'},
  {value: 'dedication', label: 'Dedicatória'},
  {value: 'epigraph', label: 'Epígrafe'},
  {value: 'acknowledgments', label: 'Agradecimentos'},
  {value: 'preface', label: 'Prefácio'},
  {value: 'aboutAuthor', label: 'Sobre o Autor'},
  {value: 'appendix', label: 'Apêndice'},
]

async function submit() {
  if (!title.value.trim()) return
  const meta = await workspace.createChapter(title.value.trim(), role.value)
  if (meta) {
    title.value = ''
    role.value = 'chapter'
    open.value = false
  }
}
</script>

<template>
  <DialogRoot v-model:open="open">
    <DialogPortal>
      <DialogOverlay class="dialog-overlay" />
      <DialogContent class="dialog-content">
        <DialogTitle class="dialog-title">Novo Documento</DialogTitle>

        <div class="form-field">
          <label for="doc-title">Título</label>
          <input id="doc-title" class="form-input" type="text" v-model="title" placeholder="Título do documento" />
        </div>

        <div class="form-field">
          <label>Tipo</label>
          <SelectRoot v-model="role">
            <SelectTrigger class="select-trigger">
              <SelectValue />
              <SelectIcon>▾</SelectIcon>
            </SelectTrigger>
            <SelectPortal>
              <SelectContent class="select-content">
                <SelectViewport>
                  <SelectItem v-for="r in roles" :key="r.value" :value="r.value" class="select-item">
                    <SelectItemText>{{ r.label }}</SelectItemText>
                  </SelectItem>
                </SelectViewport>
              </SelectContent>
            </SelectPortal>
          </SelectRoot>
        </div>

        <div class="dialog-actions">
          <button class="welcome__button" type="button" @click="open = false">Cancelar</button>
          <button class="welcome__button welcome__button--primary" type="button" :disabled="!title.trim()" @click="submit">
            Criar
          </button>
        </div>
      </DialogContent>
    </DialogPortal>
  </DialogRoot>
</template>

<style scoped>
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
  width: 320px;
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

.form-input {
  padding: 6px 8px;
  border-radius: 6px;
  border: 1px solid var(--luz-border);
  background: transparent;
  color: var(--luz-fg);
  width: 100%;
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

/* .select-content/.select-item vivem em style.css (global) — ver comentário
   lá: conteúdo teleportado pelo Reka UI não é alcançado por <style scoped>. */

.welcome__button {
  padding: 6px 14px;
  border-radius: 6px;
  border: 1px solid var(--luz-border);
  background: transparent;
  color: var(--luz-fg);
  cursor: pointer;
  font-size: 0.85rem;
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
</style>
