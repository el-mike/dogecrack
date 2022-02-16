import { getDuration } from 'core/utils';

export const timeForPipe = (start: string, end: string) => {
  const duration = getDuration(new Date(start), new Date(end));

  const days = duration.days ? `${duration.days}d, ` : '';
  const hours = (duration.hours || duration.days) ? `${duration.hours}h, ` : '';
  const minutes = (duration.minutes || duration.hours || duration.days) ? `${duration.minutes}min, ` : '';
  const seconds = `${duration.seconds || 0}s`;

  return `${days}${hours}${minutes}${seconds}`;
};
