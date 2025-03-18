// stores/counter.js
import { defineStore } from 'pinia';

export const useCounterStore = defineStore('counter', {
    state: () => ({
        count: 0,
        name: 'Vue Master',
    }),
    getters: {
        doubleCount: (state) => state.count * 2,
        // Геттер с аргументом
        greet: (state) => (greeting: any) => `${greeting}, ${state.name}!`,
    },
    actions: {
        increment() {
            this.count++;
        },
        async fetchData() {
            // Асинхронный пример
            const res = await fetch('https://jsonplaceholder.typicode.com/todos/1');
            const data = await res.json();
            console.log(data);
        },
    },
});