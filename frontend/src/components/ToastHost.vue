<script setup lang="ts">
import {useToastStore} from '../stores/toast'

const toastStore = useToastStore()
</script>

<template>
  <div class="toast-host">
    <div v-for="t in toastStore.items" :key="t.id" class="toast" :class="`toast--${t.kind}`">
      <span class="toast__message">{{ t.message }}</span>
      <button class="toast__dismiss" type="button" title="Fechar" @click="toastStore.dismiss(t.id)">×</button>
    </div>
  </div>
</template>

<style scoped>
.toast-host {
  position: fixed;
  bottom: 32px;
  right: 16px;
  display: flex;
  flex-direction: column;
  gap: 8px;
  z-index: 1000;
  max-width: 360px;
}

.toast {
  display: flex;
  align-items: flex-start;
  gap: 8px;
  padding: 10px 12px;
  border-radius: 6px;
  font-size: 0.82rem;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.25);
  color: white;
}

.toast--error {
  background: #c0392b;
}

.toast--info {
  background: #2c3e50;
}

.toast__message {
  flex: 1;
  word-break: break-word;
}

.toast__dismiss {
  border: none;
  background: transparent;
  color: inherit;
  cursor: pointer;
  font-size: 1rem;
  line-height: 1;
  padding: 0;
  opacity: 0.8;
}

.toast__dismiss:hover {
  opacity: 1;
}
</style>
