import { useState } from 'react';

import {
  navigationContext,
  NavigationContext,
} from '../navigation.context';

export const NavigationProvider: React.FC = props => {
  const [isOpen, setIsOpen] = useState(false);

  const open = () => setIsOpen(true);
  const close = () => setIsOpen(false);

  const value = {
    isOpen,
    open,
    close,
  } as NavigationContext;

  return (
    <navigationContext.Provider value={value}>{props.children}</navigationContext.Provider>
  );
};
