import {
  useContext,
  createContext,
} from 'react';

import { Enums } from 'models';

/**
 * GeneralContext - contains app-wide, domain specific data.
 */
export type GeneralContext = {
  enumsLoading: boolean;
  enums: Enums;
}

export const generalContext = createContext<GeneralContext>(null!);

export const useGeneralContext = () => useContext(generalContext);
