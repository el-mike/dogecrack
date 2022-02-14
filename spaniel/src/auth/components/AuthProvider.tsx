import {
  useState,
  useEffect,
  useCallback,
} from 'react';

import {
  User,
  UserCredentials,
} from 'models';

import { useAuthService } from 'core/hooks';

import {
  LoginFnCallback,
  LogoutFnCallback,
  AuthContext,
  authContext,
} from '../auth.context';

export const AuthProvider: React.FC = props => {
  const authService = useAuthService();

  const [user, setUser] = useState<User | null>(null);
  const [userLoading, setUserLoading] = useState<boolean>(true);
  const [loginLoading, setLoginLoading] = useState<boolean>(false);

  const loadUser = useCallback(
    (callback?: LoginFnCallback) => authService.getMe()
      .then(user => {
        setUser(user);
        callback?.(user);
      })
      .catch(() => setUser(null))
      .finally(() => setUserLoading(false)),
    [authService, setUser, setUserLoading],
  );

  useEffect(() => {
    loadUser();
    /* eslint-disable-next-line */
  }, []);

  const login = (creds: UserCredentials, callback?: LoginFnCallback) => {
    setLoginLoading(true);
  
    authService.login(creds)
      .then(() => loadUser(callback))
      .finally(() => setLoginLoading(false));
  };

  const logout = (callback?: LogoutFnCallback) => {
    authService.logout()
      .then(() => {
        setUser(null);
        callback?.()
      });
  };

  const value = {
    user,
    userLoading,
    loginLoading,
    login,
    logout,
  } as AuthContext;
  
  return (
    <authContext.Provider value={value}>{props.children}</authContext.Provider>
  );
};
