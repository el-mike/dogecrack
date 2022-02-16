import {
  useLocation,
  Navigate,
} from 'react-router-dom';

import { useAuthContext } from '../auth.context';

export const ProtectedRoute: React.FC = props => {
  const { user } = useAuthContext();
  const location = useLocation();

  return !user
    ? <Navigate to='/login' state={{ from: location }} replace />
    : <>{props.children}</>;
};
