package medium

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/indikator/aggregator_lets_go/internal/parser"
	"github.com/indikator/aggregator_lets_go/model"
)

const (
	limitQuantityStates = 25
)

type articlesParser struct {
	url string
}

type StatesMedium struct {
	Data DataStates
}

type DataStates struct {
	TagFeed TagFeed
}

type TagFeed struct {
	Items []Item
}

type Item struct {
	FeedId string
	Post   Post
}

type Post struct {
	Id                     string
	Creator                Creator
	ExtendedPreviewContent ExtendedPreviewContent
	FirstPublishedAt       int64
	Title                  string
	MediumUrl              string
}

type Creator struct {
	Name string
}

type ExtendedPreviewContent struct {
	Subtitle string
}

// create an instance of articles parser
func NewParser(cfg parser.Config) parser.ArticlesParser {
	return &articlesParser{
		url: cfg.URL,
	}
}

func init() {
	parser.RegisterParser("medium.com", NewParser)
}

// / parse all articles that were created earler than the target date
func (p *articlesParser) ParseAfter(maxDate time.Time) (articles []model.Article, err error) {

	var (
		dateLastState   int64
		initNumberState int
	)

	maxDateUnix := maxDate.UnixMilli()
	initialRequest := true

	for true {

		states, err := getStates(p.url, initialRequest, initNumberState)
		if err != nil {
			return nil, err
		}

		if len(states) == 0 && !initialRequest {
			log.Printf("Problems in site - medium.com - at %v", time.Now())
			return articles, nil
		}

		for _, state := range states {
			lenItemsStates := len(state.Data.TagFeed.Items)
			for index, itemState := range state.Data.TagFeed.Items {
				article := model.Article{
					Title:       itemState.Post.Title,
					URL:         itemState.Post.MediumUrl,
					Created:     time.Unix(itemState.Post.FirstPublishedAt/1000, 0),
					Author:      itemState.Post.Creator.Name,
					Description: itemState.Post.ExtendedPreviewContent.Subtitle,
				}
				articles = append(articles, article)
				if index == (lenItemsStates - 1) {
					dateLastState = itemState.Post.FirstPublishedAt / 1000
				}
			}
		}

		initialRequest = false

		// For implement pagination, this required in request
		initNumberState += limitQuantityStates

		if maxDateUnix >= dateLastState {
			break
		}

	}

	return
}

func getStates(url string, initialRequest bool, initNumberState int) (states []StatesMedium, err error) {

	responseBody := getBody(initialRequest, initNumberState)

	request, err := http.NewRequest(http.MethodPost, url, responseBody)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Host", "medium.com")
	request.Header.Set("Content-Type", "application/json")

	responce, err := http.DefaultClient.Do(request)
	if err != nil || responce.StatusCode != http.StatusOK {
		return nil, err
	}

	defer responce.Body.Close()

	resBody, err := ioutil.ReadAll(responce.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(resBody, &states)
	if err != nil {
		return nil, err
	}

	return states, nil

}

func getBody(initialRequest bool, initNumberState int) *bytes.Buffer {

	var postBody []byte

	paging := fmt.Sprintf(`{"to":"%d","limit":%d}}`, initNumberState, limitQuantityStates)
	if initialRequest {
		paging = fmt.Sprintf(`{"limit":%d}}`, limitQuantityStates)
	}

	postBody = []byte(fmt.Sprintf(`[{"operationName":"TopicFeedQuery","variables":{"tagSlug":"golang","mode":"NEW","paging": %s,"query":"query TopicFeedQuery($paging: PagingOptions, $tagSlug: String, $mode: TagFeedMode) {\n  tagFeed(paging: $paging, tagSlug: $tagSlug, mode: $mode) {\n    items {\n      ... on TagFeedItem {\n        feedId\n        reason\n        moduleSourceEncoding\n        postProviderExplanation {\n          reason\n          topic {\n            name\n            __typename\n          }\n          __typename\n        }\n        post {\n          ...PostPreview_post\n          __typename\n        }\n        __typename\n      }\n      __typename\n    }\n    pagingInfo {\n      next {\n        limit\n        to\n        __typename\n      }\n      __typename\n    }\n    __typename\n  }\n}\n\nfragment PostPreview_post on Post {\n  id\n  creator {\n    ...PostPreview_user\n    __typename\n    id\n  }\n  collection {\n    ...CardByline_collection\n    __typename\n    id\n  }\n  ...InteractivePostBody_postPreview\n  firstPublishedAt\n  isLocked\n  isSeries\n  isShortform\n  latestPublishedAt\n  inResponseToCatalogResult {\n    __typename\n  }\n  previewImage {\n    id\n    focusPercentX\n    focusPercentY\n    __typename\n  }\n  readingTime\n  sequence {\n    slug\n    __typename\n  }\n  title\n  uniqueSlug\n  visibility\n  ...CardByline_post\n  ...PostFooterActionsBar_post\n  ...InResponseToEntityPreview_post\n  ...PostScrollTracker_post\n  ...ReadMore_post\n  ...HighDensityPreview_post\n  __typename\n}\n\nfragment PostPreview_user on User {\n  __typename\n  name\n  username\n  ...CardByline_user\n  id\n}\n\nfragment CardByline_user on User {\n  __typename\n  id\n  name\n  username\n  mediumMemberAt\n  socialStats {\n    followerCount\n    __typename\n  }\n  ...userUrl_user\n  ...UserMentionTooltip_user\n}\n\nfragment userUrl_user on User {\n  __typename\n  id\n  customDomainState {\n    live {\n      domain\n      __typename\n    }\n    __typename\n  }\n  hasSubdomain\n  username\n}\n\nfragment UserMentionTooltip_user on User {\n  id\n  name\n  username\n  bio\n  imageId\n  mediumMemberAt\n  ...UserAvatar_user\n  ...UserFollowButton_user\n  __typename\n}\n\nfragment UserAvatar_user on User {\n  __typename\n  id\n  imageId\n  mediumMemberAt\n  name\n  username\n  ...userUrl_user\n}\n\nfragment UserFollowButton_user on User {\n  ...UserFollowButtonSignedIn_user\n  ...UserFollowButtonSignedOut_user\n  __typename\n  id\n}\n\nfragment UserFollowButtonSignedIn_user on User {\n  id\n  name\n  __typename\n}\n\nfragment UserFollowButtonSignedOut_user on User {\n  id\n  ...SusiClickable_user\n  __typename\n}\n\nfragment SusiClickable_user on User {\n  ...SusiContainer_user\n  __typename\n  id\n}\n\nfragment SusiContainer_user on User {\n  ...SignInOptions_user\n  ...SignUpOptions_user\n  __typename\n  id\n}\n\nfragment SignInOptions_user on User {\n  id\n  name\n  __typename\n}\n\nfragment SignUpOptions_user on User {\n  id\n  name\n  __typename\n}\n\nfragment CardByline_collection on Collection {\n  __typename\n  id\n  name\n  ...collectionUrl_collection\n}\n\nfragment collectionUrl_collection on Collection {\n  id\n  domain\n  slug\n  __typename\n}\n\nfragment InteractivePostBody_postPreview on Post {\n  extendedPreviewContent(\n    truncationConfig: {previewParagraphsWordCountThreshold: 400, minimumWordLengthForTruncation: 150, truncateAtEndOfSentence: true, showFullImageCaptions: true, shortformPreviewParagraphsWordCountThreshold: 30, shortformMinimumWordLengthForTruncation: 30}\n  ) {\n    bodyModel {\n      ...PostBody_bodyModel\n      __typename\n    }\n    isFullContent\n    __typename\n  }\n  __typename\n  id\n}\n\nfragment PostBody_bodyModel on RichText {\n  sections {\n    name\n    startIndex\n    textLayout\n    imageLayout\n    backgroundImage {\n      id\n      originalHeight\n      originalWidth\n      __typename\n    }\n    videoLayout\n    backgroundVideo {\n      videoId\n      originalHeight\n      originalWidth\n      previewImageId\n      __typename\n    }\n    __typename\n  }\n  paragraphs {\n    id\n    ...PostBodySection_paragraph\n    __typename\n  }\n  ...normalizedBodyModel_richText\n  __typename\n}\n\nfragment PostBodySection_paragraph on Paragraph {\n  name\n  ...PostBodyParagraph_paragraph\n  __typename\n  id\n}\n\nfragment PostBodyParagraph_paragraph on Paragraph {\n  name\n  type\n  ...ImageParagraph_paragraph\n  ...TextParagraph_paragraph\n  ...IframeParagraph_paragraph\n  ...MixtapeParagraph_paragraph\n  ...CodeBlockParagraph_paragraph\n  __typename\n  id\n}\n\nfragment ImageParagraph_paragraph on Paragraph {\n  href\n  layout\n  metadata {\n    id\n    originalHeight\n    originalWidth\n    focusPercentX\n    focusPercentY\n    alt\n    __typename\n  }\n  ...Markups_paragraph\n  ...ParagraphRefsMapContext_paragraph\n  ...PostAnnotationsMarker_paragraph\n  __typename\n  id\n}\n\nfragment Markups_paragraph on Paragraph {\n  name\n  text\n  hasDropCap\n  dropCapImage {\n    ...MarkupNode_data_dropCapImage\n    __typename\n    id\n  }\n  markups {\n    type\n    start\n    end\n    href\n    anchorType\n    userId\n    linkMetadata {\n      httpStatus\n      __typename\n    }\n    __typename\n  }\n  __typename\n  id\n}\n\nfragment MarkupNode_data_dropCapImage on ImageMetadata {\n  ...DropCap_image\n  __typename\n  id\n}\n\nfragment DropCap_image on ImageMetadata {\n  id\n  originalHeight\n  originalWidth\n  __typename\n}\n\nfragment ParagraphRefsMapContext_paragraph on Paragraph {\n  id\n  name\n  text\n  __typename\n}\n\nfragment PostAnnotationsMarker_paragraph on Paragraph {\n  ...PostViewNoteCard_paragraph\n  __typename\n  id\n}\n\nfragment PostViewNoteCard_paragraph on Paragraph {\n  name\n  __typename\n  id\n}\n\nfragment TextParagraph_paragraph on Paragraph {\n  type\n  hasDropCap\n  ...Markups_paragraph\n  ...ParagraphRefsMapContext_paragraph\n  __typename\n  id\n}\n\nfragment IframeParagraph_paragraph on Paragraph {\n  iframe {\n    mediaResource {\n      id\n      iframeSrc\n      iframeHeight\n      iframeWidth\n      title\n      __typename\n    }\n    __typename\n  }\n  layout\n  ...getEmbedlyCardUrlParams_paragraph\n  ...Markups_paragraph\n  __typename\n  id\n}\n\nfragment getEmbedlyCardUrlParams_paragraph on Paragraph {\n  type\n  iframe {\n    mediaResource {\n      iframeSrc\n      __typename\n    }\n    __typename\n  }\n  __typename\n  id\n}\n\nfragment MixtapeParagraph_paragraph on Paragraph {\n  type\n  mixtapeMetadata {\n    href\n    mediaResource {\n      mediumCatalog {\n        id\n        __typename\n      }\n      __typename\n    }\n    __typename\n  }\n  ...GenericMixtapeParagraph_paragraph\n  __typename\n  id\n}\n\nfragment GenericMixtapeParagraph_paragraph on Paragraph {\n  text\n  mixtapeMetadata {\n    href\n    thumbnailImageId\n    __typename\n  }\n  markups {\n    start\n    end\n    type\n    href\n    __typename\n  }\n  __typename\n  id\n}\n\nfragment CodeBlockParagraph_paragraph on Paragraph {\n  codeBlockMetadata {\n    lang\n    mode\n    __typename\n  }\n  __typename\n  id\n}\n\nfragment normalizedBodyModel_richText on RichText {\n  paragraphs {\n    markups {\n      type\n      __typename\n    }\n    ...getParagraphHighlights_paragraph\n    ...getParagraphPrivateNotes_paragraph\n    __typename\n  }\n  sections {\n    startIndex\n    ...getSectionEndIndex_section\n    __typename\n  }\n  ...getParagraphStyles_richText\n  ...getParagraphSpaces_richText\n  __typename\n}\n\nfragment getParagraphHighlights_paragraph on Paragraph {\n  name\n  __typename\n  id\n}\n\nfragment getParagraphPrivateNotes_paragraph on Paragraph {\n  name\n  __typename\n  id\n}\n\nfragment getSectionEndIndex_section on Section {\n  startIndex\n  __typename\n}\n\nfragment getParagraphStyles_richText on RichText {\n  paragraphs {\n    text\n    type\n    __typename\n  }\n  sections {\n    ...getSectionEndIndex_section\n    __typename\n  }\n  __typename\n}\n\nfragment getParagraphSpaces_richText on RichText {\n  paragraphs {\n    layout\n    metadata {\n      originalHeight\n      originalWidth\n      id\n      __typename\n    }\n    type\n    ...paragraphExtendsImageGrid_paragraph\n    __typename\n  }\n  ...getSeriesParagraphTopSpacings_richText\n  ...getPostParagraphTopSpacings_richText\n  __typename\n}\n\nfragment paragraphExtendsImageGrid_paragraph on Paragraph {\n  layout\n  type\n  __typename\n  id\n}\n\nfragment getSeriesParagraphTopSpacings_richText on RichText {\n  paragraphs {\n    id\n    __typename\n  }\n  sections {\n    startIndex\n    __typename\n  }\n  __typename\n}\n\nfragment getPostParagraphTopSpacings_richText on RichText {\n  paragraphs {\n    layout\n    text\n    __typename\n  }\n  sections {\n    startIndex\n    __typename\n  }\n  __typename\n}\n\nfragment CardByline_post on Post {\n  ...DraftStatus_post\n  ...Star_post\n  __typename\n  id\n}\n\nfragment DraftStatus_post on Post {\n  id\n  pendingCollection {\n    id\n    creator {\n      id\n      __typename\n    }\n    ...BoldCollectionName_collection\n    __typename\n  }\n  statusForCollection\n  creator {\n    id\n    __typename\n  }\n  isPublished\n  __typename\n}\n\nfragment BoldCollectionName_collection on Collection {\n  id\n  name\n  __typename\n}\n\nfragment Star_post on Post {\n  id\n  creator {\n    id\n    __typename\n  }\n  __typename\n}\n\nfragment PostFooterActionsBar_post on Post {\n  id\n  visibility\n  isPublished\n  allowResponses\n  postResponses {\n    count\n    __typename\n  }\n  isLimitedState\n  creator {\n    id\n    __typename\n  }\n  collection {\n    id\n    __typename\n  }\n  ...BookmarkButton_post\n  ...MultiVote_post\n  ...SharePostButtons_post\n  ...PostFooterSocialPopover_post\n  ...OverflowMenuButtonWithNegativeSignal_post\n  __typename\n}\n\nfragment BookmarkButton_post on Post {\n  visibility\n  ...SusiClickable_post\n  ...AddToCatalogBookmarkButton_post\n  __typename\n  id\n}\n\nfragment SusiClickable_post on Post {\n  id\n  mediumUrl\n  ...SusiContainer_post\n  __typename\n}\n\nfragment SusiContainer_post on Post {\n  id\n  __typename\n}\n\nfragment AddToCatalogBookmarkButton_post on Post {\n  ...AddToCatalogBase_post\n  __typename\n  id\n}\n\nfragment AddToCatalogBase_post on Post {\n  id\n  __typename\n}\n\nfragment MultiVote_post on Post {\n  id\n  creator {\n    id\n    ...SusiClickable_user\n    __typename\n  }\n  isPublished\n  ...SusiClickable_post\n  collection {\n    id\n    slug\n    __typename\n  }\n  isLimitedState\n  ...MultiVoteCount_post\n  __typename\n}\n\nfragment MultiVoteCount_post on Post {\n  id\n  ...PostVotersNetwork_post\n  __typename\n}\n\nfragment PostVotersNetwork_post on Post {\n  id\n  voterCount\n  recommenders {\n    name\n    __typename\n  }\n  __typename\n}\n\nfragment SharePostButtons_post on Post {\n  id\n  isLimitedState\n  visibility\n  mediumUrl\n  ...SharePostButton_post\n  __typename\n}\n\nfragment SharePostButton_post on Post {\n  id\n  __typename\n}\n\nfragment PostFooterSocialPopover_post on Post {\n  id\n  mediumUrl\n  title\n  ...SharePostButton_post\n  __typename\n}\n\nfragment OverflowMenuButtonWithNegativeSignal_post on Post {\n  id\n  ...OverflowMenuWithNegativeSignal_post\n  ...CreatorActionOverflowPopover_post\n  __typename\n}\n\nfragment OverflowMenuWithNegativeSignal_post on Post {\n  id\n  creator {\n    id\n    __typename\n  }\n  collection {\n    id\n    __typename\n  }\n  ...OverflowMenuItemUndoClaps_post\n  __typename\n}\n\nfragment OverflowMenuItemUndoClaps_post on Post {\n  id\n  clapCount\n  ...ClapMutation_post\n  __typename\n}\n\nfragment ClapMutation_post on Post {\n  __typename\n  id\n  clapCount\n  ...MultiVoteCount_post\n}\n\nfragment CreatorActionOverflowPopover_post on Post {\n  allowResponses\n  id\n  statusForCollection\n  isLocked\n  isPublished\n  clapCount\n  mediumUrl\n  pinnedAt\n  pinnedByCreatorAt\n  curationEligibleAt\n  mediumUrl\n  responseDistribution\n  visibility\n  inResponseToPostResult {\n    __typename\n  }\n  inResponseToCatalogResult {\n    __typename\n  }\n  pendingCollection {\n    id\n    name\n    creator {\n      id\n      __typename\n    }\n    avatar {\n      id\n      __typename\n    }\n    domain\n    slug\n    __typename\n  }\n  creator {\n    id\n    ...MutePopoverOptions_creator\n    ...auroraHooks_publisher\n    __typename\n  }\n  collection {\n    id\n    name\n    creator {\n      id\n      __typename\n    }\n    avatar {\n      id\n      __typename\n    }\n    domain\n    slug\n    ...MutePopoverOptions_collection\n    ...auroraHooks_publisher\n    __typename\n  }\n  ...useIsPinnedInContext_post\n  ...NewsletterV3EmailToSubscribersMenuItem_post\n  ...OverflowMenuItemUndoClaps_post\n  __typename\n}\n\nfragment MutePopoverOptions_creator on User {\n  id\n  __typename\n}\n\nfragment auroraHooks_publisher on Publisher {\n  __typename\n  ... on Collection {\n    isAuroraEligible\n    isAuroraVisible\n    viewerEdge {\n      id\n      isEditor\n      __typename\n    }\n    __typename\n    id\n  }\n  ... on User {\n    isAuroraVisible\n    __typename\n    id\n  }\n}\n\nfragment MutePopoverOptions_collection on Collection {\n  id\n  __typename\n}\n\nfragment useIsPinnedInContext_post on Post {\n  id\n  collection {\n    id\n    __typename\n  }\n  pendingCollection {\n    id\n    __typename\n  }\n  pinnedAt\n  pinnedByCreatorAt\n  __typename\n}\n\nfragment NewsletterV3EmailToSubscribersMenuItem_post on Post {\n  id\n  creator {\n    id\n    newsletterV3 {\n      id\n      subscribersCount\n      __typename\n    }\n    __typename\n  }\n  isNewsletter\n  isAuthorNewsletter\n  __typename\n}\n\nfragment InResponseToEntityPreview_post on Post {\n  id\n  inResponseToEntityType\n  __typename\n}\n\nfragment PostScrollTracker_post on Post {\n  id\n  collection {\n    id\n    __typename\n  }\n  sequence {\n    sequenceId\n    __typename\n  }\n  __typename\n}\n\nfragment ReadMore_post on Post {\n  mediumUrl\n  readingTime\n  ...usePostUrl_post\n  __typename\n  id\n}\n\nfragment usePostUrl_post on Post {\n  id\n  creator {\n    ...userUrl_user\n    __typename\n    id\n  }\n  collection {\n    id\n    domain\n    slug\n    __typename\n  }\n  isSeries\n  mediumUrl\n  sequence {\n    slug\n    __typename\n  }\n  uniqueSlug\n  __typename\n}\n\nfragment HighDensityPreview_post on Post {\n  id\n  title\n  previewImage {\n    id\n    focusPercentX\n    focusPercentY\n    __typename\n  }\n  extendedPreviewContent(\n    truncationConfig: {previewParagraphsWordCountThreshold: 400, minimumWordLengthForTruncation: 150, truncateAtEndOfSentence: true, showFullImageCaptions: true, shortformPreviewParagraphsWordCountThreshold: 30, shortformMinimumWordLengthForTruncation: 30}\n  ) {\n    subtitle\n    __typename\n  }\n  ...HighDensityFooter_post\n  __typename\n}\n\nfragment HighDensityFooter_post on Post {\n  id\n  readingTime\n  tags {\n    ...TopicPill_tag\n    __typename\n  }\n  ...BookmarkButton_post\n  ...ExpandablePostCardOverflowButton_post\n  ...OverflowMenuButtonWithNegativeSignal_post\n  __typename\n}\n\nfragment TopicPill_tag on Tag {\n  __typename\n  id\n  displayTitle\n  normalizedTagSlug\n}\n\nfragment ExpandablePostCardOverflowButton_post on Post {\n  creator {\n    id\n    __typename\n  }\n  ...ExpandablePostCardEditorWriterButton_post\n  ...ExpandablePostCardReaderButton_post\n  __typename\n  id\n}\n\nfragment ExpandablePostCardEditorWriterButton_post on Post {\n  id\n  collection {\n    id\n    name\n    slug\n    __typename\n  }\n  allowResponses\n  clapCount\n  visibility\n  mediumUrl\n  responseDistribution\n  ...useIsPinnedInContext_post\n  ...CopyFriendLinkMenuItem_post\n  ...NewsletterV3EmailToSubscribersMenuItem_post\n  ...OverflowMenuItemUndoClaps_post\n  __typename\n}\n\nfragment CopyFriendLinkMenuItem_post on Post {\n  id\n  __typename\n}\n\nfragment ExpandablePostCardReaderButton_post on Post {\n  id\n  collection {\n    id\n    __typename\n  }\n  creator {\n    id\n    __typename\n  }\n  clapCount\n  ...ClapMutation_post\n  __typename\n}\n"}]`, paging))

	return bytes.NewBuffer(postBody)

}
