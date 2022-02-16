import styled from 'styled-components';

import {
  Typography,
  IconButton
} from '@mui/material';

import { ContentCopy as ContentCopyIcon } from '@mui/icons-material';

export type LabeledInfoProps = {
  title: string;
  value?: string | number;
  toCopy?: LabeledInfoProps['value'];
};

const InfoContainer = styled.div``;

export const LabeledInfo: React.FC<LabeledInfoProps> = props => {
  const { title, value, toCopy: valueToCopy } = props;

  const handleCopy = () => {
    navigator.clipboard.writeText(`${valueToCopy}`);
  };

  return (
    <InfoContainer>
      <Typography variant='caption' display='flex'>{title}</Typography>
      <Typography variant='subtitle1' fontWeight='bold' display={valueToCopy ? 'inline' : 'flex'}>
        {props.children || value || '-'}
      </Typography>
      
      {!!valueToCopy && (
        <IconButton onClick={handleCopy} size='small'>
          <ContentCopyIcon fontSize='small' />
      </IconButton>
      )}
    </InfoContainer>
  );
};
