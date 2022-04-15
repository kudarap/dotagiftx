import React, { useState } from 'react'
import PropTypes from 'prop-types'
import Alert from '@mui/material/Alert'
import Link from './Link'
import useLocalStorage from './useLocalStorage'

const targetUpdateID = 20220415

function ExpiringPostsBanner({ userID }) {
  const wuid = `whatsnew_id_${userID}`
  const [clientUpdateID, setClientUpdateID] = useLocalStorage(wuid, 0)

  const [open, setOpen] = useState(targetUpdateID > clientUpdateID)
  const handleClose = () => {
    setClientUpdateID(targetUpdateID)
    setOpen(false)
  }

  const handleSubmit = () => {
    handleClose()
  }

  if (!open) {
    return null
  }

  return (
    <Alert severity="warning" onClose={handleSubmit}>
      Major update: We will role out Expiring items on May 1, 2022 —{' '}
      <Link href="/expiring-posts">Read more</Link>
    </Alert>
  )
}
ExpiringPostsBanner.propTypes = {
  userID: PropTypes.string,
}
ExpiringPostsBanner.defaultProps = {
  userID: '',
}

export default ExpiringPostsBanner
