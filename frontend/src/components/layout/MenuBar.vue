<script setup>
import { ref } from 'vue'
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
import { QuitApp, GetRecentProjects, ClearRecentProjects } from '@wails/go/main/App'
import { useWorkspaceStore } from '@/stores/workspace'

const workspace = useWorkspaceStore()
const recentProjects = ref([])

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
</script>

<template>
  <nav class="menubar">
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
</style>
