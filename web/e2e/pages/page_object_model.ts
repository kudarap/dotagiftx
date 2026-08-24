import { Page } from '@playwright/test'
import { BansPage } from './bans_page'
import { DiscordPage } from './discord_page'
import { IndexPage } from './index_page'
import { PlusPage } from './plus_page'
import { FaqsPage } from './faqs_page'
import { GuidesPage } from './guides_page'
import { HeroesPage } from './heroes_pages'
import { MiddleManPage } from './middle_man_page'
import { MobilePage } from './mobile_page'
import { ModeratorsPage } from './moderators_page'
import { RulesPage } from './rules_page'
import { TreasuresPage } from './treasures_page'
import { UpdatesPage } from './updates_page'

export class PageObjectModel {
  private readonly page: Page
  public bansPage: BansPage
  public discordPage: DiscordPage
  public indexPage: IndexPage
  public plusPage: PlusPage
  public faqsPage: FaqsPage
  public guidesPage: GuidesPage
  public heroesPage: HeroesPage
  public middleManPage: MiddleManPage
  public mobilePage: MobilePage
  public moderatorsPage: ModeratorsPage
  public rulesPage: RulesPage
  public treasuresPage: TreasuresPage
  public updatesPage: UpdatesPage

  constructor(page: Page) {
    this.page = page
    this.bansPage = new BansPage(page)
    this.discordPage = new DiscordPage(page)
    this.indexPage = new IndexPage(page)
    this.plusPage = new PlusPage(page)
    this.faqsPage = new FaqsPage(page)
    this.guidesPage = new GuidesPage(page)
    this.heroesPage = new HeroesPage(page)
    this.middleManPage = new MiddleManPage(page)
    this.mobilePage = new MobilePage(page)
    this.moderatorsPage = new ModeratorsPage(page)
    this.rulesPage = new RulesPage(page)
    this.treasuresPage = new TreasuresPage(page)
    this.updatesPage = new UpdatesPage(page)
  }
}
