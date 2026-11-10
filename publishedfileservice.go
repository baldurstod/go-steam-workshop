package steamworkshop

import (
	"errors"
	"fmt"
	"net/url"
	"reflect"
	"strconv"
)

const publishedFileServiceQueryFiles = "https://api.steampowered.com/IPublishedFileService/QueryFiles/v1/"
const publishedFileServiceGetDetails = "https://api.steampowered.com/IPublishedFileService/GetDetails/v1/"

type EPublishedFileQueryType int

const (
	K_PublishedFileQueryType_RankedByVote                                  EPublishedFileQueryType = 0
	K_PublishedFileQueryType_RankedByPublicationDate                       EPublishedFileQueryType = 1
	K_PublishedFileQueryType_AcceptedForGameRankedByAcceptanceDate         EPublishedFileQueryType = 2
	K_PublishedFileQueryType_RankedByTrend                                 EPublishedFileQueryType = 3
	K_PublishedFileQueryType_FavoritedByFriendsRankedByPublicationDate     EPublishedFileQueryType = 4
	K_PublishedFileQueryType_CreatedByFriendsRankedByPublicationDate       EPublishedFileQueryType = 5
	K_PublishedFileQueryType_RankedByNumTimesReported                      EPublishedFileQueryType = 6
	K_PublishedFileQueryType_CreatedByFollowedUsersRankedByPublicationDate EPublishedFileQueryType = 7
	K_PublishedFileQueryType_NotYetRated                                   EPublishedFileQueryType = 8
	K_PublishedFileQueryType_RankedByTotalUniqueSubscriptions              EPublishedFileQueryType = 9
	K_PublishedFileQueryType_RankedByTotalVotesAsc                         EPublishedFileQueryType = 10
	K_PublishedFileQueryType_RankedByVotesUp                               EPublishedFileQueryType = 11
	K_PublishedFileQueryType_RankedByTextSearch                            EPublishedFileQueryType = 12
	K_PublishedFileQueryType_RankedByPlaytimeTrend                         EPublishedFileQueryType = 13
	K_PublishedFileQueryType_RankedByTotalPlaytime                         EPublishedFileQueryType = 14
	K_PublishedFileQueryType_RankedByAveragePlaytimeTrend                  EPublishedFileQueryType = 15
	K_PublishedFileQueryType_RankedByLifetimeAveragePlaytime               EPublishedFileQueryType = 16
	K_PublishedFileQueryType_RankedByPlaytimeSessionsTrend                 EPublishedFileQueryType = 17
	K_PublishedFileQueryType_RankedByLifetimePlaytimeSessions              EPublishedFileQueryType = 18
	K_PublishedFileQueryType_RankedByInappropriateContentRating            EPublishedFileQueryType = 19
	K_PublishedFileQueryType_RankedByBanContentCheck                       EPublishedFileQueryType = 20
	K_PublishedFileQueryType_RankedByLastUpdatedDate                       EPublishedFileQueryType = 21
)

type EPublishedFileInfoMatchingFileType int

const (
	K_PFI_MatchingFileType_Items                   EPublishedFileInfoMatchingFileType = 0  //	Items.
	K_PFI_MatchingFileType_Collections             EPublishedFileInfoMatchingFileType = 1  //	A collection of Workshop items.
	K_PFI_MatchingFileType_Art                     EPublishedFileInfoMatchingFileType = 2  //	Artwork.
	K_PFI_MatchingFileType_Videos                  EPublishedFileInfoMatchingFileType = 3  //	Videos.
	K_PFI_MatchingFileType_Screenshots             EPublishedFileInfoMatchingFileType = 4  //	Screenshots.
	K_PFI_MatchingFileType_CollectionEligible      EPublishedFileInfoMatchingFileType = 5  //	Items that can be put inside a collection.
	K_PFI_MatchingFileType_Games                   EPublishedFileInfoMatchingFileType = 6  //	Unused.
	K_PFI_MatchingFileType_Software                EPublishedFileInfoMatchingFileType = 7  //	Unused
	K_PFI_MatchingFileType_Concepts                EPublishedFileInfoMatchingFileType = 8  //	Unused
	K_PFI_MatchingFileType_GreenlightItems         EPublishedFileInfoMatchingFileType = 9  //	Unused
	K_PFI_MatchingFileType_AllGuides               EPublishedFileInfoMatchingFileType = 10 //	Guides.
	K_PFI_MatchingFileType_WebGuides               EPublishedFileInfoMatchingFileType = 11 //	Steam web guide.
	K_PFI_MatchingFileType_IntegratedGuides        EPublishedFileInfoMatchingFileType = 12 //	Application integrated guide.
	K_PFI_MatchingFileType_UsableInGame            EPublishedFileInfoMatchingFileType = 13 //
	K_PFI_MatchingFileType_Merch                   EPublishedFileInfoMatchingFileType = 14 //	Workshop merchandise meant to be voted on for the purpose of being sold
	K_PFI_MatchingFileType_ControllerBindings      EPublishedFileInfoMatchingFileType = 15 //	Steam Controller bindings.
	K_PFI_MatchingFileType_SteamworksAccessInvites EPublishedFileInfoMatchingFileType = 16 //	Used internally.
	K_PFI_MatchingFileType_Items_Mtx               EPublishedFileInfoMatchingFileType = 17 //	Workshop items that can be sold in-game.
	K_PFI_MatchingFileType_Items_ReadyToUse        EPublishedFileInfoMatchingFileType = 18 //	Workshop items that can be used right away by the user.
	K_PFI_MatchingFileType_WorkshopShowcase        EPublishedFileInfoMatchingFileType = 19 //
	K_PFI_MatchingFileType_GameManagedItems        EPublishedFileInfoMatchingFileType = 20 //	Managed completely by the game, not the user, and not shown on the web.
)

type pair struct {
	Value  any
	active bool
}

// Field are exported to be usable in reflect
type queryFilesOptions struct {
	QueryType              pair
	Page                   pair
	NumPerPage             pair
	CreatorAppId           pair
	RequiredTags           pair
	ExcludedTags           pair
	MatchAllTags           pair
	RequiredFlags          pair
	ExcludedFlags          pair
	SearchText             pair
	FileType               pair
	ChildPublishedFileId   pair
	Days                   pair
	IncludeRecentVotesOnly pair
	CacheMaxAgeSeconds     pair
	Language               pair
	//RequiredKvTags         pair
	TotalOnly              pair
	IdsOnly                pair
	ReturnVoteData         pair
	ReturnTags             pair
	ReturnKvTags           pair
	ReturnPreviews         pair
	ReturnChildren         pair
	ReturnShortDescription pair
	ReturnForSaleData      pair
	ReturnMetadata         pair
	ReturnPlaytimeStats    pair
}

type queryFilesParam struct {
	t          string
	queryParam string
}

var queryFilesParams = map[string]queryFilesParam{
	"QueryType":              {"EPublishedFileQueryType", "query_type"},
	"Page":                   {"uint", "page"},
	"NumPerPage":             {"uint", "numperpage"},
	"CreatorAppId":           {"uint", "creator_appid"},
	"RequiredTags":           {"string", "requiredtags"},
	"ExcludedTags":           {"string", "excludedtags"},
	"MatchAllTags":           {"bool", "match_all_tags"},
	"RequiredFlags":          {"string", "required_flags"},
	"ExcludedFlags":          {"string", "omitted_flags"},
	"SearchText":             {"string", "search_text"},
	"FileType":               {"EPublishedFileInfoMatchingFileType", "filetype"},
	"ChildPublishedFileId":   {"uint64", "child_publishedfileid"},
	"Days":                   {"uint", "days"},
	"IncludeRecentVotesOnly": {"bool", "include_recent_votes_only"},
	"CacheMaxAgeSeconds":     {"uint", "cache_max_age_seconds"},
	"Language":               {"int", "language"},
	//"RequiredKvTags":         {"int", "required_kv_tags"},
	"TotalOnly":              {"bool", "totalonly"},
	"IdsOnly":                {"bool", "ids_only"},
	"ReturnVoteData":         {"bool", "return_vote_data"},
	"ReturnTags":             {"bool", "return_tags"},
	"ReturnKvTags":           {"bool", "return_kv_tags"},
	"ReturnPreviews":         {"bool", "return_previews"},
	"ReturnChildren":         {"bool", "return_children"},
	"ReturnShortDescription": {"bool", "return_short_description"},
	"ReturnForSaleData":      {"bool", "return_for_sale_data"},
	"ReturnMetadata":         {"bool", "return_metadata"},
	"ReturnPlaytimeStats":    {"uint", "return_playtime_stats"},
}

type QueryFilesOption func(*queryFilesOptions)

func getQueryFilesOptions(opts ...QueryFilesOption) queryFilesOptions {
	options := queryFilesOptions{
		NumPerPage: pair{uint(100), true},
	}

	for _, fn := range opts {
		fn(&options)
	}

	return options
}

func (options queryFilesOptions) getUrlValues() (url.Values, error) {

	v := reflect.ValueOf(options)
	typeOfS := v.Type()

	values := url.Values{}

	for i := 0; i < v.NumField(); i++ {
		name := typeOfS.Field(i).Name

		active := v.Field(i).Field(1).Bool()
		if !active {
			continue
		}

		value := v.Field(i).Field(0).Interface()
		if value == nil {
			continue
		}

		p, ok := queryFilesParams[name]
		if !ok {
			return nil, errors.New("unable to find param " + name)
		}

		switch p.t {
		case "EPublishedFileQueryType":
			v, ok := value.(EPublishedFileQueryType)
			if !ok {
				return nil, errors.New("unable to cast EPublishedFileQueryType")
			}
			values[p.queryParam] = []string{strconv.Itoa(int(v))}
		case "EPublishedFileInfoMatchingFileType":
			v, ok := value.(EPublishedFileInfoMatchingFileType)
			if !ok {
				return nil, errors.New("unable to cast EPublishedFileInfoMatchingFileType")
			}
			values[p.queryParam] = []string{strconv.Itoa(int(v))}
		case "int":
			v, ok := value.(int)
			if !ok {
				return nil, errors.New("unable to cast int")
			}
			values[p.queryParam] = []string{strconv.FormatInt(int64(v), 10)}
		case "uint":
			v, ok := value.(uint)
			if !ok {
				return nil, errors.New("unable to cast uint")
			}
			values[p.queryParam] = []string{strconv.FormatUint(uint64(v), 10)}
		case "uint64":
			v, ok := value.(uint64)
			if !ok {
				return nil, errors.New("unable to cast uint")
			}
			values[p.queryParam] = []string{strconv.FormatUint(v, 10)}
		case "string":
			v, ok := value.(string)
			if !ok {
				return nil, errors.New("unable to cast uint")
			}
			values[p.queryParam] = []string{v}
		case "bool":
			v, ok := value.(bool)
			if !ok {
				return nil, errors.New("unable to cast uint")
			}
			boolValue := "0"
			if v {
				boolValue = "1"
			}
			values[p.queryParam] = []string{boolValue}
		default:
			return nil, errors.New("unhandled type in GetUrlValues: " + p.t)
		}
	}
	return values, nil
}

func WithQueryType(queryType EPublishedFileQueryType) QueryFilesOption {
	return func(o *queryFilesOptions) {
		o.QueryType = pair{queryType, true}
	}
}

func WithPage(page uint) QueryFilesOption {
	return func(o *queryFilesOptions) {
		o.Page = pair{page, true}
	}
}

func WithNumPerPage(numPerPage uint) QueryFilesOption {
	return func(o *queryFilesOptions) {
		o.NumPerPage = pair{numPerPage, true}
	}
}

func WithCreatorAppId(creatorAppId uint) QueryFilesOption {
	return func(o *queryFilesOptions) {
		o.CreatorAppId = pair{creatorAppId, true}
	}
}

func WithRequiredTags(requiredTags string) QueryFilesOption {
	return func(o *queryFilesOptions) {
		o.RequiredTags = pair{requiredTags, true}
	}
}

func WithExcludedTags(excludedTags string) QueryFilesOption {
	return func(o *queryFilesOptions) {
		o.ExcludedTags = pair{excludedTags, true}
	}
}

func WithMatchAllTags() QueryFilesOption {
	return func(o *queryFilesOptions) {
		o.MatchAllTags = pair{true, true}
	}
}

func WithRequiredFlags(requiredFlags string) QueryFilesOption {
	return func(o *queryFilesOptions) {
		o.RequiredFlags = pair{requiredFlags, true}
	}
}

func WithExcludedFlags(excludedFlags string) QueryFilesOption {
	return func(o *queryFilesOptions) {
		o.ExcludedFlags = pair{excludedFlags, true}
	}
}

func WithSearchText(searchText string) QueryFilesOption {
	return func(o *queryFilesOptions) {
		o.SearchText = pair{searchText, true}
	}
}

func WithFileType(fileType EPublishedFileInfoMatchingFileType) QueryFilesOption {
	return func(o *queryFilesOptions) {
		o.FileType = pair{fileType, true}
	}
}

func WithChildPublishedFileId(childPublishedFileId uint64) QueryFilesOption {
	return func(o *queryFilesOptions) {
		o.ChildPublishedFileId = pair{childPublishedFileId, true}
	}
}

func WithDays(days uint) QueryFilesOption {
	return func(o *queryFilesOptions) {
		o.Days = pair{days, true}
	}
}

func WithIncludeRecentVotesOnly() QueryFilesOption {
	return func(o *queryFilesOptions) {
		o.IncludeRecentVotesOnly = pair{true, true}
	}
}

func WithCacheMaxAgeSeconds(cacheMaxAgeSeconds uint) QueryFilesOption {
	return func(o *queryFilesOptions) {
		o.CacheMaxAgeSeconds = pair{cacheMaxAgeSeconds, true}
	}
}

func WithLanguage(language int) QueryFilesOption {
	return func(o *queryFilesOptions) {
		o.Language = pair{language, true}
	}
}

/*
func WithRequiredKvTags(requiredKvTags int) QueryFilesOption {
	return func(o *queryFilesOptions) {
		o.RequiredKvTags = pair{requiredKvTags, true}
	}
}
*/

func WithTotalOnly() QueryFilesOption {
	return func(o *queryFilesOptions) {
		o.TotalOnly = pair{true, true}
	}
}

func WithIdsOnly() QueryFilesOption {
	return func(o *queryFilesOptions) {
		o.IdsOnly = pair{true, true}
	}
}

func WithReturnVoteData() QueryFilesOption {
	return func(o *queryFilesOptions) {
		o.ReturnVoteData = pair{true, true}
	}
}

func WithReturnTags() QueryFilesOption {
	return func(o *queryFilesOptions) {
		o.ReturnTags = pair{true, true}
	}
}

func WithReturnKvTags() QueryFilesOption {
	return func(o *queryFilesOptions) {
		o.ReturnKvTags = pair{true, true}
	}
}

func WithReturnPreviews() QueryFilesOption {
	return func(o *queryFilesOptions) {
		o.ReturnPreviews = pair{true, true}
	}
}

func WithReturnChildren() QueryFilesOption {
	return func(o *queryFilesOptions) {
		o.ReturnChildren = pair{true, true}
	}
}

func WithReturnShortDescription() QueryFilesOption {
	return func(o *queryFilesOptions) {
		o.ReturnShortDescription = pair{true, true}
	}
}

func WithReturnForSaleData() QueryFilesOption {
	return func(o *queryFilesOptions) {
		o.ReturnForSaleData = pair{true, true}
	}
}

func WithReturnMetadata() QueryFilesOption {
	return func(o *queryFilesOptions) {
		o.ReturnMetadata = pair{true, true}
	}
}

func WithReturnPlaytimeStats(returnPlaytimeStats uint) QueryFilesOption {
	return func(o *queryFilesOptions) {
		o.ReturnPlaytimeStats = pair{returnPlaytimeStats, true}
	}
}

func QueryFiles(apiKey string, appId uint, cursor string, opts ...QueryFilesOption) (map[string]any, error) {
	opt := getQueryFilesOptions(opts...)

	values, err := opt.getUrlValues()
	if err != nil {
		return nil, err
	}

	values["key"] = []string{apiKey}
	values["appid"] = []string{strconv.FormatUint(uint64(appId), 10)}
	values["cursor"] = []string{cursor}

	u, err := url.Parse(publishedFileServiceQueryFiles)
	if err != nil {
		return nil, errors.New("can't parse steamApiQueryFiles")
	}
	u.RawQuery = values.Encode()
	response, err := getApi(u)
	if err != nil {
		return nil, fmt.Errorf("failed to get api response in readWorkshopPage: <%w>", err)
	}

	return response, nil
}
