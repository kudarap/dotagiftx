import { APP_NAME } from '../../constants/strings'

export const PageTitles = {
  MAIN_PAGE: `Dota 2 Giftables Community Market :: ${APP_NAME}`,
  TREASURES_PAGE: `All Treasures :: ${APP_NAME}`,
  HEROES_PAGE: `Heroes :: ${APP_NAME}`,
  DOTAGIFTX_PLUS_PAGE: `Plus :: ${APP_NAME}`,
  MOBILE_APP_PAGE: `${APP_NAME} for Mobile :: ${APP_NAME}`,
  BANS_PAGE: `Banned users :: ${APP_NAME}`,
  RULES_PAGE: `Rules :: ${APP_NAME}`,
  GUIDES_PAGE: `Guides :: ${APP_NAME}`,
  FQA_PAGE: `Frequently Asked Questions :: ${APP_NAME}`,
  MIDDLEMAN_PAGE: `Middleman :: ${APP_NAME}`,
  MODERATORS_PAGE: `Moderators :: ${APP_NAME}`,
  UPDATES_PAGE: `Updates :: ${APP_NAME}`,
} as const

export type PageTitles = (typeof PageTitles)[keyof typeof PageTitles]
