import {
  JobStatusKey,
  PitbullInstanceStatusKey,
  PitbullStatusKey,
  TokenGeneratorVersionKey,
  Enum,
  Dictionary,
} from 'models';

import { InputOption } from 'common/components';

export const enumLabelMap = {
  /**
   * JobStatus enums.
   */
  [JobStatusKey.SCHEDULED]: 'Scheduled',
  [JobStatusKey.PROCESSING]: 'Processing',
  [JobStatusKey.RESCHEDULED]: 'Rescheduled',
  [JobStatusKey.REJECTED]: 'Rejected',
  [JobStatusKey.ACKNOWLEDGED]: 'Acknowledged',
  /**
   * PitbullInstanceStatus enums.
   */
  [PitbullInstanceStatusKey.WAITING_FOR_HOST]: 'Waiting for host',
  [PitbullInstanceStatusKey.HOST_STARTING]: 'Host starting',
  [PitbullInstanceStatusKey.RUNNING]: 'Running',
  [PitbullInstanceStatusKey.COMPLETED]: 'Completed',
  [PitbullInstanceStatusKey.SUCCESS]: 'Success',
  [PitbullInstanceStatusKey.FAILED]: 'Failed',

  /**
   * PitbullStatus enums.
   */
   [PitbullStatusKey.WAITING]: 'Waiting',
   [PitbullStatusKey.RUNNING]: 'Running',
   [PitbullStatusKey.FINISHED]: 'Finished',
   [PitbullStatusKey.SUCCESS]: 'Success',

  /**
   * TokenGeneratorVersion enums.
   */
  [TokenGeneratorVersionKey.V1]: "V1",
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
