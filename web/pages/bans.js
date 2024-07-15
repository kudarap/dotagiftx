import React from 'react'
import PropTypes from 'prop-types'
import useSWR from 'swr'
import moment from 'moment'
import Head from 'next/head'
import { useRouter } from 'next/router'
import debounce from 'lodash/debounce'
import startsWith from 'lodash/startsWith'
import { makeStyles } from 'tss-react/mui'
import Typography from '@mui/material/Typography'
import TextField from '@mui/material/TextField'
import {
  APP_NAME,
  DOTABUFF_PROFILE_BASE_URL,
  STEAM_PROFILE_BASE_URL,
  STEAMREP_PROFILE_BASE_URL,
} from '@/constants/strings'
import Footer from '@/components/Footer'
import Header from '@/components/Header'
import Container from '@/components/Container'
import Link from '@/components/Link'
import Avatar from '@/components/Avatar'
import { BLACKLIST, fetcherBase, parseParams } from '@/service/api'
import { retinaSrcSet } from '@/components/ItemImage'
import { USER_STATUS_MAP_LABEL, USER_STATUS_MAP_COLOR } from '@/constants/user'

const useStyles = makeStyles()(theme => ({
  main: {
    [theme.breakpoints.down('md')]: {
      marginTop: theme.spacing(2),
    },
    marginTop: theme.spacing(4),
  },
}))

const filter = {
  sort: 'updated_at:desc',
  limit: 100,
}

const STEAMURL = 'https://steamcommunity.com'

function cleanURL(url = '') {
  const s = url.split('/')
  if (s.length < 5) {
    return url
  }

  return s.slice(0, 5).join('/')
}

function isVanityURL(url = '') {
  return url.startsWith(`${STEAMURL}/id/`)
}

// returns Steam ID when available and
// resolves URL when its a vanity/custom for auto-resolve profile.
function resolveProfileURL(url = '') {
  if (url === '') {
    return false
  }

  const u = cleanURL(url)
  if (!isVanityURL(u)) {
    return u
  }

  return u.replaceAll(STEAMURL, '')
}

export default function Blacklist() {
  const { classes } = useStyles()

  const [query, setQuery] = React.useState('')
  filter.q = query
  const url = parseParams(BLACKLIST, filter)
  const { data, error } = useSWR(url, fetcherBase)

  const router = useRouter()
  let resolvedQuery = false
  if (startsWith(query, STEAMURL, 0)) {
    resolvedQuery = resolveProfileURL(query)
    if (isVanityURL(query)) {
      router.push(resolvedQuery)
    }
  }

  return (
    <>
      <Head>
        <meta charset="UTF-8" />
        <title>{APP_NAME} :: Banned users</title>
      </Head>

      <Header />

      <main className={classes.main}>
        <Container>
          <Typography variant="h5" component="h1" gutterBottom>
            Banned users
          </Typography>
          <Typography>
            These accounts were flagged as <strong>banned</strong> or <strong>suspended</strong> due
            to abusive behaviour or account involvement to a scam incident.
          </Typography>
          <br />

          <SearchBar
            placeholder="Search by Steam ID or Steam Profile URL"
            onInput={v => setQuery(v)}
            helperText="eg. 76561198088587178 or https://steamcommunity.com/id/kudarap"
          />
          <br />
          <br />
          {error && <Typography color="error">Could not load blacklisted users</Typography>}
          {!data && !error && <Typography>Loading...</Typography>}
          {!error && data && data.map(user => <UserCard data={user} />)}
          {!error && data && data.length === 0 && resolvedQuery && (
            <Typography>Please wait. Redirecting to profile...</Typography>
          )}
        </Container>
      </main>

      <Footer />
    </>
  )
}

function SearchBar({ onInput, ...other }) {
  const debounceSearch = React.useCallback(debounce(onInput, 500), [])

  const [value, setValue] = React.useState('')
  const handleInput = e => {
    const v = e.target.value
    setValue(v)
    debounceSearch(v)
  }

  return (
    <TextField
      fullWidth
      variant="outlined"
      color="secondary"
      {...other}
      value={value}
      onInput={handleInput}
    />
  )
}
SearchBar.propTypes = {
  onInput: PropTypes.func,
}
SearchBar.defaultProps = {
  onInput: () => {},
}

function UserCard({ data }) {
  return (
    <div style={{ display: 'flex', marginBottom: 14 }}>
      <Avatar
        style={{ marginTop: 4 }}
        {...retinaSrcSet(data.avatar, 40, 40)}
        component={Link}
        href={`/profiles/${data.steam_id}`}
      />
      <div style={{ marginLeft: 8 }}>
        <Typography>
          {/* <strong>{data.name}</strong> */}
          <Typography color="textSecondary" variant="body2">
            SteamID: {`${data.steam_id}`}
            {` `}
            <span
              style={{
                padding: '2px 6px',
                color: 'white',
                background: USER_STATUS_MAP_COLOR[data.status],
                marginTop: -2,
                fontSize: '0.785em',
                fontWeight: 500,
              }}>
              {USER_STATUS_MAP_LABEL[data.status]} {moment(data.updated_at).fromNow()}
            </span>
          </Typography>
          <Link variant="body2" href={`/profiles/${data.steam_id}`}>
            Profile
          </Link>
          &nbsp;&middot;&nbsp;
          <Link
            variant="body2"
            gutterBottom
            target="_blank"
            rel="noreferrer noopener"
            href={`${STEAM_PROFILE_BASE_URL}/${data.steam_id}`}>
            Steam Profile
          </Link>
          &nbsp;&middot;&nbsp;
          <Link
            variant="body2"
            gutterBottom
            target="_blank"
            rel="noreferrer noopener"
            href={`${STEAMREP_PROFILE_BASE_URL}/${data.steam_id}`}>
            SteamRep
          </Link>
          &nbsp;&middot;&nbsp;
          <Link
            variant="body2"
            gutterBottom
            target="_blank"
            rel="noreferrer noopener"
            href={`${DOTABUFF_PROFILE_BASE_URL}/${data.steam_id}`}>
            Dotabuff
          </Link>
        </Typography>
      </div>
    </div>
  )
}
UserCard.propTypes = {
  data: PropTypes.object.isRequired,
}
