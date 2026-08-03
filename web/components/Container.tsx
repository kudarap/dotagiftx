import React from 'react'
import type { CSSProperties, ReactNode } from 'react'
import { makeStyles } from 'tss-react/mui'
import MuiContainer from '@mui/material/Container'
import type { ContainerProps } from '@mui/material/Container'

import { APP_FOOTER_HEIGHT_TOTAL } from '@/constants/app'

const maxWidth = 1000

const useStyles = makeStyles()(theme => ({
  root: {
    padding: theme.spacing(0, 1),
    [theme.breakpoints.down('md')]: {
      padding: theme.spacing(1),
    },
    minHeight: 500,
  },
}))

interface ContainerPropsEx extends ContainerProps {
  children: ReactNode
  disableMinHeight?: boolean
}

export default function Container({ children, disableMinHeight, ...other }: ContainerPropsEx) {
  const { classes } = useStyles()

  const style: CSSProperties = {
    minHeight: disableMinHeight ? 0 : `calc(100vh - ${APP_FOOTER_HEIGHT_TOTAL}px)`,
    maxWidth: other.maxWidth ? undefined : maxWidth,
  }

  return (
    <MuiContainer className={classes.root} disableGutters style={style} {...other}>
      {children}
    </MuiContainer>
  )
}
