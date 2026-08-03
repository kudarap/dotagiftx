import React from 'react'
import filter from 'lodash/filter'
import { useRouter } from 'next/router'
import TextField from '@mui/material/TextField'
import Autocomplete from '@mui/material/Autocomplete'
import CircularProgress from '@mui/material/CircularProgress'
import { item, itemSearch } from '@/service/api'
import type { Item, SearchResponse } from '@/lib/types'

// const itemSearchFilter = { limit: 1000, sort: 'view_count:desc', active: true }
const itemSearchFilter = { limit: 1000, sort: 'created_at:desc', active: true }
const optionTextSeparator = ' - '

type ItemOption = Item & { text: string }

interface ItemAutoCompleteProps {
  onSelect: (item: Partial<Item>) => void
  forwardedRef?: React.Ref<HTMLInputElement>
  required?: boolean
  disabled?: boolean
}

function ItemAutoComplete({ onSelect, forwardedRef, required, ...other }: ItemAutoCompleteProps) {
  const [open, setOpen] = React.useState(false)
  const [options, setOptions] = React.useState<ItemOption[]>([])
  const [value, setValue] = React.useState('')
  const loading = open && options.length === 0

  const router = useRouter()
  const itemSlug = router.query.s
  React.useEffect(() => {
    if (!itemSlug) {
      return
    }

    ;(async () => {
      try {
        const res = (await item(String(itemSlug))) as Item
        setValue(res.name)
        onSelect(res)
      } catch (e) {
        console.log('error getting item details', (e as Error).message)
      }
    })()

    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [itemSlug])

  React.useEffect(() => {
    ;(async () => {
      try {
        const res = (await itemSearch(itemSearchFilter)) as SearchResponse<Item>
        setOptions(res.data.map(i => ({ ...i, text: `${i.hero}${optionTextSeparator}${i.name}` })))
      } catch (e) {
        console.log('error getting item list', (e as Error).message)
      }
    })()
  }, [])

  const handleInputChange = (e: React.SyntheticEvent, text: string) => {
    setValue(text)
    const name = String(text).split(optionTextSeparator)[1]
    const res = filter(options, { name })
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
      isOptionEqualToValue={(opt, val) => opt.name === val.name}
      getOptionLabel={option => option.text}
      options={options}
      loading={loading}
      renderInput={params => (
        <TextField
          {...params}
          ref={forwardedRef}
          required={required}
          color="secondary"
          label="Item name"
          helperText="Search item you want to post from your inventory."
          variant="outlined"
          slotProps={{
            ...params.slotProps,

            input: {
              ...params.slotProps.input,
              endAdornment: (
                <>
                  {loading ? <CircularProgress color="inherit" size={20} /> : null}
                  {params.slotProps.input.endAdornment}
                </>
              ),
            },
          }}
        />
      )}
      {...other}
    />
  )
}

const itemAutoCompleteRef = (props: ItemAutoCompleteProps, ref: React.Ref<HTMLInputElement>) => (
  <ItemAutoComplete forwardedRef={ref} {...props} />
)

export default React.forwardRef(itemAutoCompleteRef)
