import {
  useState,
  useEffect,
  useMemo,
} from 'react';

import {
  ListRequest,
  CrackJob,
  CrackJobsFilters,
  JobStatusKey,
  RunCrackJobPayload,
} from 'models';

import { useGeneralContext } from 'core/contexts';

import { useCrackJobService } from 'core/hooks';

import {
  crackJobsContext,
  CrackJobsContext,
} from '../crack-jobs.context';

const DEFAULT_PAGE_SIZE = 10;

export const CrackJobsProvider: React.FC = props => {
  const { enums } = useGeneralContext();

  const crackJobService = useCrackJobService();

  const defaultRequest = useMemo(
    () => ({
      page: 1,
      pageSize: DEFAULT_PAGE_SIZE,
      statuses: [enums.jobStatus[JobStatusKey.PROCESSING]],
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

  const run = (keyword: RunCrackJobPayload) => {
    crackJobService.runJob(keyword)
      .then(() => reload());
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
