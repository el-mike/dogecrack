import {
  useEffect,
  useMemo,
  useState,
} from 'react';

import {
  CrackJob,
  CrackJobsFilters,
  JobStatusKey,
  ListRequest,
  RunCrackJobPayload,
} from 'models';

import {
  NotificationVariant,
  useGeneralContext,
  useSnackbarContext
} from 'core/contexts';

import { useCrackJobService } from 'core/hooks';

import {
  crackJobsContext,
  CrackJobsContext,
} from '../crack-jobs.context';

const DEFAULT_PAGE_SIZE = 10;

export const CrackJobsProvider: React.FC = props => {
  const { enums, latestTokenGeneratorVersion } = useGeneralContext();
  const { notify } = useSnackbarContext();

  const crackJobService = useCrackJobService();

  const defaultRequest = useMemo(
    () => ({
      page: 1,
      pageSize: DEFAULT_PAGE_SIZE,
      statuses: [enums.jobStatus[JobStatusKey.PROCESSING]],
      tokenGeneratorVersion: latestTokenGeneratorVersion,
    } as ListRequest<CrackJobsFilters>),
    /* eslint-disable-next-line */
    [],
  );

  const [request, setRequest] = useState<ListRequest<CrackJobsFilters>>(defaultRequest);
  const [totalCount, setTotalCount] = useState(0);
  const [jobs, setJobs] = useState<CrackJob[]>([]);
  const [loading, setLoading] = useState(false);
  const [lastLoaded, setLastLoaded] = useState(new Date());

  const reload = () => {
    /**
     * We fire a new useEffect call by changing the reference to request object.
     */
    setRequest({ ...request });
    setLastLoaded(new Date());
  };

  const filter = (filters: CrackJobsFilters) => {
    setLastLoaded(new Date());
    setRequest({
      ...request,
      ...filters,
    });
  }

  const changePage = (page: number) => {
    setLastLoaded(new Date());
    setRequest({
      ...request,
      page,
    })
  };

  const run = (payload: RunCrackJobPayload) => {
    crackJobService.runJob(payload)
      .then(() => {
        notify({
          message: 'Job has been run',
          variant: NotificationVariant.SUCCESS,
        });
      })
      .finally(() => reload());
  };

  const cancel = (jobId: string) => {
    /* We don't need to set it back to false in this function, as it will be done in reload.  */
    setLoading(true);

    crackJobService.cancelJob({ jobId })
      .then(() => {
        notify({
          message: 'Job has been canceled',
          variant: NotificationVariant.SUCCESS,
        });
      })
      .catch(() => {
        notify({
          message: 'Canceling job failed',
          variant: NotificationVariant.ERROR,
        });
      })
      .finally(() => reload());
  };

  const recreate = (jobId: string) => {
    setLoading(true);

    crackJobService.recreateJob({ jobId })
      .then(() => {
        notify({
          message: 'Job has been recreated',
          variant: NotificationVariant.SUCCESS,
        });
      })
      .catch(() => {
        notify({
          message: 'Recreating job failed',
          variant: NotificationVariant.ERROR,
        });
      })
      .finally(() => reload());
  };

  const resetFilters = () => {
    setRequest(defaultRequest);
    setLastLoaded(new Date());
  };

  const { page, pageSize, ...filters } = request;

  const value = {
    filters,
    page,
    pageSize: DEFAULT_PAGE_SIZE,
    totalCount,
    loading,
    jobs,
    lastLoaded,
    reload,
    filter,
    changePage,
    run,
    cancel,
    recreate,
    resetFilters,
  } as CrackJobsContext;

  /**
   * We want to refresh the list every time request or lastLoaded change.
   */
  useEffect(() => {
    setLoading(true);

    crackJobService.getJobs(request)
      .then(result => {
        setJobs(result.entities || []);
        setTotalCount(result.totalCount);
      })
      .finally(() => setLoading(false));
      /* eslint-disable-next-line */
  }, [request]);

  return (
    <crackJobsContext.Provider value={value}>{props.children}</crackJobsContext.Provider>
  );
};
