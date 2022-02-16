import {
  useState,
  useEffect,
} from 'react';

import { millisecondsInSecond } from 'date-fns';

import { timeForPipe } from '../pipes';

export type TimeForProps = {
  from: string;
  to?: string;
};

export const TimeFor: React.FC<TimeForProps> = props => {
  const { from, to } = props;

  const [clearIntervalCb, setClearIntervalCb] = useState<any>(null);
  const [timeFor, setTimeFor] = useState(timeForPipe(from, to || new Date().toISOString()));

  useEffect(() => {
      const interval = setInterval(
        () => setTimeFor(timeForPipe(from, to || new Date().toISOString())),
        millisecondsInSecond,
      );

      setClearIntervalCb(interval);
  
      return () => clearInterval(clearIntervalCb);
  }, [from]);

  useEffect(() => {
    if (to && clearIntervalCb) {
      clearInterval(clearIntervalCb);
    }
  }, [to, clearIntervalCb]);

  return (
    <>{timeFor}</>
  );
};
