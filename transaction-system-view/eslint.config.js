import globals from 'globals';
import pluginJs from '@eslint/js';
import tseslint from 'typescript-eslint';
import pluginVue from 'eslint-plugin-vue';


/** @type {import('eslint').Linter.Config[]} */
export default [
    { files: ['**/*.{js,mjs,cjs,ts,vue}'] },
    { languageOptions: { globals: globals.browser } },
    pluginJs.configs.recommended,
    ...tseslint.configs.recommended,
    ...pluginVue.configs['flat/essential'],
    { files: ['**/*.vue'], languageOptions: { parserOptions: { parser: tseslint.parser } } },
    {
    // Общие правила для всех файлов
        rules: {
            // Базовые правила JavaScript/TypeScript
            'semi': ['error', 'always'],
            'indent': ['error', 4],
            'quotes': ['error', 'single'],
            'comma-dangle': ['error', 'always-multiline'],
            'object-curly-spacing': ['error', 'always'],
      
            // Правила для TypeScript
            '@typescript-eslint/explicit-function-return-type': 'off',
            '@typescript-eslint/no-explicit-any': 'warn',
      
            // Правила для Vue
            'vue/html-indent': ['error', 4],
            'vue/script-indent': ['error', 4],
            'vue/multi-word-component-names': 'off',
            'vue/attribute-hyphenation': ['error', 'always'],
            'vue/html-self-closing': ['error', {
                'html': {
                    'void': 'always',
                    'normal': 'never',
                    'component': 'always',
                },
            }],
            'vue/order-in-components': ['error', {
                'order': [
                    'el',
                    'name',
                    'key',
                    'parent',
                    'functional',
                    ['delimiters', 'comments'],
                    ['components', 'directives', 'filters'],
                    'extends',
                    'mixins',
                    ['provide', 'inject'],
                    'ROUTER_GUARDS',
                    'layout',
                    'middleware',
                    'validate',
                    'scrollToTop',
                    'transition',
                    'loading',
                    'inheritAttrs',
                    'model',
                    ['props', 'propsData'],
                    'emits',
                    'setup',
                    'asyncData',
                    'data',
                    'computed',
                    'watch',
                    'watchQuery',
                    'LIFECYCLE_HOOKS',
                    'methods',
                    ['template', 'render'],
                    'renderError',
                ],
            }],
        },
    },
];