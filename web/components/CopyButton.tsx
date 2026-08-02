import React, { useState } from 'react'
import IconButton from '@mui/material/IconButton'
import type { IconButtonProps } from '@mui/material/IconButton'
import Tooltip from '@mui/material/Tooltip'
import ContentCopyIcon from '@mui/icons-material/ContentCopy'
import CheckIcon from '@mui/icons-material/Done'

interface CopyButtonProps extends IconButtonProps {
  value: string
}

export default function CopyButton(props: CopyButtonProps) {
  const [copied, setCopied] = useState(false)
  const { value, ...buttonProps } = props
  const handleClick = () => {
    navigator.clipboard.writeText(value)
    setCopied(true)
  }
  return (
    <Tooltip title={copied ? 'Copied!' : 'Copy full reference id'}>
      <IconButton {...buttonProps} onClick={handleClick}>
        {copied ? <CheckIcon fontSize="inherit" /> : <ContentCopyIcon fontSize="inherit" />}
      </IconButton>
    </Tooltip>
  )
}
