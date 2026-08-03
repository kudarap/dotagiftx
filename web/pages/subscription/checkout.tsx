import React, { useContext, useEffect } from 'react'
import { useRouter } from 'next/router'
import { PayPalProvider, PayPalSubscriptionButton } from '@paypal/react-paypal-js/sdk-v6'
import { styled } from '@mui/material/styles'
import Typography from '@mui/material/Typography'
import Box from '@mui/material/Box'
import { Alert, Divider } from '@mui/material'
import Header from '@/components/Header'
import Container from '@/components/Container'
import Footer from '@/components/Footer'
import AppContext from '@/components/AppContext'
import Link from '@/components/Link'
import { createMySubscription, myProfile } from '@/service/api'
import { USER_SUBSCRIPTION_MAP_LABEL } from '@/constants/user'
import type { Profile } from '@/lib/types'

const PAYPAL_CLIENT_ID = process.env.NEXT_PUBLIC_PAYPAL_CLIENT_ID

const isPaypalLive = (() => {
  try {
    const url = new URL(process.env.NEXT_PUBLIC_API_URL as string)
    return url.protocol === 'https:' && url.hostname === 'api.dotagiftx.com'
  } catch {
    return false
  }
})()

interface Subscription {
  id: string
  name: string
  features: string[]
  planId: string
  planIdLive: string
}

const subscriptions: Record<string, Subscription> = {
  supporter: {
    id: 'supporter',
    name: 'Supporter',
    features: ['Supporter Badge', 'Refresher Shard'],
    planId: 'P-16467111M44423113MJNKYKI',
    planIdLive: 'P-8JJ23834W3257961PMJMEB5A',
  },
  trader: {
    id: 'trader',
    name: 'Trader',
    features: ['Trader Badge', 'Refresher Orb'],
    planId: 'P-38P22523A72261937MJNLBRI',
    planIdLive: 'P-6TG171216S461482EMJMW55Q',
  },
  partner: {
    id: 'partner',
    name: 'Partner',
    features: ['Partner Badge', 'Refresher Orb', "Shopkeeper's Contract"],
    planId: 'P-2Y98477558961784RMJNLBYI',
    planIdLive: 'P-0EB00258NU2523843MJMW6JY',
  },
}

function ButtonWrapper({
  planId,
  onSuccess,
}: {
  planId: string
  onSuccess: (res: { subscriptionId: string }) => void
}) {
  return (
    <PayPalSubscriptionButton
      createSubscription={async () => {
        const res = (await createMySubscription(planId)) as { id: string }
        return { subscriptionId: res.id }
      }}
      onApprove={async res => {
        onSuccess(res)
      }}
      onCancel={data => console.log('Cancelled:', data)}
      onError={error => console.error('Error:', error)}
      presentationMode="auto"
    />
  )
}

const FeatureList = styled('ul')(() => ({
  listStyle: 'none',
  '& li:before': {
    content: `'✔'`,
    marginRight: 8,
  },
  paddingLeft: 0,
}))

const manualPriceOverhead = 0.6

const priceTable: Record<string, number> = {
  partner: 20,
  trader: 3,
  supporter: 1,
}

const minimumCycle: Record<string, number> = {
  partner: 6,
  trader: 12,
  supporter: 12,
}

const subscriptionStyleCard: Record<string, { borderTop: string; backgroundImage: string }> = {
  partner: {
    borderTop: '2px solid #ae7f1e',
    backgroundImage: 'linear-gradient(#a6791d63, #a6791d)',
  },
  trader: {
    borderTop: '2px solid #629cbd',
    backgroundImage: 'linear-gradient(#578ba863, #578ba8)',
  },
  supporter: {
    borderTop: '2px solid #596b95',
    backgroundImage: 'linear-gradient(#4654755c, #465475)',
  },
}

interface SubscriptionProfile {
  subscription?: number | null
  subscription_type?: string
  subscription_ends_at?: string | null
  subscriptionLabel?: string
}

const defaultProfile: SubscriptionProfile = {
  subscription: null,
  subscription_type: '',
  subscription_ends_at: null,
  subscriptionLabel: '',
}

export default function Checkout() {
  const { currentAuth } = useContext(AppContext)

  const router = useRouter()
  const { query } = router
  const subscriptionId = typeof query.id === 'string' ? query.id : undefined
  const subscription = subscriptionId ? subscriptions[subscriptionId] : undefined
  const minimumSubscriptionCycle = subscriptionId ? minimumCycle[subscriptionId] : undefined
  const subscriptionPrice = subscriptionId ? priceTable[subscriptionId] : undefined
  const cardStyle = subscriptionId ? subscriptionStyleCard[subscriptionId] : undefined
  const [currentProfile, setCurrentProfile] = React.useState<SubscriptionProfile>(defaultProfile)
  useEffect(() => {
    if (!query.id) {
      return
    }

    if (!currentAuth?.user_id) {
      router.push('/login')
      return
    }

    ;(async () => {
      const res = (await myProfile.GET(true)) as Profile
      setCurrentProfile({
        ...res,
        subscriptionLabel: USER_SUBSCRIPTION_MAP_LABEL[res.subscription as number],
      })
    })()
  }, [router, query.id, currentAuth?.user_id])

  const handleSuccess = (res: { subscriptionId: string }) => {
    // send orderId to subscription verifier to ack the process
    router.push(`/subscription/thank-you?subid=${res.subscriptionId}`)
  }

  const isReady = currentAuth?.steam_id && subscription && subscriptionPrice

  return (
    <PayPalProvider
      clientId={PAYPAL_CLIENT_ID}
      environment={isPaypalLive ? 'production' : 'sandbox'}
      components={['paypal-subscriptions']}
      pageType="checkout"
    >
      <div className="container">
        <Header />

        <main>
          <Container>
            <Box sx={{ mt: 8, mb: 4, textAlign: 'center' }}>
              {!!currentProfile.subscription && subscription && (
                <Alert severity="warning" sx={{ mt: 2 }}>
                  You have an active subscription. If you wish to change from{' '}
                  <strong>
                    <u>{currentProfile.subscriptionLabel}</u>
                  </strong>{' '}
                  to{' '}
                  <strong>
                    <u>{subscription.name}</u>
                  </strong>{' '}
                  {currentProfile.subscription_type == 'paypal' && (
                    <span>
                      you must cancel the subscription from your PayPal&apos;s dashboard before
                      proceeding on this page
                    </span>
                  )}
                  {currentProfile.subscription_type == 'manual' && (
                    <span>please be noted that it will be overriden and not refundable</span>
                  )}
                  . Any remaining days from the previous subscription will not be carried over.
                </Alert>
              )}

              <Typography
                style={{
                  background: 'linear-gradient( to right, #CB8F37 20%, #F0CF59 50%, #B5793D 80% )',
                  WebkitBackgroundClip: 'text',
                  WebkitTextFillColor: 'transparent',
                  letterSpacing: 3,
                }}
                variant="h3"
                component="h1"
                color="secondary"
                sx={{
                  fontWeight: 'bold',
                  mt: 8,
                }}
              >
                Dotagift Plus
              </Typography>

              {isReady && subscription && cardStyle && subscriptionPrice && (
                <>
                  <Box
                    sx={{
                      py: 4,
                      px: 3,
                      borderTop: cardStyle.borderTop,
                      backgroundImage: cardStyle.backgroundImage,
                      borderRadius: 2,
                      width: 500,
                      m: '30px auto',
                    }}
                  >
                    <Typography
                      variant="h5"
                      sx={{
                        fontWeight: 'bold',
                      }}
                    >
                      {subscription.name} Subscription
                    </Typography>

                    <Box
                      component={FeatureList}
                      sx={{ height: 160, textAlign: 'justify', display: 'inline-table' }}
                    >
                      {subscription.features.map(v => (
                        <li key={v}>{v}</li>
                      ))}
                    </Box>

                    <Typography variant="h6" sx={{ mb: 1 }}>
                      US ${subscriptionPrice}
                      &nbsp;
                      <Typography variant="subtitle2" component="span" color="textSecondary">
                        /mo
                      </Typography>
                    </Typography>
                    <Box sx={{ maxWidth: 225, m: '0 auto 10px' }}>
                      <ButtonWrapper
                        planId={isPaypalLive ? subscription.planIdLive : subscription.planId}
                        onSuccess={handleSuccess}
                      />
                    </Box>
                  </Box>
                </>
              )}
            </Box>

            <Divider sx={{ mt: 5, mb: 5 }} />

            {isReady && subscriptionPrice && minimumSubscriptionCycle && (
              <Box
                sx={{
                  textAlign: 'center',
                  maxWidth: 600,
                  m: 'auto',
                }}
              >
                <Typography variant="h6" sx={{ mb: 1 }}>
                  PayPal not supported?
                </Typography>
                <Typography>
                  You can pay one-time via{' '}
                  <Link
                    color="secondary"
                    target="_blank"
                    rel="noreferrer noopener"
                    href="https://steamcommunity.com/market/listings/440/Mann%20Co.%20Supply%20Crate%20Key"
                  >
                    TF2 Keys
                  </Link>{' '}
                  or{' '}
                  <Link
                    color="secondary"
                    target="_blank"
                    rel="noreferrer noopener"
                    href="https://steamcommunity.com/market/listings/570/Fractal%20Horns%20of%20Inner%20Abysm"
                  >
                    TB Arcanas
                  </Link>{' '}
                  with minimum of {minimumSubscriptionCycle} months and +{manualPriceOverhead * 100}
                  % overhead for steam community market conversion and manual processing fees.
                </Typography>

                <Box>
                  <Typography>
                    <br />${subscriptionPrice} x {minimumSubscriptionCycle} months = $
                    {subscriptionPrice * minimumSubscriptionCycle}
                    <br />+{manualPriceOverhead * 100}% SCM overhead fee = $
                    {Math.round(subscriptionPrice * minimumSubscriptionCycle * manualPriceOverhead)}
                    <br />
                    <strong>
                      Total = $
                      {Math.round(
                        subscriptionPrice * minimumSubscriptionCycle * (1 + manualPriceOverhead)
                      )}
                    </strong>
                  </Typography>
                </Box>

                <Box
                  sx={{
                    textAlign: 'left',
                  }}
                >
                  <ol>
                    <li>
                      Acquire your TF2 keys and/or TB Arcanas and total should equal to{' '}
                      <strong>
                        $
                        {Math.round(
                          subscriptionPrice * minimumSubscriptionCycle * (1 + manualPriceOverhead)
                        )}
                      </strong>
                      .
                    </li>
                    <li>
                      Send{' '}
                      <Link
                        color="secondary"
                        target="_blank"
                        rel="noreferrer noopener"
                        href="https://steamcommunity.com/tradeoffer/new/?partner=128321450&token=38BJlyuW"
                      >
                        trade offer
                      </Link>{' '}
                      and indicate your subscription plan.
                    </li>
                    <li>
                      Notify us on our{' '}
                      <Link
                        color="secondary"
                        target="_blank"
                        rel="noreferrer noopener"
                        href="https://discord.gg/3JVU2EumRw"
                      >
                        Discord
                      </Link>{' '}
                      that you made a trade offer and allow us to process in 2-3 days.
                    </li>
                  </ol>
                </Box>
              </Box>
            )}
          </Container>
        </main>

        <Footer />
      </div>
    </PayPalProvider>
  )
}
