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

const isPaypalLive = process.env.NEXT_PUBLIC_API_URL.startsWith('https://dotagiftx.com')

const subscriptions = {
  supporter: {
    id: 'supporter',
    name: 'Supporter',
    features: ['Supporter Badge', 'Refresher Shard'],
    planId: 'P-616716383W896284VMJMR4CY',
    planIdLive: 'P-0EB00258NU2523843MJMW6JY',
  },
  trader: {
    id: 'trader',
    name: 'Trader',
    features: ['Trader Badge', 'Refresher Orb'],
    planId: 'P-28V29656NC814125PMJMSDWQ',
    planIdLive: 'P-6TG171216S461482EMJMW55Q',
  },
  partner: {
    id: 'partner',
    name: 'Partner',
    features: ['Partner Badge', 'Refresher Orb', "Shopkeeper's Contract", 'Dedicated Pos-5'],
    planId: 'P-2FS77965H7642004PMJMSD6Q',
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

export default function Subscription({ data }) {
  const { currentAuth } = useContext(AppContext)

  const router = useRouter()
  const { query } = router
  const [subscription, setSubscription] = useState(null)
  useEffect(() => {
    if (!query.id) {
      return
    }
    if (!currentAuth.user_id) {
      router.push('/login')
      return
    }

    setSubscription(subscriptions[query.id])
  }, [query.id, currentAuth.user_id])

  const handleSuccess = data => {
    // send orderId to subscription verifier to ack the process
    router.push(`/thanks-subscriber?id=${subscription.id}`)
  }

  const isReady = currentAuth.steam_id && subscription

  return (
    <div className="container">
      <Header />

      <main>
        <Container>
          <Box sx={{ mt: 8, mb: 4, textAlign: 'center' }}>
            <Typography variant="h4" component="h1" fontWeight="bold">
              Dotagift Plus
            </Typography>
            {isReady && (
              <>
                <Typography variant="h6">{subscription.name} Subscription</Typography>
                <FeatureList>
                  {subscription.features.map(v => (
                    <li key={v}>{v}</li>
                  ))}
                </FeatureList>
              </>
            )}
          </Box>

          <Box sx={{ textAlign: 'center' }}>
            {isReady && (
              <ButtonWrapper
                type="subscription"
                planId={isPaypalLive ? subscription.planIdLive : subscription.planId}
                customId={`STEAMID-${currentAuth.steam_id}`}
                onSuccess={handleSuccess}
              />
            )}
          </Box>
        </Container>
      </main>

      <Footer />
    </div>
  )
}
