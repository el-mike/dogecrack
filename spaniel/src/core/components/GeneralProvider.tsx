import {
  useState,
  useEffect,
  useMemo,
} from 'react';

import { Enums } from 'models';

import { useGeneralService } from 'core/hooks';

import {
  generalContext,
  GeneralContext,
} from '../contexts';

export const GeneralProvider: React.FC = props => {
  const generalService = useGeneralService();

  const [enums, setEnums] = useState<Enums | null>(null);
  const [enumsLoading, setEnumsLoading] = useState(false);

  const latestTokenGeneratorVersion = useMemo(() =>
      enums ? Math.max(...Object.values(enums.tokenGeneratorVersion)) : 0,
    [enums]
  );

  useEffect(() => {
    setEnumsLoading(true);

    generalService.getEnums()
      .then(enums => setEnums(enums))
      .finally(() => setEnumsLoading(false));
  }, []);

  const value = {
    enums,
    enumsLoading,
    latestTokenGeneratorVersion,
  } as GeneralContext;

  return (
    <generalContext.Provider value={value}>{props.children}</generalContext.Provider>
  );
};
