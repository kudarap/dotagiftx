import React from 'react'
import PropTypes from 'prop-types'
import filter from 'lodash/filter'
import { useRouter } from 'next/router'
import TextField from '@material-ui/core/TextField'
import Autocomplete from '@material-ui/lab/Autocomplete'
import CircularProgress from '@material-ui/core/CircularProgress'
import { item, itemSearch } from '@/service/api'

const itemSearchFilter = { limit: 1000, sort: 'popular' }

export default function ItemAutoComplete({ onSelect }) {
  const [open, setOpen] = React.useState(false)
  const [options, setOptions] = React.useState([])
  const [value, setValue] = React.useState('')
  const loading = open && options.length === 0

  // React.useEffect(() => {
  //   let active = true
  //
  //   if (!loading) {
  //     return undefined
  //   }
  //
  //   ;(async () => {
  //     const catalogs = await itemSearch(itemSearchFilter)
  //
  //     if (active) {
  //       setOptions(catalogs.data)
  //     }
  //   })()
  //
  //   return () => {
  //     active = false
  //   }
  // }, [loading])
  //
  // React.useEffect(() => {
  //   if (!open) {
  //     setOptions([])
  //   }
  // }, [open])

  const router = useRouter()
  const itemSlug = router.query.s
  React.useEffect(() => {
    if (!itemSlug) {
      return
    }

    ;(async () => {
      try {
        const res = await item(itemSlug)
        setValue(res.name)
        onSelect(res)
      } catch (e) {
        console.log('error getting item details', e.message)
      }
    })()
  }, [itemSlug])

  React.useEffect(() => {
    ;(async () => {
      try {
        const res = await itemSearch(itemSearchFilter)
        setOptions(res.data)
      } catch (e) {
        console.log('error getting item list', e.message)
      }
    })()
  }, [])

  const handleInputChange = (e, text) => {
    setValue(text)
    const res = filter(options, { name: text })
    if (res.length === 0) {
      onSelect({})
      return
    }

    onSelect(res[0])
  }

  return (
    <Autocomplete
      id="asynchronous-item-search"
      fullWidth
      clearOnBlur={false}
      open={open}
      onOpen={() => {
        setOpen(true)
      }}
      onClose={() => {
        setOpen(false)
      }}
      onInputChange={handleInputChange}
      inputValue={value}
      getOptionSelected={(opt, val) => opt.name === val.name}
      getOptionLabel={option => option.name}
      options={options}
      loading={loading}
      renderInput={params => (
        <TextField
          {...params}
          color="secondary"
          label="Item name"
          helperText="Search item you want to post from your inventory."
          variant="outlined"
          InputProps={{
            ...params.InputProps,
            endAdornment: (
              <>
                {loading ? <CircularProgress color="inherit" size={20} /> : null}
                {params.InputProps.endAdornment}
              </>
            ),
          }}
        />
      )}
    />
  )
}
ItemAutoComplete.propTypes = {
  onSelect: PropTypes.func,
}
ItemAutoComplete.defaultProps = {
  onSelect: () => {},
}
