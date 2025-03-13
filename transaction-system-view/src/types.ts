// src/types.ts
export interface Transaction {
    hash: string
    from: string
    to: string
    amount: number
    timestamp: number
  }
  
export interface BalanceResponse {
    address: string
    balance: number
  }