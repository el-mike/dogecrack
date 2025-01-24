import { ReactNode } from 'react';

import styled from 'styled-components';

import {
  Typography,
  IconButton,
  TypographyProps
} from '@mui/material';

import { ContentCopy as ContentCopyIcon } from '@mui/icons-material';
import { Spacer } from './Spacer';

export type LabeledInfoProps = {
  title?: string;
  titleComponent?: ReactNode;
  value?: string | number;
  valueVariant?: TypographyProps['variant'];
  toCopy?: LabeledInfoProps['value'];
  useZero?: boolean;
  endAdornment?: ReactNode;
};

const InfoContainer = styled.div``;

const ValueContainer = styled.div`
  display: flex;
  align-items: center;
`;

export const LabeledInfo: React.FC<LabeledInfoProps> = props => {
  const {
    title,
    titleComponent,
    value,
    valueVariant,
    toCopy: valueToCopy,
    useZero,
    endAdornment,
  } = props;

  const handleCopy = () => {
    navigator.clipboard.writeText(`${valueToCopy}`);
  };

  return (
    <InfoContainer>
      {!!titleComponent
        ? titleComponent
        : <Typography variant='caption' display='flex'>{title}</Typography>
      }

      <ValueContainer>
        <Typography variant={valueVariant || 'subtitle1'} fontWeight='bold' display='flex' >
          {props.children || value || (!!useZero ? 0 : '-')}
        </Typography>

        {!!valueToCopy && (
          <IconButton onClick={handleCopy} size='small'>
            <ContentCopyIcon fontSize='small' />
        </IconButton>
        )}

        {!!endAdornment && (
          <>
            <Spacer mr={1} />
            {endAdornment}
          </>
        )}
      </ValueContainer>
    </InfoContainer>
  );
};
