import React, { useContext } from 'react'
import { makeStyles } from 'tss-react/mui'
import Typography from '@mui/material/Typography'
import NewbieIcon from '@mui/icons-material/NewReleases'
import moment from 'moment'
import { Box } from '@mui/material'
import Avatar from '@/components/Avatar'
import {
  USER_AGE_CAUTION,
  USER_STATUS_BANNED,
  USER_STATUS_MAP_TEXT,
  USER_STATUS_SUSPENDED,
  USER_SUBSCRIPTION_BADGE_MODE,
} from '@/constants/user'
import Link from '@/components/Link'
import { retinaSrcSet } from '@/components/ItemImage'
import ChipLink from '@/components/ChipLink'
import { DOTABUFF_PROFILE_BASE_URL, STEAM_PROFILE_BASE_URL } from '@/constants/strings'
import { isDonationGlowExpired } from '@/service/api'
import AppContext from '@/components/AppContext'
import SubscriberBadge from '@/components/SubscriberBadge'
import { getUserBadgeFromBoons, getUserTagFromBoons } from '@/lib/badge'
import ExclusiveChip from '@/components/ExclusiveChip'
import type { UserSummary } from '@/lib/types'
import type { ReactNode } from 'react'

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

type ProfileUser = UserSummary & {
  market_stats?: {
    live: number
    reserved: number
    sold: number
    bid_completed: number
  }
  status?: number
}

interface ProfileCardProps {
  user: ProfileUser
  loading?: boolean
  hideSteamProfile?: boolean
  hideInventory?: boolean
  children?: ReactNode
}

export default function ProfileCard({
  user,
  loading,
  hideSteamProfile,
  hideInventory,
  ...other
}: ProfileCardProps) {
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
  const badgeType = (userBadge || undefined) as 'supporter' | 'trader' | 'partner' | undefined
  const tagType = (userTag || undefined) as
    | 'supporter'
    | 'trader'
    | 'partner'
    | 'middleman'
    | 'moderator'
    | undefined
  return (
    <div
      className={classes.details}
      style={
        isProfileReported ? { backgroundColor: '#2d0000', padding: 10, width: '100%' } : undefined
      }
    >
      <a href={storeProfile} target="_blank" rel="noreferrer noopener">
        <Avatar
          large
          badge={badgeType}
          className={classes.avatar}
          glow={isDonationGlowExpired(user.donated_at || '')}
          {...retinaSrcSet(user.avatar, 100, 100)}
        />
      </a>
      <Box>
        <Typography
          className={classes.profileName}
          component="h3"
          variant="h4"
          color={isProfileReported ? 'error' : 'inherit'}
        >
          {user.name}
          {!USER_SUBSCRIPTION_BADGE_MODE && !isMobile && (
            <SubscriberBadge
              type={badgeType}
              style={{
                marginLeft: '0.375rem',
                marginTop: '0.375rem',
                position: 'absolute',
              }}
              size="medium"
            />
          )}
        </Typography>

        {!USER_SUBSCRIPTION_BADGE_MODE && Boolean(userBadge) && isMobile && (
          <div>
            <SubscriberBadge
              type={badgeType}
              style={{ marginTop: '0.375rem', marginBottom: '0.575rem' }}
              size="medium"
            />
          </div>
        )}

        {isProfileReported && (
          <Typography color="error">{USER_STATUS_MAP_TEXT[user.status || 0]}</Typography>
        )}

        <Typography variant="body2" component="p" color="textSecondary">
          <Typography component="span" variant="caption">
            {user.steam_id}
          </Typography>{' '}
          &middot; Joined {moment(user.created_at).fromNow()}{' '}
          {moment().diff(moment(user.created_at), 'days') <= USER_AGE_CAUTION && (
            <NewbieIcon color="info" fontSize="inherit" sx={{ mb: -0.3 }} />
          )}
        </Typography>
        <Typography variant="body2" component="p">
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
          </Link>
        </Typography>

        <Box sx={{ '& > *': { mt: 0.5 } }}>
          {USER_SUBSCRIPTION_BADGE_MODE && userBadge && (
            <>
              <ExclusiveChip tag={badgeType} />
              &nbsp;
            </>
          )}
          {userTag && (
            <>
              <ExclusiveChip tag={tagType} />
              &nbsp;
            </>
          )}
          {!hideSteamProfile && (
            <>
              <ChipLink label="Steam Profile" href={steamProfileURL} />
              &nbsp;
            </>
          )}
          {!hideInventory && (
            <>
              <ChipLink label="Steam Inventory" href={dota2Inventory} />
              &nbsp;
            </>
          )}
          <ChipLink label="Dotabuff" href={`${DOTABUFF_PROFILE_BASE_URL}/${user.steam_id}`} />
          {other.children}
        </Box>
      </Box>
    </div>
  )
}
