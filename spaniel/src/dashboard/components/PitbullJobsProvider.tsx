import {
  useState,
  useEffect,
  useMemo,
} from 'react';

import {
  ListRequest,
  PitbullJob,
  PitbullJobsFilters,
  JobStatusKey,
} from 'models';

import { useGeneralContext } from 'core/contexts';

import { usePitbullJobService } from 'core/hooks';

import {
  pitbullJobsContext,
  PitbullJobsContext,
} from '../pitbull-jobs.context';

const DEFAULT_PAGE_SIZE = 10;

export const PitbullJobsProvider: React.FC = props => {
  const { enums } = useGeneralContext();

  const pitbullJobService = usePitbullJobService();

  const defaultRequest = useMemo(
    () => ({
      page: 1,
      pageSize: DEFAULT_PAGE_SIZE,
      statuses: [enums.jobStatus[JobStatusKey.JOB_PROCESSING]],
    } as ListRequest<PitbullJobsFilters>),
    /* eslint-disable-next-line */
    [],
  );

  const [request, setRequest] = useState<ListRequest<PitbullJobsFilters>>(defaultRequest);
  const [totalCount, setTotalCount] = useState(0);
  const [jobs, setJobs] = useState<PitbullJob[]>([]);
  const [loading, setLoading] = useState(false);
  const [lastLoaded, setLastLoaded] = useState(new Date());

  const reload = () => {
    /**
     * We fire a new useEffect call by changing the reference to request object.
     */
    setRequest({ ...request });
    setLastLoaded(new Date());
  };

  const filter = (filters: PitbullJobsFilters) => {
    console.log();
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
  } as PitbullJobsContext;

  /**
   * We want to refresh the list every time request or lastLoaded change.
   */
  useEffect(() => {
    setLoading(true);
  
    pitbullJobService.getJobs(request)
      .then(result => {
        setJobs(result.entities || []);
        setTotalCount(result.totalCount);
      })
      .finally(() => setLoading(false));
      /* eslint-disable-next-line */
  }, [request]);

  return (
    <pitbullJobsContext.Provider value={value}>{props.children}</pitbullJobsContext.Provider>
  );
};
