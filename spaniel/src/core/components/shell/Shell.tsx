import { ThemeProvider } from 'styled-components';

import { ThemeProvider as MuiThemeProvider } from '@mui/material/styles';

import { CssBaseline } from '@mui/material';

import { LoginPage } from 'auth/pages';

import { darkTheme } from 'config/theming';

import { UnauthenticatedLayout } from './UnauthenticatedLayout';

export const Shell: React.FC = () => {
  const authenticated = false;

  return (
    <ThemeShell>
      {!authenticated && (
        <UnauthenticatedLayout>
          <LoginPage />
        </UnauthenticatedLayout>
      )}
    </ThemeShell>
  );
};

const ThemeShell: React.FC = props => {
  return (
    <>
    {/**
    * Workaround - MaterialUI does not apply theme to styled-components ThemeProvider,
    * therefore we need to use both of them. 
    */}
    <ThemeProvider theme={darkTheme}>
      <MuiThemeProvider theme={darkTheme}>
        <CssBaseline />
        {props.children}
      </MuiThemeProvider>
    </ThemeProvider>
    </>
  );
};
