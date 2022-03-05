import { Dictionary } from '../structs';

import { FromKeys } from '../types';

export enum JobStatusKey {
  SCHEDULED = 'SCHEDULED',
  PROCESSING = 'PROCESSING',
  RESCHEDULED = 'RESCHEDULED',
  REJECTED = 'REJECTED',
  ACKNOWLEDGED = 'ACKNOWLEDGED',
}

export enum PitbullInstanceStatusKey {
  WAITING_FOR_HOST = 'WAITING_FOR_HOST',
  HOST_STARTING = 'HOST_STARTING',
  RUNNING = 'RUNNING',
  COMPLETED = 'COMPLETED',
  FAILED = 'FAILED',
  SUCCESS = 'SUCCESS',
}

export enum PitbullStatusKey {
  WAITING = 'WAITING',
  RUNNING = 'RUNNING',
  FINISHED = 'FINISHED',
  SUCCESS = 'SUCCESS',
}

export type Enum = Dictionary<number>;

export type JobStatusEnum = FromKeys<JobStatusKey, number>;
export type PitbullInstanceStatusEnum = FromKeys<PitbullInstanceStatusKey, number>;
export type PitbullStatusEnum = FromKeys<PitbullStatusKey, number>;

export type Enums = {
  jobStatus: JobStatusEnum;
  pitbullInstanceStatus: PitbullInstanceStatusEnum;
  pitbullStatus: PitbullStatusEnum;
};
