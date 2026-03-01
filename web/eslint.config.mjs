import { defineConfig, globalIgnores } from 'eslint/config'
import nextVitals from 'eslint-config-next/core-web-vitals'

const eslintConfig = defineConfig([
  ...nextVitals,
  // Override default ignores of eslint-config-next.
  globalIgnores([
    // Default ignores of eslint-config-next:
    '.next/**',
    'out/**',
    'build/**',
    'next-env.d.ts',
  ]),
  {
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

export default eslintConfig
