import Head from 'next/head'
import { useRouter } from 'next/router'
import Banner from "@/components/Banner";
import BuildInfo from "@/components/BuildInfo";

export default function Id() {
  const router = useRouter()
  const { id } = router.query

  return (
    <div className="container">
      <Head>
        <title>User {id}</title>
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