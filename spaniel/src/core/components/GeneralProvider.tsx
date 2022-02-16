import {
  useState,
  useEffect,
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

  useEffect(() => {
    setEnumsLoading(true);
  
    generalService.getEnums()
      .then(enums => setEnums(enums))
      .finally(() => setEnumsLoading(false));
  }, []);

  const value = {
    enums,
    enumsLoading,
  } as GeneralContext;

  return (
    <generalContext.Provider value={value}>{props.children}</generalContext.Provider>
  );
};
