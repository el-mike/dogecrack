import {
  minutesInHour,
  secondsInHour,
} from 'date-fns';

import { getDuration } from './date-utils';

import { PitbullInstance } from 'models';

export const getInstanceEstimatedCost = (instance: PitbullInstance) => {
  if (!instance.startedAt || !instance.hostInstance.dphTotal) {
    return 0;
  }

  const start = new Date(instance.startedAt);
  const end = (instance.completedAt && new Date(instance.completedAt)) || new Date();

  const duration = getDuration(start, end);

  const hoursTotal =
    (duration?.days || 0) * 24
    + (duration?.hours || 0)
    + ((duration?.minutes || 0) / minutesInHour)
    + ((duration?.seconds || 0) / secondsInHour);

  return +(hoursTotal * instance.hostInstance.dphTotal).toFixed(3);
};
