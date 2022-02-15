import { BaseEntity } from './base-entity';

export type PitbullJob = BaseEntity & {
  walletString: string;
  passlistUrl: string;

  status: number;

  instanceId: string;
  instance: any;

  startedAt: string;
  firstScheduledAt: string;
  lastScheduledAt: string;
  acknowledgedAt: string;
  rejectedAt: string;

  rescheduleCount: number;
};
