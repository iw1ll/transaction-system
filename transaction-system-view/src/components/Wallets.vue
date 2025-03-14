<template>
    <section class="section">
        <!-- Все кошельки -->
        <h2>Кошельки</h2>
        <div v-if="isLoadingWallet" class="loading">Загрузка кошельков...</div>
        <div v-else>
            <div v-for="wallet in wallets" :key="wallet.address">
                <span class="address">{{ wallet.address }}</span> <div class="balance">{{ wallet.balance }} ETH</div> 
            </div>
        </div>

        <div class="btn-action-wrapper">
            <button @click="getAllWallets">Загрузить все кошельки!</button>
            <button @click="hideWallets">Скрыть кошельки!</button>
            <button @click="hi">Hi!</button>
            <button @click="hiThen">hiThen!</button>
        </div>
    </section>
</template>


<script setup lang="ts">
import { ref, onMounted } from 'vue';
import type { BalanceResponse } from '@/types';


// Кошельки
const wallets = ref<BalanceResponse[]>([]);
const isLoadingWallet = ref<boolean>(false);

onMounted(() => {
    hi()
        .then((data: number[]) => {
            data.forEach((item: number) => {
                console.log(item);
            });
            return data;
        })
        .then((items: number[]) => {
            console.log(items.map(item => item * 2));
        });
});

// Получить все кошельки
const getAllWallets = async () => {
    try {
        isLoadingWallet.value = true;
        const response = await fetch('http://localhost:8080/api/wallets');
        wallets.value =  await response.json();
    } catch (err) {
        wallets.value = [];
        throw new Error(`HTTP error! Error: ${err}`);
    } finally {
        isLoadingWallet.value = false;
    }
};

// Стереть (временно)
const hideWallets = () => wallets.value = [];

const hiThen = () => {
    hi().then((data: number[]) => {
        data.forEach((item: number) => {
            console.log(item);
        });
        return data;
    })
        .then((items: number[]) => {
            console.log(items.map(item => item * 2));
        });
};

//Check
const hi = ():Promise<number[]> => {
    console.log('hi!');
    return new Promise<number[]>((resolve) => {
        const data: number[] = [1, 2, 3];
        resolve(data);
    });
};
</script>