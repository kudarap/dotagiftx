import React from 'react'
import PropTypes from 'prop-types'
import Head from 'next/head'
import debounce from 'lodash/debounce'
import { makeStyles } from '@material-ui/core/styles'
import Typography from '@material-ui/core/Typography'
import TextField from '@material-ui/core/TextField'
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
import useSWR from 'swr'
import { BLACKLIST, fetcherBase, parseParams } from '@/service/api'
import { retinaSrcSet } from '@/components/ItemImage'
import { USER_STATUS_MAP_LABEL } from '@/constants/user'
import moment from 'moment'

const useStyles = makeStyles(theme => ({
  main: {
    [theme.breakpoints.down('sm')]: {
      marginTop: theme.spacing(2),
    },
    marginTop: theme.spacing(4),
  },
}))

const filter = {
  sort: 'updated_at:desc',
  limit: 100,
}

export default function Blacklist() {
  const classes = useStyles()

  const [query, setQuery] = React.useState('')
  filter.q = query
  const url = parseParams(BLACKLIST, filter)
  const { data, error } = useSWR(url, fetcherBase)

  return (
    <>
      <Head>
        <title>{APP_NAME} :: Fraud</title>
      </Head>

      <Header />

      <main className={classes.main}>
        <Container>
          <Typography variant="h5" component="h1" gutterBottom>
            Blacklist
          </Typography>
          <Typography color="textSecondary">
            These accounts were flagged as banned or suspended due to scam incident and/or
            involvement to scam.
          </Typography>
          <br />

          <SearchBar placeholder="Search by Steam URL or ID..." onInput={v => setQuery(v)} />
          <br />
          <br />

          {error && <Typography color="error">Could not load blacklisted users</Typography>}
          {!data && !error && <Typography>Loading...</Typography>}
          {!error && data && data.map(user => <UserCard data={user} />)}
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
      size="small"
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
    <div style={{ display: 'flex', marginBottom: 12 }}>
      <Avatar style={{ marginTop: 4 }} {...retinaSrcSet(data.avatar, 40, 40)} />
      <div style={{ marginLeft: 8 }}>
        <Typography>
          <strong>{data.name}</strong>
          <span
            style={{
              padding: '2px 6px',
              color: 'white',
              background: '#a00',
              marginLeft: 4,
              fontSize: 10,
              fontWeight: 500,
            }}>
            {USER_STATUS_MAP_LABEL[data.status]} {moment(data.updated_at).fromNow()}
          </span>
          <Typography variant="body2" color="textSecondary">
            SteamID {data.steam_id}
          </Typography>
        </Typography>
        <Link variant="body2" href={`/profiles/${data.steam_id}/activity`}>
          History
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
      </div>
    </div>
  )
}
UserCard.propTypes = {
  data: PropTypes.object.isRequired,
}
