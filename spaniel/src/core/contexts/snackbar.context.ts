import {
  createContext,
  useContext
} from 'react';

export enum NotificationVariant {
  SUCCESS = 'success',
  ERROR = 'error',
}

export type NotificationOptions = {
  variant: NotificationVariant;
  message: string;
};

type NotificationFn = (options: NotificationOptions) => void;

export type SnackbarContext = {
  notify: NotificationFn;
};

export const snackbarContext = createContext<SnackbarContext>(null!);

export const useSnackbarContext = () => useContext(snackbarContext);
