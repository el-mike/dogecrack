import {
  useState,
  useEffect,
} from 'react';

import {
  Statistics,
} from 'models';

import {
  useGeneralService,
} from 'core/hooks';

import {
  statisticsContext,
  StatisticsContext,
} from '../statistics.context';

export const StatisticsProvider: React.FC = props => {
  const generalService = useGeneralService();

  const [statistics, setStatistics] =
    useState<Statistics>({} as Statistics);

  const load = () => {
    generalService.getStatistics()
      .then(statistics => setStatistics(statistics));
  };

  const value = {
    load,
    statistics
  } as StatisticsContext;

  useEffect(() => {
    load();
  }, []);

  return (
    <statisticsContext.Provider value={value}>{props.children}</statisticsContext.Provider>
  );
};
