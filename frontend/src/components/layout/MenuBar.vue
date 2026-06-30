<script setup>
import {
  DropdownMenuRoot,
  DropdownMenuTrigger,
  DropdownMenuPortal,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
} from 'reka-ui'
import { QuitApp } from '@wails/go/main/App'
import { useWorkspaceStore } from '@/stores/workspace'

const workspace = useWorkspaceStore()
</script>

<template>
  <nav class="menubar">
    <DropdownMenuRoot>
      <DropdownMenuTrigger class="menubar-item" as-child>
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
:global([data-state="open"]) .menubar-btn {
  background: rgba(255, 255, 255, 0.1);
}
</style>

<style>
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
}

.menu-item[data-highlighted] {
  background: var(--color-accent);
  color: #fff;
}

.menu-shortcut {
  font-size: 11px;
  color: var(--color-text-muted);
  margin-left: 20px;
}

.menu-item[data-highlighted] .menu-shortcut {
  color: rgba(255, 255, 255, 0.7);
}

.menu-separator {
  height: 1px;
  background: var(--color-border);
  margin: 4px 0;
}

@keyframes fadeIn {
  from { opacity: 0; transform: translateY(-4px); }
  to   { opacity: 1; transform: translateY(0); }
}
</style>
