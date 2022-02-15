import {
  useLocation,
  Navigate,
} from 'react-router-dom';

import { useAuth } from '../auth.context';

export const ProtectedRoute: React.FC = props => {
  const { user } = useAuth();
  const location = useLocation();

  return !user
    ? <Navigate to='/login' state={{ from: location }} replace />
    : <>{props.children}</>;
};
