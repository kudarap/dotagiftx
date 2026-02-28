import React from 'react'
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

export default function InternalUserCard({ name, id, img, boons, discordURL }) {
  const userTag = getUserTagFromBoons(boons)
  const { color } = tagSettings[userTag]

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
        <Box sx={{ mb: 1, mt: 1 }}>
          <ExclusiveChip tag={userTag} />
          &nbsp;
          <ChipLink href={`https://steamcommunity.com/profiles/${id}`} label="Steam Profile" />
          &nbsp;
          <ChipLink href={`https://steamrep.com/profiles/${id}`} label="SteamRep" />
        </Box>
        <Box>
          <Button
            startIcon={<DiscordIcon />}
            component={Link}
            target="_blank"
            rel="noreferrer noopener"
            size="small"
            href={discordURL}>
            Discord
          </Button>
          &nbsp;
          <Button
            startIcon={
              <img src="/icon_2x.png" style={{ height: 16, filter: 'brightness(10)' }} alt="dgx" />
            }
            component={Link}
            size="small"
            href={`/profiles/${id}`}>
            DotagiftX
          </Button>
        </Box>
      </Box>
    </Box>
  )
}
