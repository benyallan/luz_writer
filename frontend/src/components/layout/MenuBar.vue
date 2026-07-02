<script setup>
import { ref, computed } from 'vue'
import {
  DropdownMenuRoot,
  DropdownMenuTrigger,
  DropdownMenuPortal,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuSub,
  DropdownMenuSubTrigger,
  DropdownMenuSubContent,
} from 'reka-ui'
import { ChevronRight } from '@lucide/vue'
import { QuitApp, GetRecentProjects, ClearRecentProjects, PrecompileToLatex } from '@wails/go/main/App'
import { useWorkspaceStore } from '@/stores/workspace'
import { useEditorStore } from '@/stores/editor'

const workspace = useWorkspaceStore()
const editor = useEditorStore()
const recentProjects = ref([])

// ── Arquivo ───────────────────────────────────────────────────────────────────

async function onArquivoOpen(open) {
  if (open) {
    recentProjects.value = await GetRecentProjects()
  }
}

async function openRecent(path) {
  await workspace.openAt(path)
}

async function clearRecents() {
  await ClearRecentProjects()
  recentProjects.value = []
}

function parentDir(fullPath) {
  const parent = fullPath.split('/').slice(0, -1).join('/') || '/'
  return parent.replace(/^\/home\/[^/]+/, '~')
}

// ── Compilar ──────────────────────────────────────────────────────────────────

const canCompile = computed(() =>
  !!(workspace.rootPath && editor.filePath && editor.fileName?.endsWith('.luztxt'))
)

const toast = ref({ show: false, message: '', type: 'success' })
let toastTimer = null

function showToast(message, type = 'success') {
  if (toastTimer) clearTimeout(toastTimer)
  toast.value = { show: true, message, type }
  toastTimer = setTimeout(() => { toast.value.show = false }, 4000)
}

async function precompileToLatex() {
  if (!canCompile.value) return
  try {
    const outPath = await PrecompileToLatex(workspace.rootPath, editor.filePath, editor.content)
    showToast(`LaTeX gerado em: ${outPath}`, 'success')
  } catch (err) {
    showToast(`Erro: ${err}`, 'error')
  }
}
</script>

<template>
  <nav class="menubar">
    <!-- ── Arquivo ── -->
    <DropdownMenuRoot @update:open="onArquivoOpen">
      <DropdownMenuTrigger as-child>
        <button class="menubar-btn">Arquivo</button>
      </DropdownMenuTrigger>
      <DropdownMenuPortal>
        <DropdownMenuContent class="menu-content" :side-offset="0" align="start">

          <DropdownMenuItem class="menu-item" @click="workspace.newProject">
            Novo projeto...
            <span class="menu-shortcut">Ctrl+N</span>
          </DropdownMenuItem>
          <DropdownMenuItem class="menu-item" @click="workspace.openFolder">
            Abrir pasta...
            <span class="menu-shortcut">Ctrl+K O</span>
          </DropdownMenuItem>

          <!-- Recentes submenu -->
          <DropdownMenuSub>
            <DropdownMenuSubTrigger class="menu-item menu-subtrigger">
              Recentes
              <ChevronRight :size="13" class="submenu-arrow" />
            </DropdownMenuSubTrigger>
            <DropdownMenuPortal>
              <DropdownMenuSubContent class="menu-content" :side-offset="2" :align-offset="-5">
                <template v-if="recentProjects.length">
                  <DropdownMenuItem
                    v-for="project in recentProjects"
                    :key="project.path"
                    class="menu-item recent-item"
                    :title="project.path"
                    @click="openRecent(project.path)"
                  >
                    <span class="recent-name">{{ project.name }}</span>
                    <span class="recent-dir">{{ parentDir(project.path) }}</span>
                  </DropdownMenuItem>
                  <DropdownMenuSeparator class="menu-separator" />
                  <DropdownMenuItem class="menu-item menu-item-muted" @click="clearRecents">
                    Limpar recentes
                  </DropdownMenuItem>
                </template>
                <span v-else class="menu-empty">Nenhum projeto recente</span>
              </DropdownMenuSubContent>
            </DropdownMenuPortal>
          </DropdownMenuSub>

          <DropdownMenuSeparator class="menu-separator" />

          <DropdownMenuItem class="menu-item" @click="QuitApp">
            Sair
            <span class="menu-shortcut">Alt+F4</span>
          </DropdownMenuItem>

        </DropdownMenuContent>
      </DropdownMenuPortal>
    </DropdownMenuRoot>

    <!-- ── Compilar ── -->
    <DropdownMenuRoot>
      <DropdownMenuTrigger as-child>
        <button class="menubar-btn">Compilar</button>
      </DropdownMenuTrigger>
      <DropdownMenuPortal>
        <DropdownMenuContent class="menu-content" :side-offset="0" align="start">

          <DropdownMenuItem
            class="menu-item"
            :disabled="!canCompile"
            @click="precompileToLatex"
          >
            Pré-compilar para LaTeX
          </DropdownMenuItem>

        </DropdownMenuContent>
      </DropdownMenuPortal>
    </DropdownMenuRoot>

    <!-- Toast de notificação -->
    <Teleport to="body">
      <Transition name="toast">
        <div v-if="toast.show" class="compile-toast" :class="'toast-' + toast.type">
          {{ toast.message }}
        </div>
      </Transition>
    </Teleport>
  </nav>
</template>

<style scoped>
.menubar {
  background: var(--color-menubar);
  display: flex;
  align-items: stretch;
  height: var(--menubar-height);
  user-select: none;
}

.menubar-btn {
  height: 100%;
  padding: 0 10px;
  font-size: 12px;
  color: var(--color-text);
  border-radius: 0;
  transition: background 0.1s;
}

.menubar-btn:hover,
:global([data-state="open"]) > .menubar-btn {
  background: rgba(255, 255, 255, 0.1);
}
</style>

<style>
/* ── Shared menu styles (global, consumed by portal content) ── */
.menu-content {
  background: #252526;
  border: 1px solid var(--color-border);
  border-radius: 2px;
  padding: 4px 0;
  min-width: 220px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.4);
  z-index: 999;
  animation: fadeIn 0.08s ease-out;
}

.menu-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 4px 20px 4px 28px;
  height: 26px;
  cursor: pointer;
  color: var(--color-text);
  font-size: 13px;
  outline: none;
  user-select: none;
  gap: 8px;
}

.menu-item[data-highlighted] {
  background: var(--color-accent);
  color: #fff;
}

.menu-item[data-disabled] {
  color: var(--color-text-muted);
  cursor: default;
  pointer-events: none;
}

/* SubTrigger keeps highlight when submenu is open */
.menu-item[data-state="open"] {
  background: var(--color-hover);
}
.menu-item[data-state="open"][data-highlighted] {
  background: var(--color-accent);
  color: #fff;
}

.menu-shortcut {
  font-size: 11px;
  color: var(--color-text-muted);
  margin-left: auto;
}

.menu-item[data-highlighted] .menu-shortcut {
  color: rgba(255, 255, 255, 0.7);
}

.menu-subtrigger {
  justify-content: flex-start;
}

.submenu-arrow {
  margin-left: auto;
  color: var(--color-text-muted);
  flex-shrink: 0;
}

.menu-item[data-highlighted] .submenu-arrow {
  color: rgba(255, 255, 255, 0.7);
}

/* Recent project item */
.recent-item {
  flex-direction: column;
  align-items: flex-start;
  height: auto;
  padding: 5px 16px 5px 28px;
  gap: 1px;
}

.recent-name {
  font-size: 13px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  max-width: 240px;
}

.recent-dir {
  font-size: 11px;
  color: var(--color-text-muted);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  max-width: 240px;
  direction: rtl;
  text-align: left;
}

.menu-item[data-highlighted] .recent-dir {
  color: rgba(255, 255, 255, 0.65);
}

.menu-item-muted {
  color: var(--color-text-muted);
}

.menu-item-muted[data-highlighted] {
  color: #fff;
}

.menu-separator {
  height: 1px;
  background: var(--color-border);
  margin: 4px 0;
}

.menu-empty {
  display: block;
  padding: 6px 28px;
  font-size: 12px;
  color: var(--color-text-muted);
  user-select: none;
}

@keyframes fadeIn {
  from { opacity: 0; transform: translateY(-4px); }
  to   { opacity: 1; transform: translateY(0); }
}

/* ── Toast ── */
.compile-toast {
  position: fixed;
  bottom: 24px;
  right: 24px;
  max-width: 480px;
  padding: 10px 16px;
  border-radius: 4px;
  font-size: 12.5px;
  line-height: 1.5;
  z-index: 9999;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.5);
  word-break: break-all;
}

.toast-success {
  background: #1e3a2f;
  border: 1px solid #2ea05a;
  color: #7ee8a2;
}

.toast-error {
  background: #3a1e1e;
  border: 1px solid #a02e2e;
  color: #e87e7e;
}

.toast-enter-active,
.toast-leave-active {
  transition: opacity 0.2s, transform 0.2s;
}
.toast-enter-from,
.toast-leave-to {
  opacity: 0;
  transform: translateY(8px);
}
</style>
