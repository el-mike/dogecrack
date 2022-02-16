import {
  authService,
  pitbullJobService,
  generalService,
} from '../services';

export const useAuthService = () => authService;
export const usePitbullJobService = () => pitbullJobService;
export const useGeneralService = () => generalService;
