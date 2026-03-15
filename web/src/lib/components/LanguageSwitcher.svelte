<script lang="ts">
  import { i18n, type Locale } from '$lib/i18n.svelte';

  function toggle() {
    i18n.locale = i18n.locale === 'zh' ? 'en' : 'zh';
  }

  const label: Record<Locale, string> = { zh: '中', en: 'EN' };
</script>

<button
  onclick={toggle}
  class="h-7 px-2 text-xs font-medium border rounded hover:bg-muted transition-colors"
  aria-label="Switch language"
>
  {label[i18n.locale]}
</button>
