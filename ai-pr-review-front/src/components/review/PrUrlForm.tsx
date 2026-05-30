import { useState, type FormEvent } from 'react';
import { Button } from '@/components/ui/Button';
import { Input } from '@/components/ui/Input';
import { isValidGitHubPrUrl } from '@/utils/pr';

interface PrUrlFormProps {
  onSubmit: (prUrl: string) => void;
  loading?: boolean;
  disabled?: boolean;
  serverError?: string | null;
}

export function PrUrlForm({
  onSubmit,
  loading = false,
  disabled = false,
  serverError = null,
}: PrUrlFormProps) {
  const [prUrl, setPrUrl] = useState('');
  const [validationError, setValidationError] = useState<string | null>(null);

  function handleSubmit(e: FormEvent) {
    e.preventDefault();
    const trimmed = prUrl.trim();

    if (!trimmed) {
      setValidationError('请输入 PR 链接');
      return;
    }

    if (!isValidGitHubPrUrl(trimmed)) {
      setValidationError(
        '请输入有效的 GitHub PR 链接，例如 https://github.com/org/repo/pull/123',
      );
      return;
    }

    setValidationError(null);
    onSubmit(trimmed);
  }

  const displayError = validationError || serverError;

  return (
    <form onSubmit={handleSubmit} className="space-y-6">
      <div>
        <label htmlFor="pr-url" className="mb-2 flex items-center gap-2 text-sm font-medium text-white/70">
          <span>🔗</span> GitHub PR 链接
        </label>
        <Input
          id="pr-url"
          type="url"
          placeholder="https://github.com/org/repo/pull/123"
          value={prUrl}
          onChange={(e) => {
            setPrUrl(e.target.value);
            setValidationError(null);
          }}
          error={displayError ?? undefined}
          disabled={loading || disabled}
        />
        <div className="mt-3 flex flex-wrap gap-2">
          {['仅支持 GitHub', '同步审查', '约 1-3 分钟'].map((tag) => (
            <span
              key={tag}
              className="rounded-full border border-white/10 bg-white/5 px-2.5 py-0.5 text-[11px] text-white/40"
            >
              {tag}
            </span>
          ))}
        </div>
      </div>
      <Button type="submit" loading={loading} disabled={disabled} className="w-full sm:w-auto">
        ✨ 开始审查
      </Button>
    </form>
  );
}
