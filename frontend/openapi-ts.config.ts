import { defineConfig } from '@hey-api/openapi-ts';

export default defineConfig({
  input:
    '../backend/services/dream-wiki/openapi.yml',
  output: {
    format: 'prettier',
    path: './src/client',
  },
  plugins: [
    '@hey-api/client-axios',
    '@hey-api/schemas',
    '@hey-api/sdk',
    {
      enums: 'javascript',
      name: '@hey-api/typescript',
    },
  ],
});
