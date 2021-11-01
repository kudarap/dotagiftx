import React from 'react'
import PropTypes from 'prop-types'
import makeStyles from '@mui/styles/makeStyles'
import MuiContainer from '@mui/material/Container'

import { APP_FOOTER_HEIGHT_TOTAL } from '@/constants/app'

const useStyles = makeStyles(theme => ({
  root: {
    [theme.breakpoints.down('md')]: {
      padding: theme.spacing(1),
    },
  },
}))

export default function Container({ children, disableMinHeight, ...other }) {
  const classes = useStyles()

  return (
    <MuiContainer
      className={classes.root}
      maxWidth="md"
      disableGutters
      style={{ minHeight: disableMinHeight ? 0 : `calc(100vh - ${APP_FOOTER_HEIGHT_TOTAL}px)` }}
      {...other}>
      {children}
    </MuiContainer>
  )
}
Container.propTypes = {
  children: PropTypes.node.isRequired,
  disableMinHeight: PropTypes.bool,
}
Container.defaultProps = {
  disableMinHeight: false,
}
