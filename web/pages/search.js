import React from 'react'
import PropTypes from 'prop-types'
import Head from 'next/head'
import { makeStyles } from '@material-ui/core/styles'
import Typography from '@material-ui/core/Typography'
import LinearProgress from '@material-ui/core/LinearProgress'
import { catalogSearch } from '@/service/api'
import Footer from '@/components/Footer'
import Header from '@/components/Header'
import Container from '@/components/Container'
import CatalogList from '@/components/CatalogList'
import TablePaginationRouter from '@/components/TablePaginationRouter'
import { APP_NAME, APP_URL } from '@/constants/strings'

const useStyles = makeStyles(theme => ({
  main: {
    marginTop: theme.spacing(2.5),
  },
  listControl: {},
  paginator: {
    float: 'right',
  },
}))

export default function Search({ catalogs: initialCatalogs, filter, canonicalURL }) {
  const classes = useStyles()

  const [catalogs, setCatalogs] = React.useState(initialCatalogs)
  const [error, setError] = React.useState(null)

  // Handle catalog request on page change.
  React.useEffect(() => {
    ;(async () => {
      try {
        const res = await catalogSearch(filter)
        setCatalogs(res)
      } catch (e) {
        setError(e.message)
      }
    })()
  }, [filter])

  let metaTitle = `${APP_NAME} Search`
  let metaDesc = `Search for item name, hero, treasure`
  if (filter.q) {
    metaTitle += ` :: ${filter.q}`
    metaDesc = `${catalogs && catalogs.total_count} results for "${filter.q}"`
  }

  const linkProps = { href: '/search', query: filter }

  return (
    <>
      <Head>
        <title>{metaTitle}</title>
        <meta name="description" content={metaDesc} />
        <link rel="canonical" href={canonicalURL} />
      </Head>
      <Header />

      <main className={classes.main}>
        <Container>
          {filter.q && (
            <>
              <Typography component="h1" variant="h6">
                {catalogs && catalogs.total_count} results for &quot;{filter.q}&quot;
              </Typography>
              <br />
            </>
          )}

          {!catalogs && <LinearProgress color="secondary" />}
          {catalogs && (
            <div>
              <CatalogList items={catalogs.data} error={error} />
              {!error && (
                <TablePaginationRouter
                  linkProps={linkProps}
                  colSpan={3}
                  style={{ textAlign: 'right' }}
                  count={catalogs.total_count}
                  page={filter.page}
                />
              )}
            </div>
          )}
        </Container>
      </main>

      <Footer />
    </>
  )
}
Search.propTypes = {
  catalogs: PropTypes.object.isRequired,
  filter: PropTypes.object.isRequired,
  canonicalURL: PropTypes.string.isRequired,
}

const catalogSearchFilter = { sort: 'popular', page: 1 }

// This gets called on every request
export async function getServerSideProps({ query }) {
  const filter = { ...catalogSearchFilter, ...query }
  filter.page = Number(query.page || 1)

  let catalogs = {}
  let error = null
  try {
    catalogs = await catalogSearch(filter)
  } catch (e) {
    error = e.message
  }

  const canonicalURL = `${APP_URL}/search?q=${filter.q}`

  return {
    props: {
      canonicalURL,
      filter,
      catalogs,
      error,
    },
  }
}
