import {
  useLocation,
  Navigate,
} from 'react-router-dom';

import { useAuthContext } from '../auth.context';

export const PublicRoute: React.FC = props => {
  const { user } = useAuthContext();
  const location = useLocation();

  return !!user
    ? <Navigate to='/dashboard' state={{ from: location }} replace />
    : <>{props.children}</>;
};
