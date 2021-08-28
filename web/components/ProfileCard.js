import React from 'react'
import PropTypes from 'prop-types'
import { makeStyles } from '@material-ui/core/styles'
import Typography from '@material-ui/core/Typography'
import Avatar from '@/components/Avatar'
import { USER_STATUS_MAP_TEXT } from '@/constants/user'
import Link from '@/components/Link'
import { retinaSrcSet } from '@/components/ItemImage'
import DonatorBadge from '@/components/DonatorBadge'
import { isDonationGlowExpired } from '@/service/api'

const useStyles = makeStyles(theme => ({
  details: {
    [theme.breakpoints.down('xs')]: {
      textAlign: 'center',
      display: 'block',
    },
    display: 'inline-flex',
  },
  profileName: {
    [theme.breakpoints.down('xs')]: {
      fontSize: theme.typography.h6.fontSize,
    },
    fontSize: '1.9vw',
  },
  avatar: {
    [theme.breakpoints.down('xs')]: {
      margin: '0 auto',
    },
    width: 100,
    height: 100,
    marginRight: theme.spacing(1.5),
  },
}))

export default function ProfileCard({ user, loading, ...other }) {
  const classes = useStyles()

  const storeProfile = `/profiles/${user.steam_id}`
  const marketSummary = user.market_stats

  const isProfileReported = Boolean(user.status)

  return (
    <div
      className={classes.details}
      style={isProfileReported ? { backgroundColor: '#2d0000', padding: 10, width: '100%' } : null}>
      <a href={storeProfile} target="_blank" rel="noreferrer noopener">
        <Avatar
          className={classes.avatar}
          glow={isDonationGlowExpired(user.donated_at)}
          {...retinaSrcSet(user.avatar, 100, 100)}
        />
      </a>
      <Typography component="h1">
        <Typography
          className={classes.profileName}
          component="p"
          variant="h4"
          color={isProfileReported ? 'error' : ''}>
          {user.name}
          {Boolean(user.donation) && (
            <DonatorBadge
              style={{ marginLeft: 4, marginTop: 10, position: 'absolute' }}
              size="medium">
              DONATOR
            </DonatorBadge>
          )}
        </Typography>

        {isProfileReported && (
          <Typography color="error">{USER_STATUS_MAP_TEXT[user.status]}</Typography>
        )}

        <Typography variant="body2" component="span">
          <Link href={`/profiles/${user.steam_id}`}>
            {marketSummary ? marketSummary.live : '--'} Items
          </Link>{' '}
          &middot;{' '}
          <Link href={`/profiles/${user.steam_id}/reserved`}>
            {marketSummary ? marketSummary.reserved : '--'} Reserved
          </Link>{' '}
          &middot;{' '}
          <Link href={`/profiles/${user.steam_id}/delivered`}>
            {marketSummary ? marketSummary.sold : '--'} Delivered
          </Link>{' '}
          &middot;{' '}
          <Link href={`/profiles/${user.steam_id}/bought`}>
            {marketSummary ? marketSummary.bid_completed : '--'} Bought
          </Link>{' '}
        </Typography>

        <br />
        <Typography gutterBottom>{other.children}</Typography>
      </Typography>
    </div>
  )
}
ProfileCard.propTypes = {
  user: PropTypes.object,
  marketSummary: PropTypes.object,
  loading: PropTypes.bool,
}
ProfileCard.defaultProps = {
  user: {},
  marketSummary: null,
  loading: false,
}
