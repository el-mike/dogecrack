import {
  authService,
  pitbullJobService,
} from '../services';

export const useAuthService = () => authService;
export const usePitbullJobService = () => pitbullJobService;
