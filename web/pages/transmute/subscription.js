import React, { useEffect, useState } from 'react'
import { PayPalButtons, usePayPalScriptReducer } from '@paypal/react-paypal-js'
import { styled } from '@mui/material/styles'
import Typography from '@mui/material/Typography'
import Box from '@mui/material/Box'
import Header from '@/components/Header'
import Container from '@/components/Container'
import Footer from '@/components/Footer'
import { useRouter } from 'next/router'

const isPaypalLive = false

const subscriptions = {
  supporter: {
    name: 'Supporter',
    features: ['Supporter Badge', 'Refresher Shard'],
    planId: 'P-616716383W896284VMJMR4CY',
    planIdLive: 'P-616716383W896284VMJMR4CY',
  },
  trader: {
    name: 'Trader',
    features: ['Trader Badge', 'Refresher Orb'],
    planId: 'P-28V29656NC814125PMJMSDWQ',
    planIdLive: 'P-28V29656NC814125PMJMSDWQ',
  },
  partner: {
    name: 'Partner',
    features: ['Partner Badge', 'Refresher Orb', "Shopkeeper's Contract", 'Dedicated Pos-5'],
    planId: 'P-2FS77965H7642004PMJMSD6Q',
    planIdLive: 'P-2FS77965H7642004PMJMSD6Q',
  },
}

const ButtonWrapper = ({ type, planId, customId }) => {
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
            // send orderId to subscription verifier to ack the process
            return orderId
          })
      }}
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
  const router = useRouter()
  const { query } = router
  const [subscription, setSubscription] = useState(null)
  useEffect(() => {
    if (!query.id) {
      return
    }

    setSubscription(subscriptions[query.id])
  }, [query.id])

  return (
    <div className="container">
      <Header />

      <main>
        <Container>
          <Box sx={{ mt: 8, mb: 4, textAlign: 'center' }}>
            <Typography variant="h4" component="h1" fontWeight="bold">
              Dotagift Plus
            </Typography>
            {subscription && (
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
            {subscription && (
              <ButtonWrapper
                type="subscription"
                planId={isPaypalLive ? subscription.planIdLive : subscription.planId}
                customId="STEAMID-2000001"
              />
            )}
          </Box>
        </Container>
      </main>

      <Footer />
    </div>
  )
}
