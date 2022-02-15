import {
  useLocation,
  Navigate,
} from 'react-router-dom';

import { useAuth } from '../auth.context';

export const PublicRoute: React.FC = props => {
  const { user } = useAuth();
  const location = useLocation();

  return !!user
    ? <Navigate to='/dashboard' state={{ from: location }} replace />
    : <>{props.children}</>;
};
