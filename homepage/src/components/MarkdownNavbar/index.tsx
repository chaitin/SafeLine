import { type FC, useEffect, useRef, useState } from "react";
import { Ellipsis } from "@/components";
import { useRouter } from "next/router";
import { Box } from "@mui/material";
import Link from "next/link";
import { useDebounceFn, useThrottleFn } from "ahooks";
import { slug } from "github-slugger";

interface MarkdownNavbarProps {
  container?: HTMLElement;
  source: string;
}

const MarkdownNavbar: FC<MarkdownNavbarProps> = (props) => {
  const router = useRouter();
  let { container, source } = props;
  if (!container && typeof window !== "undefined") {
    // @ts-ignore
    container = document;
  }
  const [currentListNo, setCurrentListNo] = useState<string>("");
  const [navStructure, setNavStructure] = useState<any[]>([]);
  const scrollEventLock = useRef(false);
  const timer = useRef<any>();

  const trimArrZero = (arr: number[]) => {
    let start, end;
    for (start = 0; start < arr.length; start++) {
      if (arr[start]) {
        break;
      }
    }
    for (end = arr.length - 1; end >= 0; end--) {
      if (arr[end]) {
        break;
      }
    }
    return arr.slice(start, end + 1);
  };

  const getHeadingList = () => {
    const headingList: { dataId: string; listNo: string; offsetTop: number }[] =
      [];

    navStructure.forEach((t) => {
      const headings = document.querySelectorAll(`h${t.level}`);
      const curHeading = Array.prototype.slice
        .apply(headings)
        .find(
          (h) =>
            h.innerText.trim() === t.text.trim() &&
            !headingList.find(
              (x) => x.offsetTop === h.getBoundingClientRect().top
            )
        );
      if (curHeading) {
        let rect = curHeading.getBoundingClientRect();
        headingList.push({
          dataId: `heading-${t.index}`,
          listNo: t.listNo,
          offsetTop: rect.top,
        });
      }
    });
    return headingList;
  };

  const initHeadingsId = (navStructure: any[]) => {
    navStructure.forEach((t) => {
      const headings = document.querySelectorAll(`h${t.level}`);
      const curHeading = Array.prototype.slice
        .apply(headings)
        .find(
          (h) =>
            h.innerText.trim() === t.text.trim() &&
            (!h.dataset || !h.dataset.id)
        );

      if (curHeading) {
        curHeading.dataset.id = `heading-${t.index}`;
      }
    });
  };

  const getNavStructure = (source: string) => {
    const contentWithoutCode = source
      .replace(/^[^#]+\n/g, "")
      .replace(/(?:[^\n#]+)#+\s([^#\n]+)\n*/g, "") // 匹配行内出现 # 号的情况
      .replace(/^#\s[^#\n]*\n+/, "")
      .replace(/```[^`\n]*\n+[^```]+```\n+/g, "")
      .replace(/`([^`\n]+)`/g, "$1")
      .replace(/\*\*?([^*\n]+)\*\*?/g, "$1")
      .replace(/__?([^_\n]+)__?/g, "$1")
      .trim();

    const pattOfTitle = /#+\s([^#\n]+)\n*/g;
    const matchResult = contentWithoutCode.match(pattOfTitle);

    if (!matchResult) {
      return [];
    }

    const navData = matchResult.map((r, i) => ({
      index: i,
      level: r.match(/^#+/g)![0].length,
      text: r.replace(pattOfTitle, "$1"),
      listNo: "",
    }));

    let maxLevel = 0;
    navData.forEach((t) => {
      if (t.level > maxLevel) {
        maxLevel = t.level;
      }
    });
    let matchStack = [];
    // 此部分重构，原有方法会出现次级标题后再次出现高级标题时，listNo重复的bug
    for (let i = 0; i < navData.length; i++) {
      const t = navData[i];
      const { level } = t;
      while (
        matchStack.length &&
        matchStack[matchStack.length - 1].level > level
      ) {
        matchStack.pop();
      }
      if (matchStack.length === 0) {
        const arr = new Array(maxLevel).fill(0);
        arr[level - 1] += 1;
        matchStack.push({
          level,
          arr,
        });
        t.listNo = trimArrZero(arr).join(".");
        continue;
      }
      const arr: number[] = matchStack[matchStack.length - 1].arr;
      const newArr = arr.slice();
      newArr[level - 1] += 1;
      matchStack.push({
        level,
        arr: newArr,
      });
      t.listNo = trimArrZero(newArr).join(".");
    }
    return navData;
  };

  const refreshNav = (source: string) => {
    const navStructure = getNavStructure(source);
    setNavStructure(navStructure.filter((nav) => [2, 3].includes(nav.level)));
    setCurrentListNo(navStructure[0]?.listNo);
    initHeadingsId(navStructure);
    container?.addEventListener("scroll", winScroll, false);
  };

  const { run: winScroll } = useThrottleFn(
    () => {
      clearTimeout(timer.current);
      if (scrollEventLock.current) return;
      const newHeadingList = getHeadingList().map((h) => ({
        ...h,
        distanceToTop: Math.abs(h.offsetTop),
      }));
      const distanceList = newHeadingList.map((h) => h.distanceToTop);
      const minDistance = Math.min(...distanceList);
      const curHeading = newHeadingList.find(
        (h) => h.distanceToTop === minDistance
      );
      if (!curHeading) return;
      timer.current = setTimeout(() => {
        setCurrentListNo(curHeading.listNo);
      });
    },
    { wait: 300 }
  );

  const { run: scrollToTarget } = useDebounceFn(
    (dataId: string) => {
      const target = document.querySelector(`[data-id="${dataId}"]`);
      if (target) {
        target.scrollIntoView({ behavior: "smooth" });
        scrollEventLock.current = true;
        setTimeout(() => {
          scrollEventLock.current = false;
        }, 500);
      }
    },
    { wait: 0 }
  );

  useEffect(() => {
    refreshNav(source);
    return () => {
      container?.removeEventListener("scroll", winScroll, false);
    };
  }, [source]);

  const tBlocks = navStructure.map((t, index) => {
    const cls = `title-anchor title-level${t.level}`;
    const point = slug(t.text);
    let hash = router.asPath.split("#")[1];
    return (
      <Box
        className={cls}
        component={Link}
        href={`#${point}`}
        style={{ width: "100%" }}
        key={`title_anchor_${Math.random().toString(36).substring(2)}`}
      >
        <Ellipsis
          sx={{
            "&:hover": {
              color: "primary.main",
            },
            color:
              decodeURIComponent(hash) === point ? "primary.main" : "inherit",
          }}
        >
          {t.text}
        </Ellipsis>
      </Box>
    );
  });
  return <div className={`markdown-navigation`}>{tBlocks}</div>;
};

export default MarkdownNavbar;
