import React, { useState } from 'react'
import IconButton from '@mui/material/IconButton'
import Tooltip from '@mui/material/Tooltip'
import ContentCopyIcon from '@mui/icons-material/ContentCopy'
import CheckIcon from '@mui/icons-material/Done'

export default function CopyButton(props) {
  const [copied, setCopied] = useState(false)
  const handleClick = () => {
    navigator.clipboard.writeText(props.value)
    setCopied(true)
  }
  return (
    <Tooltip title={copied ? 'Copied!' : 'Copy full reference id'}>
      <IconButton {...props} onClick={handleClick}>
        {copied ? <CheckIcon fontSize="inherit" /> : <ContentCopyIcon fontSize="inherit" />}
      </IconButton>
    </Tooltip>
  )
}
