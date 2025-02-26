import {
  BaseEntityDto,
  BaseEntity,
} from './base-entity';

import {
  PitbullInstanceDto,
  PitbullInstance,
  mapPitbullInstance,
} from './pitbull-instance';

export type CrackJobDto = BaseEntityDto & {
  walletString: string;

  name?: string;
  keyword?: string;
  passlistUrl?: string;
  tokenlist?: string;
  customTokenlist?: boolean;
  tokenGeneratorVersion?: number;

  status: number;

  instanceId: string;
  instance: PitbullInstanceDto;

  startedAt: string;
  firstScheduledAt: string;
  lastScheduledAt: string;
  acknowledgedAt: string;
  rejectedAt: string;

  rescheduleCount: number;

  errorLog: string;
};

export type CrackJob = BaseEntity
  & Omit<CrackJobDto, 'instance'>
  & {
    instance: PitbullInstance;
  };

export type CrackJobsFilters = Partial<{
  statuses: number[];
  keyword: string;
  passlistUrl: string;
  jobId: string;
  name: string;
  tokenGeneratorVersion: number;
}>;

export type RunCrackJobPayload = Partial<{
  name: string;
  keywords: string[];
  tokenlist: string;
  passlistUrl: string;
  tokenGeneratorVersion: number;
}>;

export type CancelCrackJobPayload = {
  jobId: string;
};

export type RecreateCrackJobPayload = {
  jobId: string;
};

export const mapCrackJob = (dto: CrackJobDto) => ({
  ...dto,
  instance: mapPitbullInstance(dto.instance || {}),
} as CrackJob);
