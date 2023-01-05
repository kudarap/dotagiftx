import React from 'react'
import Typography from '@mui/material/Typography'
import MuiLink from '@mui/material/Link'
import Grid from '@mui/material/Grid'
import { makeStyles } from 'tss-react/mui'
import RarityTag from '@/components/RarityTag'
import ItemImage from '@/components/ItemImage'
import Link from '@/components/Link'
import ChipLink from '@/components/ChipLink'

const useStyles = makeStyles()(theme => ({
  mediaContainer: {
    [theme.breakpoints.down('sm')]: {
      width: '100%',
    },
  },
  media: {
    [theme.breakpoints.down('sm')]: {
      width: 300,
      height: 198,
      margin: '0 auto',
    },
    width: 165,
    height: 110,
  },
  itemStats: {
    marginBottom: theme.spacing(1),
  },
}))

export default function ItemViewCard({ item }) {
  const { classes } = useStyles()

  const wikiLink = `https://dota2.gamepedia.com/${item.name.replace(/ +/gi, '_')}`
  return (
    <Grid container spacing={1.5}>
      <Grid item className={classes.mediaContainer}>
        <div style={{ background: 'rgba(0, 0, 0, 0.15)' }}>
          {item.image && (
            <a href={wikiLink} target="_blank" rel="noreferrer noopener">
              <ItemImage
                className={classes.media}
                image={item.image}
                width={330}
                height={220}
                title={item.name}
                rarity={item.rarity}
                nextOptimized
              />
            </a>
          )}
        </div>
      </Grid>
      <Grid item>
        <div>
          <Typography component="h1" variant="h4">
            {item.name}
          </Typography>
        </div>
        <div>
          <Link href={`/search?origin=${item.origin}`}>{item.origin}</Link>{' '}
          {item.rarity !== 'regular' && (
            <>
              &mdash;
              <RarityTag
                rarity={item.rarity}
                variant="body1"
                component={Link}
                href={`/search?rarity=${item.rarity}`}
              />
            </>
          )}
        </div>
        <div>
          <Typography color="textSecondary" component="span">
            {`Hero: `}
          </Typography>
          <Link href={`/search?hero=${item.hero}`}>{item.hero}</Link>
        </div>
        <div className={classes.itemStats} spacing={1}>
          <ChipLink label="Dota 2 Wiki" href={wikiLink} />
          &nbsp;&middot;&nbsp;
          <Typography variant="body2" component={MuiLink} color="textPrimary" href="#reserved">
            {item.reserved_count} Reserved
          </Typography>
          &nbsp;&middot;&nbsp;
          <Typography variant="body2" component={MuiLink} color="textPrimary" href="#delivered">
            {item.sold_count} Delivered
          </Typography>
        </div>
      </Grid>
    </Grid>
  )
}
