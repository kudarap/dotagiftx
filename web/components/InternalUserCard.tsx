import React from 'react'
import Image from 'next/image'
import moment from 'moment'
import Typography from '@mui/material/Typography'
import Link from '@mui/material/Link'
import Box from '@mui/material/Box'
import ChipLink from '@/components/ChipLink'
import Avatar from '@/components/Avatar'
import ExclusiveChip, { tagSettings } from '@/components/ExclusiveChip'
import Button from '@/components/Button'
import DiscordIcon from '@/components/DiscordIcon'
import { CDN_URL } from '@/service/api'
import { getUserTagFromBoons } from '@/lib/badge'

export interface InternalUserCardProps {
  name?: string
  id?: string
  img?: string
  boons?: string[]
  discordURL?: string
  createdAt?: string
}

export default function InternalUserCard({
  name,
  id,
  img,
  boons,
  discordURL,
  createdAt,
}: InternalUserCardProps) {
  const userTag = getUserTagFromBoons(boons)
  if (!userTag) {
    return null
  }

  const { color } = tagSettings[userTag as keyof typeof tagSettings]

  return (
    <Box sx={{ display: 'inline-flex' }}>
      <Avatar
        src={`${CDN_URL}/${img}`}
        sx={{ width: 100, height: 100, border: `2px solid ${color}`, m: 1, mr: 2, ml: 0 }}
      />
      <Box sx={{ mr: 6, mb: 4 }}>
        <Typography variant="h5" component="h3">
          {name}
        </Typography>
        <Typography variant="body2" component="p" color="textSecondary">
          <Typography component="span" variant="caption">
            {id}
          </Typography>{' '}
          &middot; Joined {moment(createdAt).fromNow()}{' '}
        </Typography>
        <Box sx={{ mb: 1, mt: 1 }}>
          <ExclusiveChip tag={userTag as keyof typeof tagSettings} />
          &nbsp;
          <ChipLink href={`https://steamcommunity.com/profiles/${id}`} label="Steam Profile" />
        </Box>
        <Box>
          <Button
            startIcon={<DiscordIcon />}
            component={Link}
            nativeButton={false}
            target="_blank"
            rel="noreferrer noopener"
            size="small"
            href={discordURL}
          >
            Discord
          </Button>
          &nbsp;
          <Button
            startIcon={
              <Image
                src="/icon_2x.png"
                width={16}
                height={16}
                style={{ height: 16, width: 'auto', filter: 'brightness(10)' }}
                alt="dgx"
              />
            }
            component={Link}
            nativeButton={false}
            size="small"
            href={`/profiles/${id}`}
          >
            DotagiftX
          </Button>
        </Box>
      </Box>
    </Box>
  )
}
