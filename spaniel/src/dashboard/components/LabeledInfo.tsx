import styled from 'styled-components';

import {
  Typography,
  IconButton
} from '@mui/material';

import { ContentCopy as ContentCopyIcon } from '@mui/icons-material';

export type LabeledInfoProps = {
  title: string;
  value: string | number;
  allowCopy?: boolean;
};

const InfoContainer = styled.div``;

export const LabeledInfo: React.FC<LabeledInfoProps> = props => {
  const { title, value, allowCopy } = props;

  const handleCopy = () => {
    navigator.clipboard.writeText(`${value}`);
  };

  return (
    <InfoContainer>
      <Typography variant='caption' display='flex'>{title}</Typography>
      <Typography variant='subtitle1' fontWeight='bold' display={allowCopy ? 'inline' : 'flex'}>
        {value}
      </Typography>
      
      {!!allowCopy && (
        <IconButton onClick={handleCopy} size='small'>
          <ContentCopyIcon fontSize='small' />
      </IconButton>
      )}
    </InfoContainer>
  );
};
