import {
  format,
  intervalToDuration,
} from 'date-fns';

export const DEFAULT_DATE_TIME_FORMAT = 'yyyy-MM-dd HH:mm';

export const toDateTimeString = (date: Date) => {
  if (isNaN(date.getTime())) {
    return '';
  }

  try {
    return format(date, DEFAULT_DATE_TIME_FORMAT);
  } catch (error: unknown) {
    console.error(error);
    return '';
  }
};

export const getDurationFromNow = (date: Date) => {
  if (isNaN(date.getTime())) {
    return {};
  }

  return intervalToDuration({
    start: date,
    end: new Date(),
  });
};

export const getDuration = (start: Date, end: Date) => {
  if (isNaN(start.getTime()) || isNaN(end.getTime())) {
    return {};
  }

  return intervalToDuration({ start, end });
}
