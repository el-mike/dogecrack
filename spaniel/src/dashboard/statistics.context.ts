import {
  createContext,
  useContext,
} from 'react';

import {
  Statistics,
} from 'models';

export type LoadFn = () => void;

export type StatisticsContext = {
  load: LoadFn;
  statistics: Statistics;
};

export const statisticsContext = createContext<StatisticsContext>(null!);

export const useStatisticsContext = () => useContext(statisticsContext);
