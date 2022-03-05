import styled from 'styled-components';

export type JobStatusLightProps = {
  $size?: 'small' | 'medium';
  $queued?: boolean;
  $processing?: boolean;
  $rejected?: boolean;
  $acknowledged?: boolean;
};

export const JobStatusLight = styled.span<JobStatusLightProps>`
  display: flex;
  width: ${props => props.theme.spacing(props.$size === 'small' ? 1 : 2)};
  height: ${props => props.theme.spacing(props.$size === 'small' ? 1 : 2)};
  border-radius: 50%;

  ${props => props.$queued && `
    background-color: ${props.theme.palette.warning.light};
  `}

  ${props => props.$processing && `
    background-color: ${props.theme.palette.info.light};
  `}

  ${props => props.$rejected && `
    background-color: ${props.theme.palette.error.light};
  `}

  ${props => props.$acknowledged && `
    background-color: ${props.theme.palette.success.light};
  `}
`;
