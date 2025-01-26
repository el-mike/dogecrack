import {
  useCallback,
  useEffect,
  useState
} from 'react';

import { CheckedIdeas } from 'models';
import { useCrackJobService } from 'core/hooks';
import {
  NotificationVariant,
  useSnackbarContext
} from 'core/contexts';
import {
  checkedIdeasContext,
  CheckedIdeasContext,
} from '../checked-ideas.context';

export const CheckedIdeasProvider: React.FC = props => {
  const crackJobService = useCrackJobService();
  const { notify } = useSnackbarContext();

  const [checkedIdeas, setCheckedIdeas] = useState<CheckedIdeas | null>(null);
  const [loading, setLoading] = useState(false);

  const loadIdeas = useCallback(() => {
    setLoading(true);

    crackJobService.getCheckedIdeas()
      .then(ideas => setCheckedIdeas(ideas))
      .catch(() => notify({
        message: 'Loading checked ideas failed',
        variant: NotificationVariant.ERROR,
      }))
      .finally(() => setLoading(false));
  }, [crackJobService, setCheckedIdeas, setLoading, notify]);

  useEffect(() => {
    loadIdeas();
  }, []);

  const value = {
    loadIdeas,
    checkedIdeas,
    loading,
  } as CheckedIdeasContext;

  return <checkedIdeasContext.Provider value={value}>{props.children}</checkedIdeasContext.Provider>;
};
