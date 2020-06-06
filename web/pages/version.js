import Head from 'next/head'

export default function Version({ data }) {
  return (
    <div className="container">
      <Head>
        <title>User head</title>
        <link rel="icon" href="/favicon.ico" />
        <meta name="author" content="John Doe" />
      </Head>
        
        <main>
            <p>Your viewing version page</p>
            <p>
                tag: {data.version} <br />
                hash: {data.hash} <br />
                built: {data.built} <br />
            </p>
        </main>
    </div>
  )
}

// This gets called on every request
export async function getServerSideProps() {
  // Fetch data from external API
  const res = await fetch(`https://fotolink.app/api/`)
  const data = await res.json()

  // Pass data to the page via props
  return { props: { data } }
}