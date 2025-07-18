package vkma

import (
	"strconv"
)

// Referral source.
type Referral string

// Possible values.
const (
	Catalog               Referral = "catalog"
	CatalogRecent         Referral = "catalog_recent"
	CatalogFavourites     Referral = "catalog_favourites"
	CatalogRecommendation Referral = "catalog_recommendation"
	CatalogTopDau         Referral = "catalog_top_dau"
	CatalogEntertainment  Referral = "catalog_entertainment"
	CatalogCommunication  Referral = "catalog_communication"
	CatalogTools          Referral = "catalog_tools"
	CatalogShopping       Referral = "catalog_shopping"
	CatalogEvents         Referral = "catalog_events"
	CatalogEducation      Referral = "catalog_education"
	CatalogPayments       Referral = "catalog_payments"
	CatalogFinance        Referral = "catalog_finance"
	CatalogFood           Referral = "catalog_food"
	CatalogHealth         Referral = "catalog_health"
	CatalogTravel         Referral = "catalog_travel"
	CatalogTaxi           Referral = "catalog_taxi"
	CatalogJobs           Referral = "catalog_jobs"
	CatalogRealty         Referral = "catalog_realty"
	CatalogBusiness       Referral = "catalog_business"
	CatalogLifestyle      Referral = "catalog_lifestyle"
	CatalogAdmin          Referral = "catalog_admin"
	BoardTopicAll         Referral = "board_topic_all"
	BoardTopicView        Referral = "board_topic_view"
	Feed                  Referral = "feed"
	FeedPost              Referral = "feed_post"
	FeedComments          Referral = "feed_comments"
	FeaturingDiscover     Referral = "featuring_discover"
	FeaturingMenu         Referral = "featuring_menu"
	FeaturingNew          Referral = "featuring_new"
	Fave                  Referral = "fave"
	FaveLinks             Referral = "fave_links"
	FavePosts             Referral = "fave_posts"
	Group                 Referral = "group"
	GroupMenu             Referral = "group_menu"
	GroupMessages         Referral = "group_messages"
	GroupAddresses        Referral = "group_addresses"
	SnippetPost           Referral = "snippet_post"
	SnippetIm             Referral = "snippet_im"
	Clips                 Referral = "clips"
	CommentsListClip      Referral = "comments_list_clip"
	Im                    Referral = "im"
	ImChat                Referral = "im_chat"
	Notifications         Referral = "notifications"
	NotificationsGrouped  Referral = "notifications_grouped"
	NotificationsAuto     Referral = "notifications_auto"
	SuperApp              Referral = "super_app"
	HomeScreen            Referral = "home_screen"
	Menu                  Referral = "menu"
	SnippedPost           Referral = "snipped_post"
	Story                 Referral = "story"
	StoryReply            Referral = "story_reply"
	StoryViewer           Referral = "story_viewer"
	Profile               Referral = "profile"
	ArticleRead           Referral = "article_read"
	MusicPlaylist         Referral = "music_playlist"
	VideoCarousel         Referral = "video_carousel"
	PhotoBrowser          Referral = "photo_browser"
	ShoppingCenter        Referral = "shopping_center"
	MarketItem            Referral = "market_item"
	LeftNav               Referral = "left_nav"
	QuickSearch           Referral = "quick_search"
	Widget                Referral = "widget"
	Other                 Referral = "other"
	Showcase              Referral = "showcase"
)

// Role in the community from which the application is launched.
type Role string

// Possible values.
const (
	RoleNone   Role = "none"
	RoleMember Role = "member"
	RoleModer  Role = "moder"
	RoleEditor Role = "editor"
	RoleAdmin  Role = "admin"
)

// Platform from which the service is launched.
type Platform string

// Possible values.
const (
	MobileAndroid          Platform = "mobile_android"
	MobileIPhone           Platform = "mobile_iphone"
	MobileWeb              Platform = "mobile_web"
	DesktopWeb             Platform = "desktop_web"
	MobileAndroidMessenger Platform = "mobile_android_messenger"
	MobileIPhoneMessenger  Platform = "mobile_iphone_messenger"
)

// Client from which the service is launched.
type Client string

// Possible values.
const (
	ClientOk = "ok"
)

// Params service launch parameters.
type Params struct {
	VkUserID                  int      `schema:"vk_user_id"`
	VkAppID                   int      `schema:"vk_app_id"`
	VkIsAppUser               bool     `schema:"vk_is_app_user"`
	VkAreNotificationsEnabled bool     `schema:"vk_are_notifications_enabled"`
	VkIsFavorite              bool     `schema:"vk_is_favorite"`
	VkLanguage                string   `schema:"vk_language"`
	VkRef                     Referral `schema:"vk_ref"`
	VkAccessTokenSettings     string   `schema:"vk_access_token_settings"`
	VkGroupID                 int      `schema:"vk_group_id"`
	VkViewerGroupRole         Role     `schema:"vk_viewer_group_role"`
	VkPlatform                Platform `schema:"vk_platform"`
	VkTs                      string   `schema:"vk_ts"`
	VkClient                  Client   `schema:"vk_client"`
	Sign                      string   `schema:"sign"`
}

func (p *Params) set(key string, value string) {
	switch key {
	case "vk_user_id":
		if v, err := strconv.Atoi(value); err == nil {
			p.VkUserID = v
		}
	case "vk_app_id":
		if v, err := strconv.Atoi(value); err == nil {
			p.VkAppID = v
		}
	case "vk_is_app_user":
		p.VkIsAppUser = value == "1"
	case "vk_are_notifications_enabled":
		p.VkAreNotificationsEnabled = value == "1"
	case "vk_is_favorite":
		p.VkIsFavorite = value == "1"
	case "vk_language":
		p.VkLanguage = value
	case "vk_ref":
		p.VkRef = Referral(value)
	case "vk_access_token_settings":
		p.VkAccessTokenSettings = value
	case "vk_group_id":
		if v, err := strconv.Atoi(value); err == nil {
			p.VkGroupID = v
		}
	case "vk_viewer_group_role":
		p.VkViewerGroupRole = Role(value)
	case "vk_platform":
		p.VkPlatform = Platform(value)
	case "vk_ts":
		p.VkTs = value
	case "vk_client":
		p.VkClient = Client(value)
	case "sign":
		p.Sign = value
	}
}