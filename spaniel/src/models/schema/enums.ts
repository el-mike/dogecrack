import { Dictionary } from '../structs';

import { FromKeys } from '../types';

export enum JobStatusKey {
  JOB_SCHEDULED = 'JOB_SCHEDULED',
  JOB_PROCESSING = 'JOB_PROCESSING',
  JOB_RESCHEDULED = 'JOB_RESCHEDULED',
  JOB_REJECTED = 'JOB_REJECTED',
  JOB_ACKNOWLEDGED = 'JOB_ACKNOWLEDGED',
}

export enum PitbullInstanceStatusKey {
  WAITING_FOR_HOST = 'WAITING_FOR_HOST',
  HOST_STARTING = 'HOST_STARTING',
  WAITING = 'WAITING',
  RUNNING = 'RUNNING',
  FINISHED = 'FINISHED',
  SUCCESS = 'SUCCESS',
  INTERRUPTED = 'INTERRUPTED',
  FAILED = 'FAILED',
}

export type Enum = Dictionary<number>;

export type JobStatusEnum = FromKeys<JobStatusKey, number>;
export type PitbullInstanceStatusEnum = FromKeys<PitbullInstanceStatusKey, number>;

export type Enums = {
  jobStatus: JobStatusEnum;
  pitbullInstanceStatus: PitbullInstanceStatusEnum;
};
