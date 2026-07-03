<script setup lang="ts">
import {computed} from 'vue'
import {nodeViewProps, NodeViewWrapper, NodeViewContent} from '@tiptap/vue-3'
import {PopoverRoot, PopoverTrigger, PopoverPortal, PopoverContent} from 'reka-ui'

const props = defineProps(nodeViewProps)

const labels: Record<string, string> = {
  luzChapter: 'Capítulo',
  luzSection: 'Seção',
  luzSubsection: 'Subseção',
}
const label = computed(() => labels[props.node.type.name] ?? props.node.type.name)
const numbered = computed(() => props.node.attrs.numbered !== false)
const includeInToc = computed(() => props.node.attrs.includeInToc !== false)

function setAttr(key: string, event: Event) {
  props.updateAttributes({[key]: (event.target as HTMLInputElement).checked})
}
</script>

<template>
  <NodeViewWrapper :class="['luz-heading', `luz-heading--${node.type.name}`]" as="div">
    <PopoverRoot>
      <PopoverTrigger class="luz-heading__label" contenteditable="false" type="button">
        {{ label }}
        <span v-if="!numbered" class="luz-heading__badge" title="Sem numeração">*</span>
        <span v-if="!includeInToc" class="luz-heading__badge" title="Fora do sumário">⊘</span>
      </PopoverTrigger>
      <PopoverPortal>
        <PopoverContent class="luz-heading__popover" :side-offset="4">
          <label class="luz-heading__option">
            <input type="checkbox" :checked="numbered" @change="setAttr('numbered', $event)" />
            Numerado
          </label>
          <label class="luz-heading__option">
            <input type="checkbox" :checked="includeInToc" @change="setAttr('includeInToc', $event)" />
            Incluir no sumário
          </label>
        </PopoverContent>
      </PopoverPortal>
    </PopoverRoot>
    <NodeViewContent class="luz-heading__content" as="span" />
  </NodeViewWrapper>
</template>

<style scoped>
.luz-heading {
  position: relative;
  margin: 1.4em 0 0.6em;
  display: flex;
  align-items: baseline;
  gap: 8px;
}

.luz-heading--luzChapter .luz-heading__content {
  font-size: 1.6rem;
  font-weight: 700;
}

.luz-heading--luzSection .luz-heading__content {
  font-size: 1.3rem;
  font-weight: 700;
}

.luz-heading--luzSubsection .luz-heading__content {
  font-size: 1.1rem;
  font-weight: 600;
}

.luz-heading__label {
  flex: 0 0 auto;
  border: none;
  background: transparent;
  color: var(--luz-fg-muted);
  font-size: 0.68rem;
  text-transform: uppercase;
  letter-spacing: 0.04em;
  cursor: pointer;
  opacity: 0.35;
  padding: 2px 4px;
  border-radius: 4px;
}

.luz-heading:hover .luz-heading__label,
.luz-heading:focus-within .luz-heading__label {
  opacity: 1;
  background: var(--luz-bg-hover);
}

.luz-heading__badge {
  margin-left: 2px;
}

.luz-heading__content {
  outline: none;
}

/* :global() porque o PopoverContent do Reka UI é teleportado — <style
   scoped> não alcança conteúdo dentro do Teleport (o atributo de escopo do
   Vue cai no wrapper de posicionamento, não nestes elementos). */
:global(.luz-heading__popover) {
  background: var(--luz-bg-editor);
  border: 1px solid var(--luz-border);
  border-radius: 6px;
  padding: 8px 10px;
  display: flex;
  flex-direction: column;
  gap: 6px;
  font-size: 0.8rem;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

:global(.luz-heading__option) {
  display: flex;
  align-items: center;
  gap: 6px;
  cursor: pointer;
}
</style>
