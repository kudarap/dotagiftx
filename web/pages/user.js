import Head from 'next/head'
import Banner from "../components/Banner";
import BuildInfo from "../components/BuildInfo";

export default function User() {
  return (
    <div className="container">
      <Head>
        <title>User head</title>
        <link rel="icon" href="/favicon.ico" />
        <meta name="author" content="John Doe" />
      </Head>

        <main>
            <p>Your viewing user pagex</p>

            <Banner />

            <BuildInfo />
        </main>
    </div>
  )
}