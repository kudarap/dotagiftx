import React from 'react'
import PropTypes from 'prop-types'
import has from 'lodash/has'

const tabPanelIndex = {}
export default function TabPanel(props) {
  const { children, value, index, ...other } = props

  // Check for indexed component, it will prevent render from
  // loading everything on mount.
  if (value !== index && !has(tabPanelIndex, index)) {
    return null
  }
  tabPanelIndex[index] = true

  return (
    <div hidden={value !== index} {...other}>
      {children}
    </div>
  )
}
TabPanel.propTypes = {
  children: PropTypes.node.isRequired,
  index: PropTypes.any.isRequired,
  value: PropTypes.any.isRequired,
}
