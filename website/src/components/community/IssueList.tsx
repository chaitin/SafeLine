import React, { useState } from "react";
import {
  Box,
  List,
  ListItem,
  Button,
  Stack,
  Typography,
  Grid,
  Paper,
  IconButton,
  InputBase,
  alpha,
} from "@mui/material";
import Icon from "@/components/Icon";
import SearchIcon from '@mui/icons-material/Search';
import { Labels } from './DiscussionList'
import { formatDate } from '@/common/utils'
import { getIssues } from "@/api";

export default IssueList;

type Issue = {
  id: string
  labels: { name: string, color: string }[]
  thumbs_up: number
  title: string
  url: string
  comment_count: number
  created_at: number
  author: {
    avatar_url: string,
    login: string,
  }
}
interface IssueListProps {
  value: Issue[];
}

const ROADMAP_TABS = [
  { title: '正在考虑', key: 'inConsideration', color: '#FFBF00' },
  { title: '进行中', key: 'inProgress', color: '#0FC6C2' },
  { title: '最近完成', key: 'released', color: '#245CFF' },
]

const isExistInLabels = (labels: Issue['labels'], label: string) => {
  return !!labels?.find((item: { name: string }) => item.name.includes(label))
}

/**
 * 
 * @param issues 
 * @returns 
 * 正在考虑 = 带 enhancement，且没有 inprogress 和 released，且 open
 * 进行中 = 带 enhancement 和 inprogress，且 open
 * 最近完成 = 带 enhancement 和 released，且 open
 * 按点赞数量降序排序
 */
const handleSortIssues = (issues: Issue[]) => {
  const list = issues.filter((item: Issue) => isExistInLabels(item.labels, 'enhancement')).sort((item1, item2) => item2.thumbs_up - item1.thumbs_up)
  const inConsideration: Issue[] = []
  const inProgress: Issue[] = []
  const released: Issue[] = []
  list.forEach((item: Issue) => {
    const { labels } = item
    if (isExistInLabels(labels, 'inprogress')) {
      inProgress.push(item)
    } else if (isExistInLabels(labels, 'released')) {
      released.push(item)
    } else {
      inConsideration.push(item)
    }
  })
  return {
    inConsideration,
    inProgress,
    released,
  }
}

function IssueList({ value }: IssueListProps) {
  const [searchText, setSearchText] = useState<string>('');
  const [issues, setIssues] = useState<Record<string, Array<Issue>>>(handleSortIssues(value || []))

  const handleSearch = async() => {
    const result = await getIssues(searchText)
    setIssues(handleSortIssues(result || []))
  }

  const handleKeyDown = (event: any) => {
    if (event.key === 'Enter') {
      handleSearch();
    }
  };

  return (
    <>
      <Box sx={{ pb: 3 }}>
        <Grid container>
          <Grid item xs={4} display="flex" alignItems="center">
            <Stack direction="row">
              <Typography variant="h6" sx={{ mr: 2 }}>Roadmap</Typography>
            </Stack>
          </Grid>
          <Grid item xs={8}>
            <Stack direction="row" justifyContent={{ xs: "flex-start", md: "flex-end" }}>
              <Paper
                sx={{
                  display: 'flex',
                  alignItems: 'center',
                  width: 460,
                  height: 36,
                  boxShadow: "none",
                  border: "1px solid #D9D9D9",
                }}
              >
                <IconButton sx={{ p: '10px' }} aria-label="menu" onClick={handleSearch}>
                  <SearchIcon />
                </IconButton>
                <InputBase
                  sx={{ ml: 1, flex: 1 }}
                  placeholder="搜索"
                  inputProps={{ 'aria-label': 'search google maps' }}
                  value={searchText}
                  onChange={(e) =>setSearchText(e.target.value)}
                  onKeyDown={handleKeyDown}
                />
              </Paper>
              <Button
                variant="contained"
                target="_blank"
                sx={{
                  width: { xs: "102px" },
                  ml: { xs: 2 },
                  height: "36px",
                  whiteSpace: 'nowrap',
                  fontSize: { xs: '12px', md: '14px' }
                }}
                href="https://github.com/chaitin/SafeLine/issues"
              >
                提交新反馈
              </Button>
            </Stack>
          </Grid>
        </Grid>
      </Box>
      <Box
        sx={{
          display: "flex",
          justifyContent: "space-between",
          width: "100%",
          flexWrap: "wrap",
        }}
      >
        {ROADMAP_TABS.map((tab, index) => (
          <Box
            key={tab.key}
            sx={{
              width: { xs: "100%", sm: "32%" },
              height: { xs: "auto" },
              flexShrink: 0,
              borderRadius: "12px",
              mb: { xs: 3, md: 0 },
              backgroundColor: alpha(tab.color, 0.08),
            }}
          >
            <Box px={2} py={2}>
              <Typography variant="h6" color={tab.color}>{tab.title}</Typography>
              <List sx={{ py: 0, maxHeight: "830px", overflowY: "auto" }}>
                {issues[tab.key].map((issue) => (
                  <ListItem key={issue.id} sx={{ pt: 2, pb: 0, px: 0 }}>
                    <IssueItem issue={issue} />
                  </ListItem>
                ))}
              </List>
            </Box>
          </Box>
        ))}
      </Box>
    </>
  );
}


interface IssueItemProps {
  issue: Issue;
}

export const IssueItem: React.FC<IssueItemProps> = ({ issue }) => {
  if (!issue) return null

  return (
    <Box
      sx={{
        backgroundColor: "common.white",
        borderRadius: "12px",
        px: "12px",
        py: "12px",
        flex: 1,
      }}
    >
      <Typography variant="h6">{issue.title}</Typography>
      <Labels value={issue.labels} />
      <Box mt={1}>
        <Typography component="span" variant="subtitle2" sx={{ color: alpha('#000', 0.5) }}>
          {getIssueDesc(issue)}
        </Typography>
      </Box>
      <Stack direction="row" alignItems="center" mt={1}>
        <Box mr={3}>
          <Icon type="icon-pinglun" color="common.black" sx={{ mr: 1 }} />
          <Typography component="span" variant="body1" sx={{}}>
            {issue.comment_count}
          </Typography>
        </Box>
        <Box>
          <Icon type="icon-zan" color="warning.main" sx={{ mr: 1 }} />
          <Typography component="span" variant="body1" sx={{}}>
            {issue.thumbs_up}
          </Typography>
        </Box>
      </Stack>
    </Box>
  );
};

function getIssueDesc(issue: Issue) {
  const { created_at, author } = issue
  const time = formatDate(created_at)
  const authorName = author?.login || ''
  return `opened on ${time} by ${authorName}`
}
