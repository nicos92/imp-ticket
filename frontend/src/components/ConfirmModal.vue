<script lang="ts" setup>
defineProps<{
  visible: boolean
  cantidad: number
}>()

const emit = defineEmits<{
  confirmar: []
  cancelar: []
}>()
</script>

<template>
  <Teleport to="body">
    <div v-if="visible" class="modal-overlay" @click.self="emit('cancelar')">
      <div class="modal">
        <h2>Confirmar impresión</h2>
        <p>
          ¿Imprimir <strong>{{ cantidad }}</strong> etiquetas
          ({{ cantidad / 2 }} pares)?
        </p>
        <div class="modal-acciones">
          <button class="btn" @click="emit('cancelar')">Cancelar</button>
          <button class="btn btn-confirm" @click="emit('confirmar')">Confirmar</button>
        </div>
      </div>
    </div>
  </Teleport>
</template>

<style scoped>
.modal-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.6);
  display: flex;
  align-items: center;
  justify-content: center;
}

.modal {
  background: #1b2636;
  border: 1px solid #2d3a4d;
  border-radius: 10px;
  padding: 2rem;
  max-width: 420px;
  width: 90%;
}

.modal-acciones {
  display: flex;
  gap: 1rem;
  justify-content: flex-end;
  margin-top: 1.5rem;
}

.btn {
  background: #2d3a4d;
  color: white;
  border: none;
  border-radius: 6px;
  padding: 0.9rem 1.5rem;
  font-size: 1rem;
  cursor: pointer;
  transition: background 0.2s;
}

.btn:hover {
  background: #3b4c63;
}

.btn-confirm {
  background: #2f855a;
}

.btn-confirm:hover {
  background: #38a169;
}
</style>
