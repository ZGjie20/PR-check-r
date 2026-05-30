import { SEVERITY_ORDER } from '@/constants/severity';
import type { ReviewResultBySeverity, Severity } from '@/types/review';

const SEVERITY_RANK: Record<Severity, number> = {
  high: 3,
  medium: 2,
  low: 1,
};

const HUNK_HEADER_PATTERN = /^@@ -\d+(?:,\d+)? \+(\d+)(?:,\d+)? @@/;
const CODE_HINT_PATTERN =
  /`([^`]{2,80})`|(\b[\w.$:{}()[\]=;+-]+\b(?:\s*[=({][^;]{0,60})?)/g;
const SEARCH_RADIUS = 8;

export interface DiffDisplayLine {
  text: string;
  file?: string;
  newLine?: number;
  changeType: 'add' | 'delete' | 'context' | 'meta';
  severity?: Severity;
}

interface ParsedDiffLine {
  file: string;
  newLine: number;
  content: string;
  changeType: 'add' | 'context';
}

interface IssueCandidate {
  file: string;
  line: number;
  severity: Severity;
  message: string;
  suggestion: string;
}

function normalizePath(path: string): string {
  return path.replace(/\\/g, '/').replace(/^\.\//, '').replace(/^b\//, '');
}

function stripNoNewlineSuffix(line: string): string {
  const suffix = '\\ No newline at end of file';
  if (line.endsWith(suffix)) {
    return line.slice(0, -suffix.length);
  }
  return line;
}

function normalizeContent(content: string): string {
  return content.trim().replace(/\s+/g, ' ').toLowerCase();
}

function compactContent(content: string): string {
  return content.replace(/\s+/g, '').toLowerCase();
}

function extractCodeHints(message: string, suggestion: string): string[] {
  const text = `${message} ${suggestion}`;
  const hints = new Set<string>();

  for (const match of text.matchAll(CODE_HINT_PATTERN)) {
    const hint = (match[1] ?? match[2] ?? '').trim();
    if (hint.length < 2) {
      continue;
    }
    if (/^[\u4e00-\u9fff]+$/.test(hint)) {
      continue;
    }
    if (/^(var|关键字|修复|改为|建议|例如|通过|环境变量)$/.test(hint)) {
      continue;
    }
    hints.add(hint);
  }

  if (/token/i.test(text)) {
    hints.add('token:');
    hints.add('token');
  }

  return [...hints].sort((a, b) => b.length - a.length);
}

function lineMatchesHints(content: string, hints: string[]): boolean {
  if (hints.length === 0) {
    return false;
  }

  const normalized = normalizeContent(content);
  const compact = compactContent(content);

  return hints.some((hint) => {
    const normalizedHint = normalizeContent(hint);
    const compactHint = compactContent(hint);

    if (normalizedHint.length >= 4) {
      if (
        normalized.includes(normalizedHint) ||
        compact.includes(compactHint)
      ) {
        return true;
      }
    }

    if (/^token:?$/i.test(hint) && /token\s*:/i.test(content)) {
      return true;
    }

    return false;
  });
}

function findFileLines(
  fileIndex: Map<string, ParsedDiffLine[]>,
  issueFile: string,
): ParsedDiffLine[] {
  const normalizedIssueFile = normalizePath(issueFile);
  const direct = fileIndex.get(normalizedIssueFile);
  if (direct) {
    return direct;
  }

  const basename = normalizedIssueFile.split('/').pop() ?? normalizedIssueFile;
  const matches: ParsedDiffLine[] = [];

  for (const [file, lines] of fileIndex.entries()) {
    const fileBasename = file.split('/').pop() ?? file;
    if (fileBasename === basename || file === normalizedIssueFile) {
      matches.push(...lines);
    }
  }

  return matches;
}

function resolveIssueLine(
  issue: IssueCandidate,
  fileLines: ParsedDiffLine[],
): number | undefined {
  const addedLines = fileLines.filter((line) => line.changeType === 'add');
  const hints = extractCodeHints(issue.message, issue.suggestion);

  const exactAdded = addedLines.find((line) => line.newLine === issue.line);
  if (exactAdded) {
    return exactAdded.newLine;
  }

  const hintedAdded = addedLines.filter((line) =>
    lineMatchesHints(line.content, hints),
  );
  if (hintedAdded.length === 1) {
    return hintedAdded[0].newLine;
  }
  if (hintedAdded.length > 1) {
    return hintedAdded.sort(
      (a, b) =>
        Math.abs(a.newLine - issue.line) - Math.abs(b.newLine - issue.line),
    )[0].newLine;
  }

  const nearbyAdded = addedLines
    .filter(
      (line) =>
        hints.length > 0 &&
        Math.abs(line.newLine - issue.line) <= SEARCH_RADIUS &&
        lineMatchesHints(line.content, hints),
    )
    .sort(
      (a, b) =>
        Math.abs(a.newLine - issue.line) - Math.abs(b.newLine - issue.line),
    );

  if (nearbyAdded.length > 0) {
    return nearbyAdded[0].newLine;
  }

  if (hints.length === 0) {
    const closestAdded = addedLines
      .filter((line) => Math.abs(line.newLine - issue.line) <= 2)
      .sort(
        (a, b) =>
          Math.abs(a.newLine - issue.line) - Math.abs(b.newLine - issue.line),
      )[0];
    if (closestAdded) {
      return closestAdded.newLine;
    }
  }

  const exactContext = fileLines.find(
    (line) => line.newLine === issue.line && line.changeType === 'context',
  );
  if (exactContext && lineMatchesHints(exactContext.content, hints)) {
    return exactContext.newLine;
  }

  return undefined;
}

function parseDiffIndex(diff: string): Map<string, ParsedDiffLine[]> {
  const fileIndex = new Map<string, ParsedDiffLine[]>();
  const rawLines = diff.split('\n');
  let currentFile = '';
  let newLine = 0;

  for (const rawLine of rawLines) {
    if (rawLine.startsWith('diff --git ')) {
      currentFile = '';
      newLine = 0;
      continue;
    }

    if (rawLine.startsWith('+++ ')) {
      let path = rawLine.slice(4).trim();
      if (path !== '/dev/null') {
        currentFile = normalizePath(path);
        if (!fileIndex.has(currentFile)) {
          fileIndex.set(currentFile, []);
        }
      }
      newLine = 0;
      continue;
    }

    if (!currentFile) {
      continue;
    }

    if (rawLine.startsWith('@@')) {
      const match = rawLine.match(HUNK_HEADER_PATTERN);
      newLine = match ? Number.parseInt(match[1], 10) : 0;
      continue;
    }

    if (
      rawLine.startsWith('---') ||
      rawLine.startsWith('index ') ||
      rawLine.startsWith('new file') ||
      rawLine.startsWith('deleted file') ||
      rawLine.startsWith('rename from') ||
      rawLine.startsWith('rename to') ||
      rawLine.startsWith('similarity index') ||
      rawLine === ''
    ) {
      continue;
    }

    const cleanedLine = stripNoNewlineSuffix(rawLine);
    const prefix = cleanedLine[0];
    if (!prefix || (prefix !== '+' && prefix !== '-' && prefix !== ' ')) {
      continue;
    }

    const content = cleanedLine.slice(1);

    if (prefix === '+') {
      fileIndex.get(currentFile)!.push({
        file: currentFile,
        newLine,
        content,
        changeType: 'add',
      });
      newLine += 1;
      continue;
    }

    if (prefix === '-') {
      continue;
    }

    fileIndex.get(currentFile)!.push({
      file: currentFile,
      newLine,
      content,
      changeType: 'context',
    });
    newLine += 1;
  }

  return fileIndex;
}

function flattenIssues(reviewResult: ReviewResultBySeverity): IssueCandidate[] {
  const issues: IssueCandidate[] = [];

  for (const severity of SEVERITY_ORDER) {
    for (const issue of reviewResult[severity]) {
      issues.push({
        file: issue.file,
        line: issue.line,
        severity,
        message: issue.message,
        suggestion: issue.suggestion,
      });
    }
  }

  return issues;
}

export function buildIssueLineMap(
  reviewResult: ReviewResultBySeverity,
  diff: string,
): Map<string, Map<number, Severity>> {
  const map = new Map<string, Map<number, Severity>>();
  const fileIndex = parseDiffIndex(diff);
  const issues = flattenIssues(reviewResult);

  for (const issue of issues) {
    const fileLines = findFileLines(fileIndex, issue.file);
    const resolvedLine = resolveIssueLine(issue, fileLines);
    if (resolvedLine === undefined) {
      continue;
    }

    const normalizedFile = normalizePath(issue.file);
    if (!map.has(normalizedFile)) {
      map.set(normalizedFile, new Map());
    }

    const fileMap = map.get(normalizedFile)!;
    const existing = fileMap.get(resolvedLine);
    if (!existing || SEVERITY_RANK[issue.severity] > SEVERITY_RANK[existing]) {
      fileMap.set(resolvedLine, issue.severity);
    }

    const actualFile = fileLines.find((line) => line.newLine === resolvedLine)?.file;
    if (actualFile && actualFile !== normalizedFile) {
      if (!map.has(actualFile)) {
        map.set(actualFile, new Map());
      }
      const actualFileMap = map.get(actualFile)!;
      const actualExisting = actualFileMap.get(resolvedLine);
      if (
        !actualExisting ||
        SEVERITY_RANK[issue.severity] > SEVERITY_RANK[actualExisting]
      ) {
        actualFileMap.set(resolvedLine, issue.severity);
      }
    }
  }

  return map;
}

function lookupSeverity(
  issueMap: Map<string, Map<number, Severity>>,
  file: string,
  line: number,
): Severity | undefined {
  const normalizedFile = normalizePath(file);
  const direct = issueMap.get(normalizedFile)?.get(line);
  if (direct) {
    return direct;
  }

  const basename = normalizedFile.split('/').pop() ?? normalizedFile;
  const candidates: Array<Map<number, Severity>> = [];

  for (const [issueFile, lineMap] of issueMap.entries()) {
    const issueBasename = issueFile.split('/').pop() ?? issueFile;
    if (issueBasename === basename) {
      candidates.push(lineMap);
    }
  }

  if (candidates.length === 1) {
    return candidates[0].get(line);
  }

  return undefined;
}

export function parseDiffWithIssues(
  diff: string,
  issueMap: Map<string, Map<number, Severity>>,
): DiffDisplayLine[] {
  const lines = diff.split('\n');
  const result: DiffDisplayLine[] = [];
  let currentFile = '';
  let newLine = 0;

  for (const rawLine of lines) {
    if (rawLine.startsWith('diff --git ')) {
      currentFile = '';
      newLine = 0;
      result.push({ text: rawLine, changeType: 'meta' });
      continue;
    }

    if (rawLine.startsWith('+++ ')) {
      let path = rawLine.slice(4).trim();
      currentFile = path !== '/dev/null' ? normalizePath(path) : '';
      newLine = 0;
      result.push({ text: rawLine, changeType: 'meta' });
      continue;
    }

    if (rawLine.startsWith('@@')) {
      const match = rawLine.match(HUNK_HEADER_PATTERN);
      newLine = match ? Number.parseInt(match[1], 10) : 0;
      result.push({ text: rawLine, file: currentFile, changeType: 'meta' });
      continue;
    }

    if (
      rawLine.startsWith('---') ||
      rawLine.startsWith('index ') ||
      rawLine.startsWith('new file') ||
      rawLine.startsWith('deleted file') ||
      rawLine.startsWith('rename from') ||
      rawLine.startsWith('rename to') ||
      rawLine.startsWith('similarity index') ||
      rawLine === ''
    ) {
      result.push({ text: rawLine, changeType: 'meta' });
      continue;
    }

    const cleanedLine = stripNoNewlineSuffix(rawLine);
    const prefix = cleanedLine[0];

    if (prefix === '+') {
      const severity = lookupSeverity(issueMap, currentFile, newLine);
      result.push({
        text: rawLine,
        file: currentFile,
        newLine,
        changeType: 'add',
        severity,
      });
      newLine += 1;
      continue;
    }

    if (prefix === '-') {
      result.push({
        text: rawLine,
        file: currentFile,
        changeType: 'delete',
      });
      continue;
    }

    if (prefix === ' ') {
      result.push({
        text: rawLine,
        file: currentFile,
        newLine,
        changeType: 'context',
      });
      newLine += 1;
      continue;
    }

    result.push({ text: rawLine, changeType: 'meta' });
  }

  return result;
}

export function countHighlightedLines(lines: DiffDisplayLine[]): number {
  return lines.filter((line) => line.severity).length;
}
