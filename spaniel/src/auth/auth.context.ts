import { createContext } from 'react';

import {
  User,
  UserCredentials,
} from 'models';

export type LoginFnCallback = (user: User) => void;
export type LogoutFnCallback = () => void;

export type LoginFn = (creds: UserCredentials, callback?: LoginFnCallback) => void;
export type LogoutFn = (callback?: LogoutFnCallback) => void;

export type AuthContext = {
  user: User;
  userLoading: boolean;
  loginLoading: boolean;
  login: LoginFn;
  logout: LogoutFn;
}

export const authContext = createContext<AuthContext>(null!);
