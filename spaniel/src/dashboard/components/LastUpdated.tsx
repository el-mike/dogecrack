import { TimeAgo } from 'core/components';

import  { LabeledInfo } from './LabeledInfo';

export type LastUpdatedProps = {
  updatedAt: string;
}

export const LastUpdated: React.FC<LastUpdatedProps> = props => {
  const { updatedAt } = props;

  return (
    <LabeledInfo title='Last updated:'>
      <TimeAgo from={updatedAt} />
    </LabeledInfo>
  );
};

