import {
  authService,
  crackJobService,
  generalService,
} from '../services';

export const useAuthService = () => authService;
export const useCrackJobService = () => crackJobService;
export const useGeneralService = () => generalService;
