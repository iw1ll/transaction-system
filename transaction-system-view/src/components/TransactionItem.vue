<!-- src/components/TransactionItem.vue -->
<template>
    <div class="transaction">
      <div class="transaction-header">
        <span class="hash">{{ transaction.hash }}</span>
        <span class="amount">{{ transaction.amount }} ETH</span>
      </div>
      <div class="transaction-details">
        <div class="address">
          <span class="label">От:</span>
          {{ transaction.from }}
        </div>
        <div class="address">
          <span class="label">Кому:</span>
          {{ transaction.to }}
        </div>
        <div class="timestamp">
          {{ new Date(transaction.timestamp).toLocaleString() }}
        </div>
      </div>
    </div>
  </template>
  
  <script setup lang="ts">
  import type { Transaction } from '@/types'
  
  defineProps<{
    transaction: Transaction
  }>()
  </script>
  
  <style scoped>
  .transaction {
    padding: 1rem;
    margin-bottom: 1rem;
    border: 1px solid #e2e8f0;
    border-radius: 0.5rem;
  }
  
  .transaction-header {
    display: flex;
    justify-content: space-between;
    margin-bottom: 0.5rem;
  }
  
  .hash {
    font-family: monospace;
    color: #4a5568;
    font-size: 0.875rem;
  }
  
  .amount {
    font-weight: 600;
    color: #48bb78;
  }
  
  .transaction-details {
    display: grid;
    gap: 0.5rem;
  }
  
  .address {
    display: flex;
    gap: 0.5rem;
    color: #4a5568;
  }
  
  .label {
    font-weight: 500;
    color: #2d3748;
  }
  
  .timestamp {
    font-size: 0.875rem;
    color: #718096;
  }
  </style>