import {
  JobStatusKey,
  PitbullInstanceStatusKey,
  Enum,
  Dictionary,
} from 'models';

import { InputOption } from 'common/components';

export const enumLabelMap = {
  /**
   * JobStatus enums.
   */
  [JobStatusKey.JOB_SCHEDULED]: 'Scheduled',
  [JobStatusKey.JOB_PROCESSING]: 'Processing',
  [JobStatusKey.JOB_RESCHEDULED]: 'Rescheduled',
  [JobStatusKey.JOB_REJECTED]: 'Rejected',
  [JobStatusKey.JOB_ACKNOWLEDGED]: 'Acknowledged',
  /**
   * PitbullInstanceStatus enums.
   */
  [PitbullInstanceStatusKey.WAITING_FOR_HOST]: 'Waiting for host',
  [PitbullInstanceStatusKey.HOST_STARTING]: 'Host starting',
  [PitbullInstanceStatusKey.WAITING]: 'Waiting',
  [PitbullInstanceStatusKey.RUNNING]: 'Running',
  [PitbullInstanceStatusKey.FINISHED]: 'Finished',
  [PitbullInstanceStatusKey.SUCCESS]: 'Success',
  [PitbullInstanceStatusKey.INTERRUPTED]: 'Interrupted',
} as Dictionary<string>;

export const getLabelForEnum = (
  e: Enum,
  enumValue: number,
) => {
  const enumKey = Object.keys(e).find(key => e[key] === enumValue) || '';

  return (enumKey && enumLabelMap[enumKey]) || '';
};

export const getEnumAsInputOptions = (
  e: Enum,
) => Object
  .values(e)
  .map(enumValue => ({
    label: getLabelForEnum(e, enumValue),
    value: enumValue,
  } as InputOption<number>))

  