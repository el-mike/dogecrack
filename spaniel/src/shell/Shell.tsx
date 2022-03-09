import { useEffect } from 'react';

import { ThemeProvider } from 'styled-components';

import {
  BrowserRouter,
  Routes,
  Route,
  Navigate,
} from 'react-router-dom';

import { ThemeProvider as MuiThemeProvider } from '@mui/material/styles';

import { CssBaseline } from '@mui/material';

import { darkTheme } from 'config/theming';

import { shepherdApiService } from 'core/services';

import { useGeneralContext } from 'core/contexts';
import { GeneralProvider } from 'core/components';

import {
  AuthProvider,
  ProtectedRoute,
  PublicRoute,
  useAuthContext,
  LoginPage,
} from 'auth';

import { DashboardPage } from 'dashboard/pages';

import { SettingsPage } from 'settings/pages';

import {
  ShellLoader,
  AuthenticatedLayout,
  UnauthenticatedLayout,
} from './components';

export const InnerShell: React.FC = () => {
  const {
    userLoading,
    clear,
  } = useAuthContext();

  const {
    enumsLoading,
  } = useGeneralContext();

  useEffect(() => {
    /**
     * On 401 Unauthorized, we want to simply clear the user data.
     */
    shepherdApiService.setInterceptors(() => clear());

    /* eslint-disable-next-line */
  }, []);

  const initializing = userLoading
    || enumsLoading;

  if (initializing) {
    return (
      <ShellLoader />
    );
  }

  return (
    <Routes>
      {/* Default path should redirect to /dashboard. */}
      <Route path='/' element={<Navigate to='/dashboard' />} />

      <Route element={<UnauthenticatedLayout />}>
        <Route
          path='/login'
          element={
            <PublicRoute>
              <LoginPage />
            </PublicRoute>
          }
        />
      </Route>

      <Route element={<AuthenticatedLayout />}>
        <Route
          path='/dashboard'
          element={
            <ProtectedRoute>
              <DashboardPage />
            </ProtectedRoute>
          }
        />

        <Route
          path='/passchecks'
          element={
            <ProtectedRoute>
              <DashboardPage />
            </ProtectedRoute>
          }
        />

        <Route
          path='/settings'
          element={
            <ProtectedRoute>
              <SettingsPage />
            </ProtectedRoute>
          }
        />
      </Route>
    </Routes>
  );
};

export const Shell: React.FC = () => {
  return (
    <ThemeShell>
      <BrowserRouter>
        <AuthProvider>
          <GeneralProvider>
            <InnerShell />
          </GeneralProvider>
        </AuthProvider>
      </BrowserRouter>
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
