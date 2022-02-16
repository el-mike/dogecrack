import {
  useContext,
  createContext,
} from 'react';

export type ToggleNavigationFn = () => void;

export type NavigationContext = {
  isOpen: boolean;
  open: ToggleNavigationFn;
  close: ToggleNavigationFn;
};

export const navigationContext = createContext<NavigationContext>(null!);

export const useNavigationContext = () => useContext(navigationContext);
