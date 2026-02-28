import React from 'react'
import Badge from '@mui/material/Badge'
import { styled } from '@mui/material/styles'
import Tab from '@mui/material/Tab'

const StyledTab = styled(props => <Tab {...props} disableRipple />)(({ theme }) => ({
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

export default function DashTab(props) {
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
