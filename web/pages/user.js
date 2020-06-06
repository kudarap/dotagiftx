import Head from 'next/head'

export default function User() {
  return (
    <div className="container">
      <Head>
        <title>User head</title>
        <link rel="icon" href="/favicon.ico" />
        <meta name="author" content="John Doe" />
      </Head>

        <main>
            <p>Your viewing user page</p>
        </main>
    </div>
  )
}