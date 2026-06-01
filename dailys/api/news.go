package api

import (
	"context"
	"fmt"
	"html"
	"net/url"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/mmcdole/gofeed"
)

const newsUserAgent = "DailyDashBoard/1.0"

var newsHTMLTagPattern = regexp.MustCompile(`<[^>]+>`)

type NewsFeed struct {
	Name string
	URL  string
}

type NewsItem struct {
	Source      string
	Title       string
	Link        string
	Summary     string
	PublishedAt time.Time
}

func DefaultNewsFeeds() []NewsFeed {
	return []NewsFeed{
		{Name: "BBC News", URL: "https://feeds.bbci.co.uk/news/rss.xml"},
		{Name: "NPR News", URL: "https://feeds.npr.org/1001/rss.xml"},
		{Name: "New York Times", URL: "https://rss.nytimes.com/services/xml/rss/nyt/HomePage.xml"},
	}
}

func DailyNews(feeds []NewsFeed, limit int) ([]NewsItem, error) {
	if len(feeds) == 0 {
		return nil, fmt.Errorf("no news feeds configured")
	}

	parser := gofeed.NewParser()
	parser.Client = httpClient
	parser.UserAgent = newsUserAgent

	var items []NewsItem
	var feedErrors []string

	for _, configuredFeed := range feeds {
		feedURL := strings.TrimSpace(configuredFeed.URL)
		if feedURL == "" {
			continue
		}

		ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
		feed, err := parser.ParseURLWithContext(feedURL, ctx)
		cancel()

		if err != nil {
			feedErrors = append(feedErrors, fmt.Sprintf("%s: %v", newsFeedName(configuredFeed, nil), err))
			continue
		}

		source := newsFeedName(configuredFeed, feed)
		for _, item := range feed.Items {
			if item == nil {
				continue
			}

			newsItem := newsItemFromFeed(source, item)
			if newsItem.Title == "" {
				continue
			}

			items = append(items, newsItem)
		}
	}

	if len(items) == 0 {
		if len(feedErrors) > 0 {
			return nil, fmt.Errorf("could not fetch any news feeds: %s", strings.Join(feedErrors, "; "))
		}

		return nil, fmt.Errorf("no news stories were found")
	}

	items = dedupeNewsItems(items)
	sortNewsItems(items)

	if limit > 0 && len(items) > limit {
		items = items[:limit]
	}

	return items, nil
}

func newsFeedName(configuredFeed NewsFeed, feed *gofeed.Feed) string {
	if configuredFeed.Name != "" {
		return configuredFeed.Name
	}

	if feed != nil && strings.TrimSpace(feed.Title) != "" {
		return cleanNewsText(feed.Title)
	}

	parsedURL, err := url.Parse(configuredFeed.URL)
	if err != nil || parsedURL.Hostname() == "" {
		return "News"
	}

	return strings.TrimPrefix(parsedURL.Hostname(), "www.")
}

func newsItemFromFeed(source string, item *gofeed.Item) NewsItem {
	summary := cleanNewsText(firstNewsText(item.Description, item.Content))
	if summary == cleanNewsText(item.Title) {
		summary = ""
	}

	return NewsItem{
		Source:      source,
		Title:       cleanNewsText(item.Title),
		Link:        strings.TrimSpace(item.Link),
		Summary:     summary,
		PublishedAt: newsItemTime(item),
	}
}

func firstNewsText(values ...string) string {
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value != "" {
			return value
		}
	}

	return ""
}

func cleanNewsText(value string) string {
	value = newsHTMLTagPattern.ReplaceAllString(value, " ")
	value = html.UnescapeString(value)

	return strings.Join(strings.Fields(value), " ")
}

func newsItemTime(item *gofeed.Item) time.Time {
	if item.PublishedParsed != nil {
		return *item.PublishedParsed
	}

	if item.UpdatedParsed != nil {
		return *item.UpdatedParsed
	}

	return time.Time{}
}

func dedupeNewsItems(items []NewsItem) []NewsItem {
	seen := make(map[string]bool)
	deduped := make([]NewsItem, 0, len(items))

	for _, item := range items {
		key := newsItemKey(item)
		if key == "" || seen[key] {
			continue
		}

		seen[key] = true
		deduped = append(deduped, item)
	}

	return deduped
}

func newsItemKey(item NewsItem) string {
	if item.Link != "" {
		return strings.ToLower(item.Link)
	}

	if item.Title != "" {
		return strings.ToLower(item.Source + ":" + item.Title)
	}

	return ""
}

func sortNewsItems(items []NewsItem) {
	sort.SliceStable(items, func(i int, j int) bool {
		left := items[i].PublishedAt
		right := items[j].PublishedAt

		switch {
		case left.IsZero() && !right.IsZero():
			return false
		case !left.IsZero() && right.IsZero():
			return true
		case !left.Equal(right):
			return left.After(right)
		default:
			return items[i].Title < items[j].Title
		}
	})
}
