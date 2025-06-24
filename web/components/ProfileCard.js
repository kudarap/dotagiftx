import React, { useContext } from 'react'
import PropTypes from 'prop-types'
import { makeStyles } from 'tss-react/mui'
import Typography from '@mui/material/Typography'
import Avatar from '@/components/Avatar'
import { USER_STATUS_BANNED, USER_STATUS_MAP_TEXT, USER_STATUS_SUSPENDED } from '@/constants/user'
import Link from '@/components/Link'
import { retinaSrcSet } from '@/components/ItemImage'
import ChipLink from '@/components/ChipLink'
import {
  DOTABUFF_PROFILE_BASE_URL,
  STEAM_PROFILE_BASE_URL,
  STEAMREP_PROFILE_BASE_URL,
} from '@/constants/strings'
import { isDonationGlowExpired } from '@/service/api'
import AppContext from '@/components/AppContext'
import SubscriberBadge from './SubscriberBadge'
import { getUserBadgeFromBoons, getUserTagFromBoons } from '@/lib/badge'
import ExclusiveChip from '@/components/ExclusiveChip'

const useStyles = makeStyles()(theme => ({
  details: {
    [theme.breakpoints.down('sm')]: {
      textAlign: 'center',
      display: 'block',
    },
    display: 'inline-flex',
  },
  profileName: {
    [theme.breakpoints.down('sm')]: {
      fontSize: theme.typography.h6.fontSize,
    },
    fontSize: '1.8vw',
  },
  avatar: {
    [theme.breakpoints.down('sm')]: {
      margin: `0 auto ${theme.spacing(1)}`,
    },
    width: 100,
    height: 100,
    marginRight: theme.spacing(1.5),
  },
  badge: {},
}))

export default function ProfileCard({ user, loading, ...other }) {
  const { classes } = useStyles()

  const { isMobile } = useContext(AppContext)

  const storeProfile = `/profiles/${user.steam_id}`
  const steamProfileURL = `${STEAM_PROFILE_BASE_URL}/${user.steam_id}`
  const dota2Inventory = `${steamProfileURL}/inventory#570`
  const marketSummary = user.market_stats

  const isProfileReported =
    user.status === USER_STATUS_SUSPENDED || user.status === USER_STATUS_BANNED

  const userBadge = getUserBadgeFromBoons(user.boons)
  const userTag = getUserTagFromBoons(user.boons)
  return (
    <div
      className={classes.details}
      style={isProfileReported ? { backgroundColor: '#2d0000', padding: 10, width: '100%' } : null}>
      <a href={storeProfile} target="_blank" rel="noreferrer noopener">
        <Avatar
          large
          badge={userBadge}
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
          {userBadge && !isMobile && (
            <SubscriberBadge
              type={userBadge}
              style={{ marginLeft: '0.375rem', marginTop: '0.375rem', position: 'absolute' }}
              size="medium"
            />
          )}
        </Typography>

        {Boolean(userBadge) && isMobile && (
          <div>
            <SubscriberBadge
              type={userBadge}
              style={{ marginTop: '0.375rem', marginBottom: '0.575rem' }}
              size="medium"
            />
          </div>
        )}

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
        <Typography gutterBottom>
          {userTag && (
            <>
              <ExclusiveChip tag={userTag} />
              &nbsp;
            </>
          )}
          <ChipLink label="Steam Inventory" href={dota2Inventory} />
          &nbsp;
          <ChipLink label="SteamRep" href={`${STEAMREP_PROFILE_BASE_URL}/${user.steam_id}`} />
          &nbsp;
          <ChipLink label="Dotabuff" href={`${DOTABUFF_PROFILE_BASE_URL}/${user.steam_id}`} />
          {other.children}
        </Typography>
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
