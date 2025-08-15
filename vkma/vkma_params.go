// Package vkma provides types and functionality for handling VK Mini Apps
// launch parameters. It includes all possible values for referral sources,
// user roles, platforms, and clients as defined in the VK API documentation.
package vkma

import (
	"strconv"
)

// Referral represents the source from which the VK Mini App was launched.
// These values indicate where in the VK ecosystem the app was opened from.
type Referral string

// Comprehensive list of possible referral sources in VK ecosystem.
// These constants cover all documented launch points for VK Mini Apps.
const (
	// Catalog and its subcategories
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

	// Board topics
	BoardTopicAll  Referral = "board_topic_all"
	BoardTopicView Referral = "board_topic_view"

	// Feed related
	Feed          Referral = "feed"
	FeedPost      Referral = "feed_post"
	FeedComments  Referral = "feed_comments"

	// Featuring sections
	FeaturingDiscover Referral = "featuring_discover"
	FeaturingMenu     Referral = "featuring_menu"
	FeaturingNew      Referral = "featuring_new"

	// Favorites
	Fave      Referral = "fave"
	FaveLinks Referral = "fave_links"
	FavePosts Referral = "fave_posts"

	// Group related
	Group          Referral = "group"
	GroupMenu      Referral = "group_menu"
	GroupMessages  Referral = "group_messages"
	GroupAddresses Referral = "group_addresses"

	// Snippets
	SnippetPost Referral = "snippet_post"
	SnippetIm   Referral = "snippet_im"

	// Clips and videos
	Clips            Referral = "clips"
	CommentsListClip Referral = "comments_list_clip"

	// Messaging
	Im     Referral = "im"
	ImChat Referral = "im_chat"

	// Notifications
	Notifications        Referral = "notifications"
	NotificationsGrouped Referral = "notifications_grouped"
	NotificationsAuto    Referral = "notifications_auto"

	// App sections
	SuperApp   Referral = "super_app"
	HomeScreen Referral = "home_screen"
	Menu       Referral = "menu"

	// Content types
	SnippedPost  Referral = "snipped_post"
	Story        Referral = "story"
	StoryReply   Referral = "story_reply"
	StoryViewer  Referral = "story_viewer"
	Profile      Referral = "profile"
	ArticleRead  Referral = "article_read"
	MusicPlaylist Referral = "music_playlist"
	VideoCarousel Referral = "video_carousel"
	PhotoBrowser  Referral = "photo_browser"

	// Shopping and commerce
	ShoppingCenter Referral = "shopping_center"
	MarketItem     Referral = "market_item"

	// Navigation
	LeftNav     Referral = "left_nav"
	QuickSearch Referral = "quick_search"

	// Miscellaneous
	Widget   Referral = "widget"
	Other    Referral = "other"
	Showcase Referral = "showcase"
)

// Role represents the user's role in the community from which the app was launched.
type Role string

// Possible user roles in communities
const (
	RoleNone   Role = "none"    // No special role
	RoleMember Role = "member"  // Regular community member
	RoleModer  Role = "moder"   // Community moderator
	RoleEditor Role = "editor"  // Community editor
	RoleAdmin  Role = "admin"   // Community administrator
)

// Platform represents the device platform from which the app was launched.
type Platform string

// Supported platforms for VK Mini Apps
const (
	MobileAndroid          Platform = "mobile_android"            // Android native app
	MobileIPhone           Platform = "mobile_iphone"             // iOS native app
	MobileWeb              Platform = "mobile_web"                // Mobile web browser
	DesktopWeb             Platform = "desktop_web"               // Desktop web browser
	MobileAndroidMessenger Platform = "mobile_android_messenger"  // Android Messenger app
	MobileIPhoneMessenger  Platform = "mobile_iphone_messenger"   // iOS Messenger app
)

// Client represents the specific client application from which the app was launched.
type Client string

// Supported VK clients
const (
	ClientOk Client = "ok"  // Odnoklassniki client
)

// Params contains all possible launch parameters for a VK Mini App.
// The struct tags match the expected query parameter names from VK.
type Params struct {
	VkUserID                  int      `schema:"vk_user_id"`                   // User ID in VK
	VkAppID                   int      `schema:"vk_app_id"`                    // App ID in VK
	VkIsAppUser               bool     `schema:"vk_is_app_user"`               // Is user logged in through VK
	VkAreNotificationsEnabled bool     `schema:"vk_are_notifications_enabled"` // Are notifications enabled
	VkIsFavorite              bool     `schema:"vk_is_favorite"`              // Is app in user's favorites
	VkLanguage                string   `schema:"vk_language"`                 // User's language preference
	VkRef                     Referral `schema:"vk_ref"`                      // Launch referral source
	VkAccessTokenSettings     string   `schema:"vk_access_token_settings"`    // Granted permissions
	VkGroupID                 int      `schema:"vk_group_id"`                 // Community ID if launched from group
	VkViewerGroupRole         Role     `schema:"vk_viewer_group_role"`        // User's role in the community
	VkPlatform                Platform `schema:"vk_platform"`                 // Platform type
	VkTs                      string   `schema:"vk_ts"`                       // Timestamp of launch
	VkClient                  Client   `schema:"vk_client"`                   // Specific client used
	Sign                      string   `schema:"sign"`                        // Security signature
}

// set assigns a value to the appropriate field in Params based on the key.
// It handles type conversion and validation for all supported parameters.
//
// Parameters:
//   - key: The parameter name (must match schema tags exactly)
//   - value: The string value to be parsed and assigned
//
// The method silently ignores unsupported keys and parsing errors,
// leaving fields at their zero values when parsing fails.
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
		p.VkIsAppUser = value == "1"  // VK uses "1" for true, empty or other for false
	case "vk_are_notifications_enabled":
		p.VkAreNotificationsEnabled = value == "1"
	case "vk_is_favorite":
		p.VkIsFavorite = value == "1"
	case "vk_language":
		p.VkLanguage = value  // No validation for language codes
	case "vk_ref":
		p.VkRef = Referral(value)  // No validation against known referrals
	case "vk_access_token_settings":
		p.VkAccessTokenSettings = value  // Comma-separated permissions
	case "vk_group_id":
		if v, err := strconv.Atoi(value); err == nil {
			p.VkGroupID = v
		}
	case "vk_viewer_group_role":
		p.VkViewerGroupRole = Role(value)  // No validation against known roles
	case "vk_platform":
		p.VkPlatform = Platform(value)  // No validation against known platforms
	case "vk_ts":
		p.VkTs = value  // Timestamp as string
	case "vk_client":
		p.VkClient = Client(value)  // No validation against known clients
	case "sign":
		p.Sign = value  // Security signature as-is
	}
}