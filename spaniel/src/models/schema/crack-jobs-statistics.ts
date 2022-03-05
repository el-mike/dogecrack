export type CrackJobsStatisticsDto = {
  all: number;
  acknowledged: number;
  processing: number;
  queued: number;
  rejected: number;
};

export type CrackJobsStatistics = CrackJobsStatisticsDto;
