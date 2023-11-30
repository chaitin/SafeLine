import React, { useState } from "react";
import {
  Box,
  Table,
  TableRow,
  TableBody,
  TableContainer,
  Stack,
  Chip,
  Button,
  Typography,
  alpha,
  AvatarGroup,
  Avatar,
  Grid,
  Paper,
  IconButton,
  InputBase,
  SxProps,
  Link,
} from "@mui/material";
import SearchIcon from '@mui/icons-material/Search';
import TableCell from "@mui/material/TableCell";
import Icon from "@/components/Icon";
import ArrowDropUpOutlinedIcon from '@mui/icons-material/ArrowDropUpOutlined';
import { formatDate } from '@/common/utils'
import { getDiscussions } from "@/api";

type User = {
  avatar_url: string
  login: string
}

export type Discussion = {
  id: string
  labels: { name: string, color: string }[]
  thumbs_up: number
  title: string
  url: string
  comment_count: number
  created_at: number
  author: User
  comment_users: User[]
  category: { name: string, emoji_html: string }
  is_answered: boolean
  upvote_count: number
}

interface DiscussionListProps {
  value: Discussion[];
}

export default function DiscussionList({ value }: DiscussionListProps) {
  const [discussionList, setDiscussionList] = useState<Array<Discussion>>(value || []);
  const [searchText, setSearchText] = useState<string>('');

  const handleSearch = async() => {
    const result = await getDiscussions(searchText)
    setDiscussionList(result || [])
  }

  const handleKeyDown = (event: any) => {
    if (event.key === 'Enter') {
      handleSearch();
    }
  };

  return (
    <>
      <Box
        sx={{
          pb: 3,
          borderBottom: "1px solid #E3E8EF",
        }}
      >
        <Grid container>
          <Grid item xs={4} display="flex" alignItems="center">
            <Stack direction="row">
              <Typography variant="h6" sx={{ mr: 2 }}>讨论</Typography>
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
                href="https://github.com/chaitin/SafeLine/discussions"
              >
                查看更多
              </Button>
            </Stack>
          </Grid>
        </Grid>
      </Box>
      <TableContainer>
        <Table
          sx={{
            ".MuiTableCell-root": {
              py: 2,
              color: "common.black",
              borderRight: "none",
              borderColor: "#E3E8EF",
              "&:first-of-type": {
                paddingLeft: "16px",
              },
            },
          }}
        >
          <TableBody>
            {discussionList.map((discussion) => (
              <TableRow key={discussion.id}>
                <TableCell sx={{ width: "85%" }}>
                  <Stack direction="row" alignItems="center">
                    <Box mr={2}>
                      <Typography variant="subtitle2" sx={{ color: "error.main" }} display="flex" alignItems="center">
                        <ArrowDropUpOutlinedIcon />
                        {discussion.upvote_count}
                      </Typography>
                    </Box>
                    <Box
                      mr={3}
                      sx={{
                        width: "39px",
                        height: "39px",
                        background: "#F6F7F8",
                        borderRadius: "4px",
                      }}
                      className="flex justify-center items-center"
                    >
                      <Box dangerouslySetInnerHTML={{ __html: discussion.category?.emoji_html }}></Box>
                    </Box>
                    <Box>
                      <Link href={discussion.url} target="_blank">
                        <Typography variant="h6">{discussion.title}</Typography>
                      </Link>
                      <Labels value={discussion.labels} />
                      <Box>
                        <Typography component="span" variant="subtitle2" sx={{ color: alpha('#000', 0.5) }}>
                          {getDiscussionDesc(discussion)}
                        </Typography>
                        {isQA(discussion.category?.name) ? (
                          <Typography
                            component="span"
                            variant="subtitle2"
                            ml={1}
                            sx={{ color: discussion.is_answered ? "primary.main" : alpha('#000', 0.5) }}
                          >
                            · {discussion.is_answered ? 'Answered' : 'Unanswered'}
                          </Typography>
                        ) : null}
                      </Box>
                    </Box>
                  </Stack>
                </TableCell>
                <TableCell sx={{ width: "15%" }}>
                  <Stack direction="row" alignItems="center">
                    {renderCommentIcon(discussion)}
                    <Typography component="span" variant="body1" sx={{ ml: 1 }} >
                      {discussion.comment_count}
                    </Typography>
                  </Stack>
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </TableContainer>
    </>
  );
}

function isQA(name: string) {
  return name && name.includes('Q&A')
}

function getDiscussionDesc(discussion: Discussion) {
  const { created_at, category, author } = discussion
  const time = formatDate(created_at)
  const action = isQA(category.name) ? 'asked' : 'started'
  const authorName = author?.login || ''
  return `${authorName} ${action} on ${time} in ${category.name}`
}

const renderCommentIcon = ({ category, is_answered }: { category: { name: string }, is_answered: boolean }) => {
  if (isQA(category.name)) {
    if (is_answered) {
      return (<Icon type="icon-a-yiwanchengtianchong" color="#02BFA5" />)
    }
    return (<Icon type="icon-yiwancheng" color="common.black" />)
  }
  return (<Icon type="icon-pinglun" color="common.black" />)
}

interface LabelsProps {
  value: { color: string, name: string }[];
  sx?: SxProps;
}

export const Labels: React.FC<LabelsProps> = ({ value, sx }) => {
  if (!value || !value.length) return null

  return (
    <Box sx={{ ...sx }}>
      {value.map((l, index) => (
        <Chip
          key={index}
          variant="filled"
          label={l.name}
          size="small"
          sx={{ backgroundColor: `#${l.color}`, color: "common.white", mr: 1, height: 16, fontSize: "12px" }}
        />
      ))}
    </Box>
  );
};

interface AvatarsProps {
  value: User[];
  author: User,
}

export const Avatars: React.FC<AvatarsProps> = ({ value }) => {
  const visibleAvatars = value.slice(0, 5)
  return (
    <AvatarGroup
      sx={{
        ".MuiAvatar-root": {
          width: "24px",
          height: "24px",
        },
      }}
    >
      {visibleAvatars.map((avatar, index) => (
        <Avatar key={index} alt={avatar.login} src={avatar.avatar_url} />
      ))}
    </AvatarGroup>
  );
};
