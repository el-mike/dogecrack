import { useTimeAgo } from 'core/hooks';

import  { LabeledInfo } from './LabeledInfo';

export type LastUpdatedProps = {
  updatedAt: string;
}

export const LastUpdated: React.FC<LastUpdatedProps> = props => {
  const { updatedAt } = props;

  const value = useTimeAgo(updatedAt);

  return (
    <LabeledInfo
      title='Last updated:'
      value={value}
    />
  );
};

