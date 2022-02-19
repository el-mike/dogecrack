import { createTheme } from '@mui/material';

export const darkTheme = createTheme({
  typography: {
    subtitle2: {
      fontWeight: 600,
    }
  },
  palette: {
    primary: {
      main: '#1976d2',
      light: '#63a4ff',
      dark: '#004ba0'
    },
    mode: 'dark',
  }
});
