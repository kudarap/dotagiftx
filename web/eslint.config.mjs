import { defineConfig } from 'eslint/config'
import react from 'eslint-plugin-react'
import prettier from 'eslint-plugin-prettier'
import globals from 'globals'
import path from 'node:path'
import { fileURLToPath } from 'node:url'
import js from '@eslint/js'
import { FlatCompat } from '@eslint/eslintrc'

const __filename = fileURLToPath(import.meta.url)
const __dirname = path.dirname(__filename)
const compat = new FlatCompat({
  baseDirectory: __dirname,
  recommendedConfig: js.configs.recommended,
  allConfig: js.configs.all,
})

export default defineConfig([
  {
    extends: [
      ...compat.extends('eslint:recommended'),
      ...compat.extends('airbnb'),
      ...compat.extends('plugin:react/recommended'),
      ...compat.extends('prettier'),
    ],

    plugins: {
      react,
      prettier,
    },

    languageOptions: {
      globals: {
        ...globals.browser,
      },

      ecmaVersion: 2020,
      sourceType: 'script',
    },

    rules: {
      eqeqeq: 'off',
      'no-plusplus': 'off',
      'no-continue': 'off',
      'react/jsx-props-no-spreading': 'off',
      'react/jsx-closing-bracket-location': 'off',
      'react/forbid-prop-types': 'off',
      'react/jsx-filename-extension': 'off',
      'react/react-in-jsx-scope': 'off',

      'import/no-unresolved': [
        'error',
        {
          ignore: ['^@'],
        },
      ],

      'import/extensions': [
        'warn',
        {
          ignore: ['^@'],
        },
      ],

      'react/prop-types': 'warn',
      'react/require-default-props': 'warn',
    },
  },
])
