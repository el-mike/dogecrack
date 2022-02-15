import {
  useState,
  useEffect,
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

  const loadUser = (callback?: LoginFnCallback) => authService.getMe()
    .then(user => {
      setUser(user);
      callback?.(user);
    })
    .catch(() => setUser(null))
    .finally(() => setUserLoading(false));

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

  const clear = () => setUser(null);

  const logout = (callback?: LogoutFnCallback) => {
    authService.logout()
      .then(() => {
        clear();
        callback?.()
      });
  };

  const value = {
    user,
    userLoading,
    loginLoading,
    login,
    logout,
    clear,
  } as AuthContext;
  
  return (
    <authContext.Provider value={value}>{props.children}</authContext.Provider>
  );
};
