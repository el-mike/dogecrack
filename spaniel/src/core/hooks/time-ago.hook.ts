import {
  useState,
  useEffect,
} from 'react';

import {
  millisecondsInSecond,
  secondsInMinute
} from 'date-fns';

import { getDurationFromNow } from 'core/utils';


export const useTimeAgo = (updatedAt: string) => {
  const duration = getDurationFromNow(new Date(updatedAt));
  const secondsPassed =
    ((duration.minutes || 0) * secondsInMinute)
    + (duration.seconds || 0);

  const [seconds, setSeconds] = useState(secondsPassed);

  useEffect(() => {
    if (duration?.hours === 0) {
      const interval = setInterval(
        () => setSeconds(seconds + 1),
        millisecondsInSecond,
      );
  
      return () => clearInterval(interval);
    }
  }, [duration]);

  return duration?.hours === 0
    ? seconds > 60
      ? `${Math.floor(seconds / 60)} minutes ago`
      : `${seconds} seconds ago`
    : '> hour ago';
};
