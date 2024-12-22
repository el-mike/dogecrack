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
  keyword?: string;
  passlistUrl?: string;
  tokens?: string[];

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
  jobId: string;
}>;

export type RunCrackJobPayload = {
  keyword?: string;
  passlistUrl?: string;
};

export const mapCrackJob = (dto: CrackJobDto) => ({
  ...dto,
  instance: mapPitbullInstance(dto.instance || {}),
} as CrackJob);
