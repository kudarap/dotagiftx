import React, { useContext, useEffect, useState } from 'react'
import { useRouter } from 'next/router'
import { PayPalButtons, usePayPalScriptReducer } from '@paypal/react-paypal-js'
import { styled } from '@mui/material/styles'
import Typography from '@mui/material/Typography'
import Box from '@mui/material/Box'
import Header from '@/components/Header'
import Container from '@/components/Container'
import Footer from '@/components/Footer'
import AppContext from '@/components/AppContext'
import { Divider } from '@mui/material'
import Link from '@/components/Link'

const isPaypalLive = process.env.NEXT_PUBLIC_API_URL.startsWith('https://api.dotagiftx.com')

const subscriptions = {
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
    features: ['Partner Badge', 'Refresher Orb', "Shopkeeper's Contract", 'Dedicated Pos-5'],
    planId: 'P-2Y98477558961784RMJNLBYI',
    planIdLive: 'P-0EB00258NU2523843MJMW6JY',
  },
}

const ButtonWrapper = ({ type, planId, customId, onSuccess }) => {
  const [{ options }, dispatch] = usePayPalScriptReducer()

  useEffect(() => {
    dispatch({
      type: 'resetOptions',
      value: {
        ...options,
        intent: 'subscription',
      },
    })
  }, [type])

  return (
    <PayPalButtons
      createSubscription={(data, actions) => {
        return actions.subscription
          .create({
            plan_id: planId,
            custom_id: customId,
          })
          .then(orderId => {
            return orderId
          })
      }}
      onApprove={onSuccess}
      style={{
        label: 'subscribe',
        color: 'blue',
      }}
    />
  )
}

const FeatureList = styled('ul')(({ theme }) => ({
  listStyle: 'none',
  '& li:before': {
    content: `'✔'`,
    marginRight: 8,
  },
  paddingLeft: 0,
}))

const priceTable = {
  partner: 20,
  trader: 3,
  supporter: 1,
}

const minimumCycle = {
  partner: 6,
  trader: 12,
  supporter: 12,
}

const manualPriceOverhead = 0.6

export default function Subscription({ data }) {
  const { currentAuth } = useContext(AppContext)

  const router = useRouter()
  const { query } = router
  const [subscription, setSubscription] = useState()
  const [minimumSubscriptionCycle, setMinimumSubscriptionCycle] = useState()
  const [subscriptionPrice, setSubscriptionPrice] = useState()
  useEffect(() => {
    if (!query.id) {
      return
    }
    if (!currentAuth.user_id) {
      router.push('/login')
      return
    }

    setSubscription(subscriptions[query.id])
    setMinimumSubscriptionCycle(minimumCycle[query.id])
    setSubscriptionPrice(priceTable[query.id])
  }, [query.id, currentAuth.user_id])

  const handleSuccess = data => {
    console.log(data)
    // send orderId to subscription verifier to ack the process
    router.push(`/thanks-subscriber?sub=${subscription.id}&subid=${data.subscriptionID}`)
  }

  const isReady = currentAuth.steam_id && subscription

  return (
    <div className="container">
      <Header />

      <main>
        <Container>
          <Box sx={{ mt: 8, mb: 4, textAlign: 'center' }}>
            <Typography
              sx={{ mt: 8 }}
              style={{
                background: 'linear-gradient( to right, #CB8F37 20%, #F0CF59 50%, #B5793D 80% )',
                WebkitBackgroundClip: 'text',
                WebkitTextFillColor: 'transparent',
                letterSpacing: 3,
                // textTransform: 'uppercase',
              }}
              variant="h3"
              component="h1"
              fontWeight="bold"
              color="secondary">
              Dotagift Plus
            </Typography>
            {isReady && (
              <>
                <Typography variant="h6">{subscription.name} Subscription</Typography>
                <Typography color="textSecondary">${subscriptionPrice} monthly</Typography>
                <FeatureList>
                  {subscription.features.map(v => (
                    <li key={v}>{v}</li>
                  ))}
                </FeatureList>
              </>
            )}
          </Box>

          <Box sx={{ maxWidth: 500, m: '0 auto' }}>
            {isReady && (
              <ButtonWrapper
                type="subscription"
                planId={isPaypalLive ? subscription.planIdLive : subscription.planId}
                customId={`STEAMID-${currentAuth.steam_id}`}
                onSuccess={handleSuccess}
              />
            )}
          </Box>

          <Divider sx={{ mt: 5, mb: 5 }} />

          {isReady && (
            <Box sx={{ maxWidth: 600, m: 'auto' }} textAlign="center">
              <Typography variant="h6" sx={{ mb: 1 }}>
                Paypal is not supported in your country?
              </Typography>
              <Typography>
                You can pay one-time via{' '}
                <Link
                  color="secondary"
                  target="_blank"
                  rel="noreferrer noopener"
                  href="https://steamcommunity.com/market/listings/440/Mann%20Co.%20Supply%20Crate%20Key">
                  TF2 Keys
                </Link>
                {' '}or{' '}
                <Link
                  color="secondary"
                  target="_blank"
                  rel="noreferrer noopener"
                  href="https://steamcommunity.com/market/listings/570/Fractal%20Horns%20of%20Inner%20Abysm">
                  TB Arcanas
                </Link>
                {' '}
                with minimum of {minimumSubscriptionCycle} months and +{manualPriceOverhead * 100}% overhead for steam
                community market conversion and manual processing fees.
              </Typography>

              <Box>
                <Typography>
                  <br />${subscriptionPrice} x {minimumSubscriptionCycle} months = $
                  {subscriptionPrice * minimumSubscriptionCycle}
                  <br />
                  +{manualPriceOverhead * 100}% SCM overhead fee = $
                  {Math.round((subscriptionPrice * minimumSubscriptionCycle) * manualPriceOverhead)}
                  <br />
                  <strong>
                    Total = ${Math.round(subscriptionPrice * minimumSubscriptionCycle * (1 + manualPriceOverhead))}
                  </strong>
                </Typography>
              </Box>

              <Box textAlign="left">
                <ol>
                  <li>
                    Acquire your TF2 keys and/or TB Arcanas and total should equal to{' '}
                    <strong>
                      ${Math.round(subscriptionPrice * minimumSubscriptionCycle * (1 + manualPriceOverhead))}
                    </strong>
                    .
                  </li>
                  <li>
                    Send{' '}
                    <Link
                      color="secondary"
                      target="_blank"
                      rel="noreferrer noopener"
                      href="https://steamcommunity.com/tradeoffer/new/?partner=128321450&token=38BJlyuW">
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
                      href="https://discord.gg/3JVU2EumRw">
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
  )
}
