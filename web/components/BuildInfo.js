import fetch from 'unfetch'
import useSWR from 'swr'

const fetcher = url => fetch(url).then(r => r.json())

export default function BuildInfo() {
  const { data, error } = useSWR('https://fotolink.app/api', fetcher)

  if (error) return <div>failed to load</div>
  if (!data) return <div>loading...</div>

  return (
    <div className="container">
      <p>
        tag: {data.version} <br />
        hash: {data.hash} <br />
        built: {data.built} <br />
      </p>
    </div>
  )
}
