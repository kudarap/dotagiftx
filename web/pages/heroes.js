import React, { useState } from 'react'
import Head from 'next/head'
import Paper from '@mui/material/Paper'
import Typography from '@mui/material/Typography'
import Header from '@/components/Header'
import Container from '@/components/Container'
import Footer from '@/components/Footer'
import { Grid } from '@mui/material'
import { styled } from '@mui/system'
import Image from 'next/image'
import Link from '@/components/Link'
import SearchInput from '@/components/SearchInput'
import { APP_NAME } from '@/constants/strings'

const allHeroes = [
  {
    name: 'Abaddon',
    image: 'abaddon.png',
    attribute: 'universal',
  },
  {
    name: 'Alchemist',
    image: 'alchemist.png',
    attribute: 'strength',
  },
  {
    name: 'Ancient Apparition',
    image: 'ancient_apparition.png',
    attribute: 'intelligence',
  },
  {
    name: 'Anti-Mage',
    image: 'antimage.png',
    attribute: 'agility',
  },
  {
    name: 'Arc Warden',
    image: 'arc_warden.png',
    attribute: 'universal',
  },
  {
    name: 'Axe',
    image: 'axe.png',
    attribute: 'strength',
  },
  {
    name: 'Bane',
    image: 'bane.png',
    attribute: 'universal',
  },
  {
    name: 'Batrider',
    image: 'batrider.png',
    attribute: 'universal',
  },
  {
    name: 'Beastmaster',
    image: 'beastmaster.png',
    attribute: 'universal',
  },
  {
    name: 'Bloodseeker',
    image: 'bloodseeker.png',
    attribute: 'agility',
  },
  {
    name: 'Bounty Hunter',
    image: 'bounty_hunter.png',
    attribute: 'agility',
  },
  {
    name: 'Brewmaster',
    image: 'brewmaster.png',
    attribute: 'universal',
  },
  {
    name: 'Bristleback',
    image: 'bristleback.png',
    attribute: 'strength',
  },
  {
    name: 'Broodmother',
    image: 'broodmother.png',
    attribute: 'agility',
  },
  {
    name: 'Centaur Warrunner',
    image: 'centaur.png',
    attribute: 'strength',
  },
  {
    name: 'Chaos Knight',
    image: 'chaos_knight.png',
    attribute: 'strength',
  },
  {
    name: 'Chen',
    image: 'chen.png',
    attribute: 'intelligence',
  },
  {
    name: 'Clinkz',
    image: 'clinkz.png',
    attribute: 'agility',
  },
  {
    name: 'Clockwerk',
    image: 'rattletrap.png',
    attribute: 'strength',
  },
  {
    name: 'Crystal Maiden',
    image: 'crystal_maiden.png',
    attribute: 'intelligence',
  },
  {
    name: 'Dark Seer',
    image: 'dark_seer.png',
    attribute: 'intelligence',
  },
  {
    name: 'Dark Willow',
    image: 'dark_willow.png',
    attribute: 'intelligence',
  },
  {
    name: 'Dawnbreaker',
    image: 'dawnbreaker.png',
    attribute: 'strength',
  },
  {
    name: 'Dazzle',
    image: 'dazzle.png',
    attribute: 'universal',
  },
  {
    name: 'Death Prophet',
    image: 'death_prophet.png',
    attribute: 'universal',
  },
  {
    name: 'Disruptor',
    image: 'disruptor.png',
    attribute: 'intelligence',
  },
  {
    name: 'Doom',
    image: 'doom_bringer.png',
    attribute: 'strength',
  },
  {
    name: 'Dragon Knight',
    image: 'dragon_knight.png',
    attribute: 'strength',
  },
  {
    name: 'Drow Ranger',
    image: 'drow_ranger.png',
    attribute: 'agility',
  },
  {
    name: 'Earth Spirit',
    image: 'earth_spirit.png',
    attribute: 'strength',
  },
  {
    name: 'Earthshaker',
    image: 'earthshaker.png',
    attribute: 'strength',
  },
  {
    name: 'Elder Titan',
    image: 'elder_titan.png',
    attribute: 'strength',
  },
  {
    name: 'Ember Spirit',
    image: 'ember_spirit.png',
    attribute: 'agility',
  },
  {
    name: 'Enchantress',
    image: 'enchantress.png',
    attribute: 'intelligence',
  },
  {
    name: 'Enigma',
    image: 'enigma.png',
    attribute: 'universal',
  },
  {
    name: 'Faceless Void',
    image: 'faceless_void.png',
    attribute: 'agility',
  },
  {
    name: 'Grimstroke',
    image: 'grimstroke.png',
    attribute: 'intelligence',
  },
  {
    name: 'Gyrocopter',
    image: 'gyrocopter.png',
    attribute: 'agility',
  },
  {
    name: 'Hoodwink',
    image: 'hoodwink.png',
    attribute: 'agility',
  },
  {
    name: 'Huskar',
    image: 'huskar.png',
    attribute: 'strength',
  },
  {
    name: 'Invoker',
    image: 'invoker.png',
    attribute: 'intelligence',
  },
  {
    name: 'Io',
    image: 'wisp.png',
    attribute: 'universal',
  },
  {
    name: 'Jakiro',
    image: 'jakiro.png',
    attribute: 'intelligence',
  },
  {
    name: 'Juggernaut',
    image: 'juggernaut.png',
    attribute: 'agility',
  },
  {
    name: 'Keeper of the Light',
    image: 'keeper_of_the_light.png',
    attribute: 'intelligence',
  },
  {
    name: 'Kez',
    image: 'kez.png',
    attribute: 'agility',
  },
  {
    name: 'Kunkka',
    image: 'kunkka.png',
    attribute: 'strength',
  },
  {
    name: 'Legion Commander',
    image: 'legion_commander.png',
    attribute: 'strength',
  },
  {
    name: 'Leshrac',
    image: 'leshrac.png',
    attribute: 'intelligence',
  },
  {
    name: 'Lich',
    image: 'lich.png',
    attribute: 'intelligence',
  },
  {
    name: 'Lifestealer',
    image: 'life_stealer.png',
    attribute: 'strength',
  },
  {
    name: 'Lina',
    image: 'lina.png',
    attribute: 'intelligence',
  },
  {
    name: 'Lion',
    image: 'lion.png',
    attribute: 'intelligence',
  },
  {
    name: 'Lone Druid',
    image: 'lone_druid.png',
    attribute: 'agility',
  },
  {
    name: 'Luna',
    image: 'luna.png',
    attribute: 'agility',
  },
  {
    name: 'Lycan',
    image: 'lycan.png',
    attribute: 'strength',
  },
  {
    name: 'Magnus',
    image: 'magnataur.png',
    attribute: 'universal',
  },
  {
    name: 'Marci',
    image: 'marci.png',
    attribute: 'universal',
  },
  {
    name: 'Mars',
    image: 'mars.png',
    attribute: 'strength',
  },
  {
    name: 'Medusa',
    image: 'medusa.png',
    attribute: 'agility',
  },
  {
    name: 'Meepo',
    image: 'meepo.png',
    attribute: 'agility',
  },
  {
    name: 'Mirana',
    image: 'mirana.png',
    attribute: 'agility',
  },
  {
    name: 'Monkey King',
    image: 'monkey_king.png',
    attribute: 'agility',
  },
  {
    name: 'Morphling',
    image: 'morphling.png',
    attribute: 'agility',
  },
  {
    name: 'Muerta',
    image: 'muerta.png',
    attribute: 'intelligence',
  },
  {
    name: 'Naga Siren',
    image: 'naga_siren.png',
    attribute: 'agility',
  },
  {
    name: "Nature's Prophet",
    image: 'furion.png',
    attribute: 'universal',
  },
  {
    name: 'Necrophos',
    image: 'necrolyte.png',
    attribute: 'intelligence',
  },
  {
    name: 'Night Stalker',
    image: 'night_stalker.png',
    attribute: 'strength',
  },
  {
    name: 'Nyx Assassin',
    image: 'nyx_assassin.png',
    attribute: 'universal',
  },
  {
    name: 'Ogre Magi',
    image: 'ogre_magi.png',
    attribute: 'strength',
  },
  {
    name: 'Omniknight',
    image: 'omniknight.png',
    attribute: 'strength',
  },
  {
    name: 'Oracle',
    image: 'oracle.png',
    attribute: 'intelligence',
  },
  {
    name: 'Outworld Destroyer',
    image: 'obsidian_destroyer.png',
    attribute: 'intelligence',
  },
  {
    name: 'Pangolier',
    image: 'pangolier.png',
    attribute: 'universal',
  },
  {
    name: 'Phantom Assassin',
    image: 'phantom_assassin.png',
    attribute: 'agility',
  },
  {
    name: 'Phantom Lancer',
    image: 'phantom_lancer.png',
    attribute: 'agility',
  },
  {
    name: 'Phoenix',
    image: 'phoenix.png',
    attribute: 'strength',
  },
  {
    name: 'Primal Beast',
    image: 'primal_beast.png',
    attribute: 'strength',
  },
  {
    name: 'Puck',
    image: 'puck.png',
    attribute: 'intelligence',
  },
  {
    name: 'Pudge',
    image: 'pudge.png',
    attribute: 'strength',
  },
  {
    name: 'Pugna',
    image: 'pugna.png',
    attribute: 'intelligence',
  },
  {
    name: 'Queen of Pain',
    image: 'queenofpain.png',
    attribute: 'intelligence',
  },
  {
    name: 'Razor',
    image: 'razor.png',
    attribute: 'agility',
  },
  {
    name: 'Riki',
    image: 'riki.png',
    attribute: 'agility',
  },
  {
    name: 'Ringmaster',
    image: 'ringmaster.png',
    attribute: 'intelligence',
  },
  {
    name: 'Rubick',
    image: 'rubick.png',
    attribute: 'intelligence',
  },
  {
    name: 'Sand King',
    image: 'sand_king.png',
    attribute: 'universal',
  },
  {
    name: 'Shadow Demon',
    image: 'shadow_demon.png',
    attribute: 'intelligence',
  },
  {
    name: 'Shadow Fiend',
    image: 'nevermore.png',
    attribute: 'agility',
  },
  {
    name: 'Shadow Shaman',
    image: 'shadow_shaman.png',
    attribute: 'intelligence',
  },
  {
    name: 'Silencer',
    image: 'silencer.png',
    attribute: 'intelligence',
  },
  {
    name: 'Skywrath Mage',
    image: 'skywrath_mage.png',
    attribute: 'intelligence',
  },
  {
    name: 'Slardar',
    image: 'slardar.png',
    attribute: 'strength',
  },
  {
    name: 'Slark',
    image: 'slark.png',
    attribute: 'agility',
  },
  {
    name: 'Snapfire',
    image: 'snapfire.png',
    attribute: 'universal',
  },
  {
    name: 'Sniper',
    image: 'sniper.png',
    attribute: 'agility',
  },
  {
    name: 'Spectre',
    image: 'spectre.png',
    attribute: 'universal',
  },
  {
    name: 'Spirit Breaker',
    image: 'spirit_breaker.png',
    attribute: 'strength',
  },
  {
    name: 'Storm Spirit',
    image: 'storm_spirit.png',
    attribute: 'intelligence',
  },
  {
    name: 'Sven',
    image: 'sven.png',
    attribute: 'strength',
  },
  {
    name: 'Techies',
    image: 'techies.png',
    attribute: 'universal',
  },
  {
    name: 'Templar Assassin',
    image: 'templar_assassin.png',
    attribute: 'agility',
  },
  {
    name: 'Terrorblade',
    image: 'terrorblade.png',
    attribute: 'agility',
  },
  {
    name: 'Tidehunter',
    image: 'tidehunter.png',
    attribute: 'strength',
  },
  {
    name: 'Timbersaw',
    image: 'shredder.png',
    attribute: 'strength',
  },
  {
    name: 'Tinker',
    image: 'tinker.png',
    attribute: 'intelligence',
  },
  {
    name: 'Tiny',
    image: 'tiny.png',
    attribute: 'strength',
  },
  {
    name: 'Treant Protector',
    image: 'treant.png',
    attribute: 'strength',
  },
  {
    name: 'Troll Warlord',
    image: 'troll_warlord.png',
    attribute: 'agility',
  },
  {
    name: 'Tusk',
    image: 'tusk.png',
    attribute: 'strength',
  },
  {
    name: 'Underlord',
    image: 'abyssal_underlord.png',
    attribute: 'strength',
  },
  {
    name: 'Undying',
    image: 'undying.png',
    attribute: 'strength',
  },
  {
    name: 'Ursa',
    image: 'ursa.png',
    attribute: 'agility',
  },
  {
    name: 'Vengeful Spirit',
    image: 'vengefulspirit.png',
    attribute: 'agility',
  },
  {
    name: 'Venomancer',
    image: 'venomancer.png',
    attribute: 'universal',
  },
  {
    name: 'Viper',
    image: 'viper.png',
    attribute: 'agility',
  },
  {
    name: 'Visage',
    image: 'visage.png',
    attribute: 'universal',
  },
  {
    name: 'Void Spirit',
    image: 'void_spirit.png',
    attribute: 'universal',
  },
  {
    name: 'Warlock',
    image: 'warlock.png',
    attribute: 'intelligence',
  },
  {
    name: 'Weaver',
    image: 'weaver.png',
    attribute: 'agility',
  },
  {
    name: 'Windranger',
    image: 'windrunner.png',
    attribute: 'universal',
  },
  {
    name: 'Winter Wyvern',
    image: 'winter_wyvern.png',
    attribute: 'intelligence',
  },
  {
    name: 'Witch Doctor',
    image: 'witch_doctor.png',
    attribute: 'intelligence',
  },
  {
    name: 'Wraith King',
    image: 'skeleton_king.png',
    attribute: 'strength',
  },
  {
    name: 'Zeus',
    image: 'zuus.png',
    attribute: 'intelligence',
  },
]

const Item = styled(Paper)(({ theme }) => ({
  backgroundColor: '#1A20278C',
  ...theme.typography.body,
  padding: theme.spacing(1),
  paddingTop: theme.spacing(1),
  textAlign: 'center',
  color: theme.palette.text.primary,
}))

export default function Treasures() {
  const [heroes, setHeroes] = useState(allHeroes)
  const [searchTerm, setSearchTerm] = useState()
  const handleChange = term => {
    setSearchTerm(term)
    setHeroes(allHeroes.filter(v => !!v.name.match(new RegExp(term, 'gi'))))
  }

  return (
    <div className="container">
      <Head>
        <meta charSet="UTF-8" />
        <title>{APP_NAME} :: Heroes</title>
        <link rel="icon" href="/favicon.ico" />
      </Head>

      <Header />

      <main>
        <div
          style={{
            width: '100%',
            height: 500,
            maskImage: 'linear-gradient(to top, transparent 25%, black 100%)',
            WebkitMaskImage: 'linear-gradient(to top, transparent 0%, black 90%)',
            position: 'relative',
            zIndex: 0,
          }}>
          <div
            style={{
              background:
                'url(https://cdn.cloudflare.steamstatic.com/steam/apps/570/library_hero.jpg?t=1724395576617) no-repeat center center',
              backgroundColor: '#2a2638ff',
              backgroundSize: 'cover',
              width: '100%',
              height: '100%',
            }}></div>
        </div>

        <Container style={{ position: 'relative' }}>
          <Typography
            sx={{ mt: -55, mb: 4 }}
            variant="h3"
            component="h1"
            fontWeight="bold"
            color="pimary">
            All Heroes
          </Typography>

          <SearchInput
            value={searchTerm}
            onChange={handleChange}
            placeholder="Search heroes names..."
            label=""
          />

          <Grid container spacing={1} sx={{ mt: 2 }}>
            {heroes.map(hero => {
              return (
                <Grid item xs={4} md={2} key={hero.name}>
                  <Link href={`/search?hero=${hero.name}`} underline="none">
                    <Item>
                      <div>
                        <Image
                          src={'/assets/heroes/' + hero.image}
                          alt={hero.name}
                          width={256 * 0.7}
                          height={144 * 0.7}
                        />
                      </div>
                      <Typography noWrap>{hero.name}</Typography>
                    </Item>
                  </Link>
                </Grid>
              )
            })}
          </Grid>
        </Container>
      </main>

      <Footer />
    </div>
  )
}
