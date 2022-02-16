import {
  useState,
  useEffect,
} from 'react';

import {
  millisecondsInSecond,
  secondsInMinute
} from 'date-fns';

import {
  getDurationFromNow,
  toDateTimeString,
} from 'core/utils';

export type TimeAgoProps = {
  from: string;
};

export const TimeAgo: React.FC<TimeAgoProps> = props => {
  const { from } = props;

  const timestamp = new Date(from);

  const duration = getDurationFromNow(timestamp);
  const secondsPassed = ((duration.minutes || 0) * secondsInMinute) + (duration.seconds || 0);

  const [seconds, setSeconds] = useState(secondsPassed);

  useEffect(() => {
    if (duration?.hours === 0) {
      const interval = setInterval(
        () => setSeconds(seconds + 1),
        millisecondsInSecond,
      );
  
      return () => clearInterval(interval);
    }
  }, [timestamp]);

  const timeAgo = duration?.hours === 0
    ? seconds > 60
      ? `${Math.floor(seconds / 60)} minutes ago`
      : `${seconds} seconds ago`
    : `${toDateTimeString(timestamp)}`;

  return (
    <>{timeAgo}</>
  );
};
