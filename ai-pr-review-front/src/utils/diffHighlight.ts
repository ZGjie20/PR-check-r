import { SEVERITY_ORDER } from '@/constants/severity';
import type { ReviewResultBySeverity, Severity } from '@/types/review';

const SEVERITY_RANK: Record<Severity, number> = {
  high: 3,
  medium: 2,
  low: 1,
};

export interface DiffDisplayLine {
  text: string;
  file?: string;
  newLine?: number;
  changeType: 'add' | 'delete' | 'context' | 'meta';
  severity?: Severity;
}

export function buildIssueLineMap(
  reviewResult: ReviewResultBySeverity,
): Map<string, Map<number, Severity>> {
  const map = new Map<string, Map<number, Severity>>();

  for (const severity of SEVERITY_ORDER) {
    for (const issue of reviewResult[severity]) {
      if (!map.has(issue.file)) {
        map.set(issue.file, new Map());
      }
      const fileMap = map.get(issue.file)!;
      const existing = fileMap.get(issue.line);
      if (!existing || SEVERITY_RANK[severity] > SEVERITY_RANK[existing]) {
        fileMap.set(issue.line, severity);
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
  const direct = issueMap.get(file)?.get(line);
  if (direct) {
    return direct;
  }

  const basename = file.split('/').pop() ?? file;
  for (const [issueFile, lineMap] of issueMap.entries()) {
    const issueBasename = issueFile.split('/').pop() ?? issueFile;
    if (issueBasename === basename || issueFile === file) {
      const severity = lineMap.get(line);
      if (severity) {
        return severity;
      }
    }
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
      const match = rawLine.match(/^diff --git a\/(.+?) b\/(.+)$/);
      currentFile = match ? match[2] : '';
      newLine = 0;
      result.push({ text: rawLine, changeType: 'meta' });
      continue;
    }

    if (rawLine.startsWith('@@')) {
      const match = rawLine.match(/\+(\d+)(?:,(\d+))?/);
      newLine = match ? Number.parseInt(match[1], 10) : 0;
      result.push({ text: rawLine, file: currentFile, changeType: 'meta' });
      continue;
    }

    if (
      rawLine.startsWith('+++') ||
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

    const prefix = rawLine[0];

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
      const severity = lookupSeverity(issueMap, currentFile, newLine);
      result.push({
        text: rawLine,
        file: currentFile,
        newLine,
        changeType: 'context',
        severity,
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
