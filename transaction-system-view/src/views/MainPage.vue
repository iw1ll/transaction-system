<template>
    <div class="container">

        <!-- <section class="section">
            <h2>Последние транзакции</h2>
            <div class="form-group">
                <input v-model.number="transactionsLimit" type="number" min="1" max="100"/>
                <button @click="fetchTransactions">Загрузить</button>
            </div>
            <div v-if="isLoading" class="loading">Загрузка...</div>
            <div v-else>
                <TransactionItem 
                    v-for="tx in transactions"
                    :key="tx.hash"
                    :transaction="tx"
                />
            </div>
            <div v-if="transactionsError" class="error">{{ transactionsError }}</div>
        </section> -->

        <section class="section">
            <h2>Кошельки</h2>
            <div v-for="wallet in wallets" :key="wallet.address">
                <span class="address">{{ wallet.address }}</span> <div class="balance">{{ wallet.balance }} ETH</div> 
            </div>
        </section>

        <!-- Форма отправки средств -->
        <section class="section">
            <h2>Отправить средства</h2>
            <div class="form-group">
                <input v-model="sendData.from" placeholder="Отправитель"/>
                <input v-model="sendData.to" placeholder="Получатель"/>
                <input v-model.number="sendData.amount" type="number" placeholder="Сумма"/>
                <button @click="handleSend">Отправить</button>
            </div>
            <div v-if="sendError" class="error">{{ sendError }}</div>
        </section>
  
        <!-- Проверка баланса -->
        <section class="section">
            <h2>Баланс кошелька</h2>
            <div class="form-group">
                <input v-model="walletAddress" placeholder="Адрес кошелька"/>
                <button @click="fetchBalance">Проверить</button>
            </div>
            <div v-if="balance !== null" class="balance">
                Текущий баланс: {{ balance }} ETH
            </div>
            <div v-if="balanceError" class="error">{{ balanceError }}</div>
        </section>
  
        <!-- История транзакций -->
        <section class="section">
            <h2>Последние транзакции</h2>
            <div class="form-group">
                <input v-model.number="transactionsLimit" type="number" min="1" max="100"/>
                <button @click="fetchTransactions">Загрузить</button>
            </div>
            <div v-if="isLoading" class="loading">Загрузка...</div>
            <div v-else>
                <TransactionItem 
                    v-for="tx in transactions"
                    :key="tx.hash"
                    :transaction="tx"
                />
            </div>
            <div v-if="transactionsError" class="error">{{ transactionsError }}</div>
        </section>
    </div>
</template>
  
<script setup lang="ts">
import { ref, onMounted } from 'vue';
import TransactionItem from '@/components/TransactionItem.vue';
  
// Типы
interface Transaction {
    hash: string
    from: string
    to: string
    amount: number
    timestamp: number
}
  
interface BalanceResponse {
    address: string
    balance: number
}
  
// Состояние отправки средств
const sendData = ref({
    from: '',
    to: '',
    amount: 0,
});
const sendError = ref<string>('');
  
// Состояние баланса
const walletAddress = ref<string>('');
const balance = ref<number | null>(null);
const balanceError = ref<string>('');
  
// Состояние транзакций
const transactionsLimit = ref<number>(1);
const transactions = ref<Transaction[]>([]);
const wallets = ref<BalanceResponse[]>([]);
const isLoading = ref<boolean>(false);
const transactionsError = ref<string>('');

onMounted(() => {
    getAllWallets();
});

// Отправка транзакции
const getAllWallets = async () => {
    try {
        const response = await fetch('http://localhost:8080/api/wallets');
        wallets.value =  await response.json();
    } catch (err) {
        wallets.value = [];
        throw new Error(`HTTP error! Error: ${err}`);
    }
};

// Отправка транзакции
const handleSend = async () => {
    try {
        const response = await fetch('http://localhost:8080/api/send', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(sendData.value),
        });
  
        if (!response.ok) {
            throw new Error(`HTTP error! Status: ${response.status}`);
        }
  
        sendError.value = '';
        alert('Транзакция успешно отправлена!');
        fetchTransactions();
        getAllWallets();
    } catch (err) {
        sendError.value = err instanceof Error ? err.message : 'Неизвестная ошибка';
    }
};
  
// Получение баланса
const fetchBalance = async () => {
    try {
        const response = await fetch(
            `http://localhost:8080/api/wallet/${walletAddress.value}/balance`,
        );
  
        if (!response.ok) {
            throw new Error(`HTTP error! Status: ${response.status}`);
        }
  
        const data: BalanceResponse = await response.json();
        balance.value = data.balance;
        balanceError.value = '';
    } catch (err) {
        balance.value = null;
        balanceError.value = err instanceof Error ? err.message : 'Неизвестная ошибка';
    }
};
  
// Получение транзакций
const fetchTransactions = async () => {
    try {
        isLoading.value = true;
        const response = await fetch(
            `http://localhost:8080/api/transactions?count=${transactionsLimit.value}`,
        );
  
        if (!response.ok) {
            throw new Error(`HTTP error! Status: ${response.status}`);
        }
  
        transactions.value = await response.json();
        transactionsError.value = '';
    } catch (err) {
        transactions.value = [];
        transactionsError.value = err instanceof Error ? err.message : 'Неизвестная ошибка';
    } finally {
        isLoading.value = false;
    }
};
</script>
  
  <style scoped>
  .container {
    max-width: 800px;
    margin: 0 auto;
    padding: 2rem;
  }
  
  .section {
    margin-bottom: 2rem;
    padding: 1.5rem;
    border: 1px solid #e2e8f0;
    border-radius: 0.5rem;
  }
  
  .form-group {
    display: flex;
    gap: 1rem;
    margin-bottom: 1rem;
  }
  
  input {
    padding: 0.5rem;
    border: 1px solid #cbd5e0;
    border-radius: 0.375rem;
    flex-grow: 1;
  }
  
  button {
    padding: 0.5rem 1rem;
    background-color: #4299e1;
    color: white;
    border: none;
    border-radius: 0.375rem;
    cursor: pointer;
    transition: background-color 0.2s;
  }
  
  button:hover {
    background-color: #3182ce;
  }
  
  .error {
    color: #e53e3e;
    margin-top: 0.5rem;
  }
  
  .balance {
    color: #48bb78;
    font-weight: 500;
    margin-top: 0.5rem;
  }

  .address {
    color: #111814;
    font-weight: 500;
    margin-top: 0.5rem;
  }
  
  .loading {
    color: #718096;
    text-align: center;
    padding: 1rem;
  }
  </style>