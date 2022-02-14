/**
 * NOTE:
 * 
 * Atm, craco cannot be used, as it's not ready for CRA 5.0. For now,
 * @emotion/react and @emotion/styled have been added, to allow the app to work.
 * As soon as craco is ready for CRA 5.0, we should remove aforementioned deps and
 * use this config patch instead.
 */
const CracoAlias = require('craco-alias');

module.exports = {
  plugins: [
    {
      plugin: CracoAlias,
      options: {
        source: 'tsconfig',
        /* tsConfigPath should point to the file where "paths" are specified */
        tsConfigPath: './tsconfig.json',
      },
    },
  ],
  webpack: {
    alias: {
      '@mui/styled-engine': '@mui/styled-engine-sc',
    },
  },
};
