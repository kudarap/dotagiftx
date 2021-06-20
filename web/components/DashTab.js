import React from 'react'
import { withStyles } from '@material-ui/core/styles'
import Badge from '@material-ui/core/Badge'
import Tab from '@material-ui/core/Tab'

const StyledTab = withStyles(theme => ({
  root: {
    textTransform: 'none',
    color: theme.palette.text.primary,
    fontWeight: theme.typography.fontWeightRegular,
    fontSize: theme.typography.pxToRem(14),
    // marginRight: theme.spacing(1),
    '&:focus': {
      opacity: 1,
    },
    minWidth: 120,
  },
}))(props => <Tab disableRipple {...props} />)

const StyledBadge = withStyles(theme => ({
  badge: {
    top: 12,
    position: 'relative',
    border: `2px solid ${theme.palette.background.paper}`,
    padding: '0 4px',
  },
}))(Badge)

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
