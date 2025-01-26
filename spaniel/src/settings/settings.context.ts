import {
  createContext,
  useContext,
} from 'react';

import { Settings } from 'models';

export type UpdateFn = (settings: Settings) => void;

export type SettingsContext = {
  settings: Settings;
  loading: boolean;
  update: UpdateFn;
};

export const settingsContext = createContext<SettingsContext>(null!);

export const useSettingsContext = () => useContext(settingsContext);
