import {
  Typography,
  Accordion as MuiAccordion,
  AccordionSummary as  MuiAccordionSummary,
  AccordionDetails as MuiAccordionDetails,
} from '@mui/material';

import { ExpandMore as ExpandMoreIcon } from '@mui/icons-material';

export type AccordionProps = {
  title: string;
  detailsComponent?: React.FC;
}

export const Accordion: React.FC<AccordionProps> = props => {
  const { title, detailsComponent } = props;

  const Details = detailsComponent || MuiAccordionDetails;

  return (
      <MuiAccordion disableGutters>
        <MuiAccordionSummary expandIcon={<ExpandMoreIcon />}>
          <Typography variant='overline'>{title}</Typography>
        </MuiAccordionSummary>

        <Details>
          {props.children}
        </Details>
      </MuiAccordion>
  );
};

