<script setup lang="ts">
import {ref, watch} from 'vue'
import {nodeViewProps, NodeViewWrapper} from '@tiptap/vue-3'
import {PopoverRoot, PopoverTrigger, PopoverPortal, PopoverContent} from 'reka-ui'

const props = defineProps(nodeViewProps)

const text = ref<string>(props.node.attrs.text ?? '')

watch(
  () => props.node.attrs.text,
  value => {
    text.value = value ?? ''
  },
)

function commit() {
  props.updateAttributes({text: text.value})
}
</script>

<template>
  <NodeViewWrapper as="span" class="luz-footnote">
    <PopoverRoot>
      <PopoverTrigger class="luz-footnote__marker" contenteditable="false" type="button">
        [{{ node.attrs.number }}]
      </PopoverTrigger>
      <PopoverPortal>
        <PopoverContent class="luz-footnote__popover" :side-offset="4">
          <textarea
            v-model="text"
            rows="3"
            placeholder="Texto da nota de rodapé"
            @blur="commit"
          ></textarea>
        </PopoverContent>
      </PopoverPortal>
    </PopoverRoot>
  </NodeViewWrapper>
</template>

<style scoped>
.luz-footnote__marker {
  border: none;
  background: transparent;
  color: inherit;
  vertical-align: super;
  font-size: 0.72em;
  cursor: pointer;
  padding: 0 1px;
  color: var(--luz-fg-muted);
}

.luz-footnote__marker:hover {
  color: var(--luz-fg);
}

/* :global() porque o PopoverContent do Reka UI é teleportado — <style
   scoped> não alcança conteúdo dentro do Teleport (o atributo de escopo do
   Vue cai no wrapper de posicionamento, não neste elemento). */
:global(.luz-footnote__popover) {
  background: var(--luz-bg-editor);
  border: 1px solid var(--luz-border);
  border-radius: 6px;
  padding: 8px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

:global(.luz-footnote__popover textarea) {
  width: 240px;
  padding: 6px;
  border-radius: 4px;
  border: 1px solid var(--luz-border);
  background: transparent;
  color: var(--luz-fg);
  font: inherit;
  resize: vertical;
}
</style>
