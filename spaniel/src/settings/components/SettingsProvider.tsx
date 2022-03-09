import {
  useState,
  useEffect,
} from 'react';

import { Settings } from 'models';

import { useGeneralService } from 'core/hooks';

import {
  settingsContext,
  SettingsContext,
} from '../settings.context';

export const SettingsProvider: React.FC = props => {
  const [settings, setSettings] = useState<Settings>({} as Settings);
  const [loading, setLoading] = useState(false);

  const generalService = useGeneralService();

  const update = (payload: Settings) => {
    setLoading(true);

    generalService.updateSettings(payload)
      .then(updated => setSettings(updated))
      .finally(() => setLoading(false));
  };

  const value = {
    settings,
    loading,
    update,
  } as SettingsContext;

  useEffect(() => {
    setLoading(true);

    generalService.getSettings()
      .then(settings => setSettings(settings))
      .finally(() => setLoading(false));
  }, []);


  return (
    <settingsContext.Provider value={value}>{props.children}</settingsContext.Provider>
  );
};
