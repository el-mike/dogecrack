import {
  authService,
  crackJobService,
  pitbullInstanceService,
  generalService,
} from '../services';

export const useAuthService = () => authService;
export const useCrackJobService = () => crackJobService;
export const usePitbullInstanceService = () => pitbullInstanceService;
export const useGeneralService = () => generalService;
