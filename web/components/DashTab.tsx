import React from 'react'
import Badge from '@mui/material/Badge'
import { styled } from '@mui/material/styles'
import Tab from '@mui/material/Tab'
import type { TabProps } from '@mui/material/Tab'
import type { ReactNode } from 'react'

const StyledTab = styled((props: TabProps) => <Tab {...props} disableRipple />)(({ theme }) => ({
  textTransform: 'none',
  color: theme.palette.text.primary,
  fontWeight: theme.typography.fontWeightRegular,
  fontSize: theme.typography.pxToRem(14),
  // marginRight: theme.spacing(1),
  '&:focus': {
    opacity: 1,
  },
  minWidth: 120,
}))

const StyledBadge = styled(Badge)(({ theme }) => ({
  '.MuiBadge-badge': {
    top: 10,
    position: 'relative',
    border: `2px solid ${theme.palette.background.paper}`,
    padding: '0 4px',
  },
}))

interface DashTabProps extends TabProps {
  label: string
  badgeContent?: ReactNode
}

export default function DashTab(props: DashTabProps) {
  const { label, badgeContent, ...other } = props
  return (
    <StyledTab
      {...other}
      label={
        badgeContent ? (
          <StyledBadge badgeContent={badgeContent} max={999}>
            {label}
          </StyledBadge>
        ) : (
          label
        )
      }
    />
  )
}
