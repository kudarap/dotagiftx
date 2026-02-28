import React, { useContext } from 'react'
import PropTypes from 'prop-types'
import { makeStyles } from 'tss-react/mui'
import ItemImage from '@/components/ItemImage'
import AppContext from '@/components/AppContext'

const useStyles = makeStyles()(theme => ({
  root: {
    [theme.breakpoints.down('sm')]: {
      background: 'rgba(0, 0, 0, 0.15)',
    },
  },
  media: {
    [theme.breakpoints.down('sm')]: {
      margin: '0 auto !important',
      width: 300,
      height: 170,
    },
    width: 165,
    height: 110,
    margin: theme.spacing(0, 1.5, 1.5, 0),
  },
}))

export default function ItemImageDialog({ item }) {
  const { classes } = useStyles()
  const { isMobile } = useContext(AppContext)

  let width = 165
  let height = 110
  if (isMobile) {
    width = 300
    height = 170
  }

  return (
    <div className={classes.root}>
      <ItemImage
        className={classes.media}
        image={item.image}
        title={item.name}
        rarity={item.rarity}
        width={width}
        height={height}
      />
    </div>
  )
}
ItemImageDialog.propTypes = {
  item: PropTypes.object.isRequired,
}
