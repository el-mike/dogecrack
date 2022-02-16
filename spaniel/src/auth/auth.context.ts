import {
  createContext,
  useContext,
} from 'react';

import {
  User,
  UserCredentials,
} from 'models';

export type LoginFnCallback = (user: User) => void;
export type LogoutFnCallback = () => void;

export type LoginFn = (creds: UserCredentials, callback?: LoginFnCallback) => void;
export type LogoutFn = (callback?: LogoutFnCallback) => void;
export type ClearAuthFn = () => void;

export type AuthContext = {
  user: User;
  userLoading: boolean;
  loginLoading: boolean;
  login: LoginFn;
  logout: LogoutFn;
  clear: ClearAuthFn;
};

export const authContext = createContext<AuthContext>(null!);

export const useAuthContext = () => useContext(authContext);
