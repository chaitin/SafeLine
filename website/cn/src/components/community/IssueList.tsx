import React, { useState } from "react";
import {
  Box,
  Link,
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
  updated_at: number
  author: {
    avatar_url: string,
    login: string,
  }
}

export type Issues = {
  in_consideration?: Issue[]
  in_progress?: Issue[]
  released?: Issue[]
}

interface IssueListProps {
  value: Issues;
}

const ROADMAP_TABS = [
  { title: '正在考虑', key: 'in_consideration', color: '#FFBF00' },
  { title: '进行中', key: 'in_progress', color: '#0FC6C2' },
  { title: '最近完成', key: 'released', color: '#245CFF' },
]

function IssueList({ value }: IssueListProps) {
  const [searchText, setSearchText] = useState<string>('');
  const [issues, setIssues] = useState<Record<string, Array<Issue>>>(value || {})

  const handleSearch = async() => {
    const result = await getIssues(searchText)
    setIssues(result || {})
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
          <Grid item xs={12} sm={4} mb={{ xs: 4, sm: 0 }} display="flex" justifyContent={{ xs: "center", sm: "flex-start" }}>
            <Stack direction="row">
              <Typography variant="h5" sx={{ mr: 2 }}>开发计划</Typography>
            </Stack>
          </Grid>
          <Grid item xs={12} sm={8}>
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
                  ml: { xs: 2 },
                  minWidth: "80px",
                  height: "36px",
                  whiteSpace: 'nowrap',
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
              <Typography variant="h6" color={tab.color} mb={2}>
                {tab.title}
                （{issues[tab.key]?.length || 0}）
              </Typography>
              <List sx={{ py: 0, maxHeight: "790px", overflowY: "auto" }}>
                {issues[tab.key]?.map((issue, index) => (
                  <ListItem key={issue.id} sx={{ pt: index > 0 ? 2 : 0, pb: 0, px: 0 }}>
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
      <Link href={issue.url} target="_blank">
        <Typography variant="h6">{issue.title}</Typography>
      </Link>
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
