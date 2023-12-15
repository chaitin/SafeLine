package service

import (
	"context"
	"log"
	"net/http"
	"slices"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/shurcooL/githubv4"
)

type LabelName = string

const (
	LabelNameEnhancement LabelName = "enhancement"
	LabelNameInProgress  LabelName = "in progress"
	LabelNameReleased    LabelName = "released"
)

type RoadmapLabelName = string

const (
	RoadmapLabelNameInConsideration RoadmapLabelName = "in_consideration"
	RoadmapLabelNameInProgress      RoadmapLabelName = "in_progress"
	RoadmapLabelNameReleased        RoadmapLabelName = "released"
)

type Label struct {
	Name  string `json:"name"`
	Color string `json:"color"`
}

type User struct {
	Login     string `json:"login"`
	AvatarUrl string `json:"avatar_url"`
}

type IssueState = string

const (
	IssueStateOpened IssueState = "OPEN"
	IssueStateClosed IssueState = "CLOSED"
)

// Issue represents a GitHub issue with minimal fields.
type Issue struct {
	ID            string     `json:"id"`
	Title         string     `json:"title"`
	Body          string     `json:"-"`
	State         IssueState `json:"state"`
	Url           string     `json:"url"`
	Labels        []Label    `json:"labels"`
	CommentCount  int        `json:"comment_count"`
	ThumbsUpCount int        `json:"thumbs_up"`
	Author        User       `json:"author"`
	CreatedAt     int64      `json:"created_at"`
	UpdatedAt     int64      `json:"updated_at"`
}

func (i Issue) InConsideration() bool {
	if i.State != IssueStateOpened {
		return false
	}
	if !slices.Contains(i.LabelNames(), LabelNameEnhancement) {
		return false
	}
	if slices.Contains(i.LabelNames(), LabelNameInProgress) {
		return false
	}
	if slices.Contains(i.LabelNames(), LabelNameReleased) {
		return false
	}
	return true
}

func (i Issue) InProgress() bool {
	return i.State == IssueStateOpened && slices.Contains(i.LabelNames(), LabelNameInProgress)
}

func (i Issue) Released() bool {
	return slices.Contains(i.LabelNames(), LabelNameReleased)
}

func (i Issue) LabelNames() []string {
	var names []string
	for _, v := range i.Labels {
		names = append(names, v.Name)
	}
	return names
}

// Discussion represents a GitHub discussion.
type Discussion struct {
	ID            string   `json:"id"`
	Url           string   `json:"url"`
	Title         string   `json:"title"`
	BodyText      string   `json:"-"`
	Labels        []Label  `json:"labels"`
	ThumbsUpCount int      `json:"thumbs_up"`
	CommentCount  int      `json:"comment_count"`
	UpvoteCount   int      `json:"upvote_count"`
	Author        User     `json:"author"`
	CommentUsers  []User   `json:"comment_users"`
	CreatedAt     int64    `json:"created_at"`
	IsAnswered    bool     `json:"is_answered"`
	Category      Category `json:"category"`
}

type Category struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Emoji     string `json:"emoji"`
	EmojiHTML string `json:"emoji_html" graphql:"emojiHTML"`
}

type Repo struct {
	ID        string `json:"id"`
	StarCount int    `json:"star_count"`
}

type GitHubAPI interface {
	Query(ctx context.Context, q interface{}, variables map[string]interface{}) error
}

// GitHubService provides methods to interact with the GitHub API.
type GitHubService struct {
	token    string
	cache    sync.Map
	cacheTTL time.Duration
	owner    string
	repo     string
}

type GithubConfig struct {
	Token    string
	Owner    string
	Repo     string
	CacheTTL time.Duration
}

// headerTransport is custom Transport used to add header information to the request
type headerTransport struct {
	transport *http.Transport
	headers   map[string]string
}

// RoundTrip implements the http.RoundTripper interface
func (h *headerTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	for key, value := range h.headers {
		req.Header.Add(key, value)
	}
	return h.transport.RoundTrip(req)
}

// NewGitHubService creates a new instance of GitHubService with authorized client.
func NewGitHubService(cfg *GithubConfig) *GitHubService {
	s := &GitHubService{
		token:    cfg.Token,
		cache:    sync.Map{},
		cacheTTL: cfg.CacheTTL,
		owner:    cfg.Owner,
		repo:     cfg.Repo,
	}
	go s.loop()
	return s
}

func (s *GitHubService) loop() {
	s.refreshCache()

	t := time.NewTicker(s.cacheTTL * time.Minute)
	defer t.Stop()
	for range t.C {
		s.refreshCache()
	}
}

func (s *GitHubService) client(proxy bool) GitHubAPI {
	httpClient := &http.Client{
		Transport: &headerTransport{
			transport: &http.Transport{},
			headers:   map[string]string{"Authorization": "Bearer " + s.token},
		},
	}

	if proxy {
		httpClient.Transport.(*headerTransport).transport = &http.Transport{
			Proxy: http.ProxyFromEnvironment,
		}
	}

	return githubv4.NewClient(httpClient)
}

func (s *GitHubService) request(ctx context.Context, q interface{}, variables map[string]interface{}) (err error) {
	err = s.client(true).Query(ctx, q, variables)
	if err != nil {
		log.Printf("request using proxy fails and falls back to non-proxy mode: %#v", err)
		err = s.client(false).Query(ctx, q, variables)
	}
	return
}

func (s *GitHubService) refreshCache() {
	issues, err := s.fetchIssues(context.Background(), nil)
	if err != nil {
		log.Printf("failed to fetch issues %v", err)
		return
	}
	s.cache.Store("issues", issues)

	discussions, err := s.fetchDiscussions(context.Background(), nil)
	if err != nil {
		log.Printf("failed to fetch discussions %v", err)
		return
	}
	s.cache.Store("discussions", discussions)

	repo, err := s.fetchRepo(context.Background())
	if err != nil {
		log.Printf("failed to fetch repo %v", err)
		return
	}

	s.cache.Store("repo", repo)
}

// GetIssues tries to get the issues from cache; if not available, fetches from GitHub API.
func (s *GitHubService) GetIssues(ctx context.Context, filter string) (map[string][]*Issue, error) {
	cachedIssues, found := s.cache.Load("issues")
	if found {
		return s.filterIssues(cachedIssues.([]*Issue), filter)
	}

	issues, err := s.fetchIssues(ctx, nil)
	if err != nil {
		return nil, err
	}

	return s.filterIssues(issues, filter)
}

func (s *GitHubService) filterIssues(issues []*Issue, filter string) (map[string][]*Issue, error) {
	filteredIssues := issues
	if filter != "" {
		filteredIssues = make([]*Issue, 0)
		for _, issue := range issues {
			if strings.Contains(issue.Title, filter) || strings.Contains(issue.Body, filter) {
				filteredIssues = append(filteredIssues, issue)
			}
		}
	}
	out := make(map[string][]*Issue)
	for _, issue := range filteredIssues {
		if issue.InConsideration() {
			out[RoadmapLabelNameInConsideration] = append(out[RoadmapLabelNameInConsideration], issue)
		}
		if issue.InProgress() {
			out[RoadmapLabelNameInProgress] = append(out[RoadmapLabelNameInProgress], issue)
		}
		if issue.Released() {
			out[RoadmapLabelNameReleased] = append(out[RoadmapLabelNameReleased], issue)
		}
	}
	sort.Slice(out[RoadmapLabelNameInConsideration], func(i, j int) bool {
		return out[RoadmapLabelNameInConsideration][i].ThumbsUpCount > out[RoadmapLabelNameInConsideration][j].ThumbsUpCount
	})
	sort.Slice(out[RoadmapLabelNameInProgress], func(i, j int) bool {
		return out[RoadmapLabelNameInProgress][i].ThumbsUpCount > out[RoadmapLabelNameInProgress][j].ThumbsUpCount
	})
	sort.Slice(out[RoadmapLabelNameReleased], func(i, j int) bool {
		return out[RoadmapLabelNameReleased][i].UpdatedAt > out[RoadmapLabelNameReleased][j].UpdatedAt
	})
	return out, nil
}

// GetRepositoryIssues queries GitHub for issues of a repository.
func (s *GitHubService) fetchIssues(ctx context.Context, afterCursor *githubv4.String) ([]*Issue, error) {
	var query struct {
		Repository struct {
			Issues struct {
				Nodes []struct {
					ID        string
					Title     string
					Body      string
					Url       string
					State     string
					CreatedAt githubv4.DateTime
					UpdatedAt githubv4.DateTime
					Author    User
					Labels    struct {
						Nodes []struct {
							Color string
							Name  string
						}
					} `graphql:"labels(first: 10)"`
					Comments struct {
						TotalCount int
					}
					Reactions struct {
						TotalCount int
					} `graphql:"reactions(content: THUMBS_UP)"`
				}
				PageInfo struct {
					EndCursor   githubv4.String
					HasNextPage bool
				}
			} `graphql:"issues(first: 100, after: $afterCursor, orderBy: {field: UPDATED_AT, direction: DESC})"`
		} `graphql:"repository(owner: $owner, name: $name)"`
	}
	variables := map[string]interface{}{
		"owner":       githubv4.String(s.owner),
		"name":        githubv4.String(s.repo),
		"afterCursor": afterCursor,
	}

	err := s.request(ctx, &query, variables)
	if err != nil {
		return nil, err
	}

	issues := make([]*Issue, 0)
	for _, node := range query.Repository.Issues.Nodes {
		issue := &Issue{
			ID:            node.ID,
			Title:         node.Title,
			Body:          node.Body,
			Url:           node.Url,
			State:         node.State,
			CreatedAt:     node.CreatedAt.Unix(),
			UpdatedAt:     node.UpdatedAt.Unix(),
			Author:        node.Author,
			CommentCount:  node.Comments.TotalCount,
			ThumbsUpCount: node.Reactions.TotalCount,
		}
		issue.Labels = make([]Label, len(node.Labels.Nodes))
		for i, label := range node.Labels.Nodes {
			issue.Labels[i] = Label{Name: label.Name, Color: label.Color}
		}
		issues = append(issues, issue)
	}

	if query.Repository.Issues.PageInfo.HasNextPage {
		moreIssues, err := s.fetchIssues(ctx, &query.Repository.Issues.PageInfo.EndCursor)
		if err != nil {
			return nil, err
		}
		issues = append(issues, moreIssues...)
	}

	return issues, nil
}

// GetDiscussions tries to get the discussions from cache; if not available, fetches from GitHub API.
func (s *GitHubService) GetDiscussions(ctx context.Context, filter string) ([]*Discussion, error) {
	if cachedData, found := s.cache.Load("discussions"); found {
		return s.filterDiscussions(cachedData.([]*Discussion), filter)
	}

	discussions, err := s.fetchDiscussions(ctx, nil)
	if err != nil {
		return nil, err
	}
	return s.filterDiscussions(discussions, filter)
}

func (s *GitHubService) filterDiscussions(discussions []*Discussion, filter string) ([]*Discussion, error) {
	if filter != "" {
		filteredDiscussions := make([]*Discussion, 0)
		for _, discussion := range discussions {
			if strings.Contains(discussion.Title, filter) || strings.Contains(discussion.BodyText, filter) {
				filteredDiscussions = append(filteredDiscussions, discussion)
			}
		}
		return filteredDiscussions, nil
	}
	return discussions, nil
}

// fetchDiscussionsFromGitHub queries GitHub for discussions of a repository.
func (s *GitHubService) fetchDiscussions(ctx context.Context, afterCursor *githubv4.String) ([]*Discussion, error) {
	var query struct {
		Repository struct {
			Discussions struct {
				Nodes []struct {
					ID          string
					Url         string
					UpvoteCount int
					Title       string
					BodyText    string
					Author      User
					CreatedAt   githubv4.DateTime
					IsAnswered  bool
					Labels      struct {
						Nodes []struct {
							Color string
							Name  string
						}
					} `graphql:"labels(first: 10)"`
					Reactions struct {
						TotalCount int
					} `graphql:"reactions(content: THUMBS_UP)"`
					Comments struct {
						Nodes []struct {
							Author User
						}
					} `graphql:"comments(first: 10)"`
					Category Category
				}
				PageInfo struct {
					EndCursor   githubv4.String
					HasNextPage bool
				}
			} `graphql:"discussions(first: 100, after: $afterCursor, orderBy: {field: CREATED_AT, direction: DESC})"`
		} `graphql:"repository(owner: $owner, name: $name)"`
	}
	variables := map[string]interface{}{
		"owner":       githubv4.String(s.owner),
		"name":        githubv4.String(s.repo),
		"afterCursor": afterCursor,
	}

	err := s.request(ctx, &query, variables)
	if err != nil {
		return nil, err
	}

	discussions := make([]*Discussion, 0)
	for _, node := range query.Repository.Discussions.Nodes {
		discussion := &Discussion{
			ID:       node.ID,
			Url:      node.Url,
			Title:    node.Title,
			BodyText: node.BodyText,
		}
		discussion.Labels = make([]Label, len(node.Labels.Nodes))
		for i, label := range node.Labels.Nodes {
			discussion.Labels[i] = Label{Name: label.Name, Color: label.Color}
		}
		exist := make(map[string]struct{})
		discussion.CommentUsers = make([]User, 0, len(node.Comments.Nodes))
		discussion.CommentUsers = append(discussion.CommentUsers, node.Author)
		exist[node.Author.Login] = struct{}{}
		for _, comment := range node.Comments.Nodes {
			if _, ok := exist[comment.Author.Login]; ok {
				continue
			}
			exist[comment.Author.Login] = struct{}{}
			discussion.CommentUsers = append(discussion.CommentUsers, comment.Author)
		}
		discussion.ThumbsUpCount = node.Reactions.TotalCount
		discussion.CommentCount = len(node.Comments.Nodes)
		discussion.UpvoteCount = node.UpvoteCount
		discussion.Author = node.Author
		discussion.CreatedAt = node.CreatedAt.Unix()
		discussion.IsAnswered = node.IsAnswered
		discussion.Category = node.Category
		discussions = append(discussions, discussion)
	}

	if query.Repository.Discussions.PageInfo.HasNextPage {
		moreDiscussions, err := s.fetchDiscussions(ctx, &query.Repository.Discussions.PageInfo.EndCursor)
		if err != nil {
			return nil, err
		}
		discussions = append(discussions, moreDiscussions...)
	}

	return discussions, nil
}

func (s *GitHubService) GetRepo(ctx context.Context) (*Repo, error) {
	if cachedData, found := s.cache.Load("repo"); found {
		return cachedData.(*Repo), nil
	}

	repo, err := s.fetchRepo(ctx)
	if err != nil {
		return nil, err
	}

	s.cache.Store("repo", repo)

	return repo, nil
}

func (s *GitHubService) fetchRepo(ctx context.Context) (*Repo, error) {
	var query struct {
		Repository struct {
			ID             string
			StargazerCount int
		} `graphql:"repository(owner: $owner, name: $name)"`
	}
	variables := map[string]interface{}{
		"owner": githubv4.String(s.owner),
		"name":  githubv4.String(s.repo),
	}
	err := s.request(ctx, &query, variables)
	if err != nil {
		return nil, err
	}
	return &Repo{ID: query.Repository.ID, StarCount: query.Repository.StargazerCount}, nil
}
