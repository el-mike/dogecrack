import { Outlet } from 'react-router-dom';

export const AuthenticatedLayout: React.FC = props => {
  return (
  <>
    <Outlet />
  </>
  );
};
