<script lang="ts">
  import { t } from '$lib/i18n.svelte';
  import { Button } from '$lib/components/ui/button';

  let { open = $bindable(false) } = $props();
</script>

{#if open}
  <!-- svelte-ignore a11y_no_static_element_interactions -->
  <div
    class="fixed inset-0 z-50 flex justify-end"
    onkeydown={(e) => { if (e.key === 'Escape') open = false; }}
  >
    <!-- svelte-ignore a11y_click_events_have_key_events -->
    <div class="absolute inset-0 bg-black/30" onclick={() => open = false}></div>

    <div class="relative w-full max-w-xl bg-background border-l shadow-xl flex flex-col">
      <div class="border-b px-6 py-3 flex items-center justify-between shrink-0">
        <h2 class="font-semibold">{t('nav.help')}</h2>
        <Button variant="ghost" size="sm" class="h-7 text-xs" onclick={() => open = false}>
          ✕
        </Button>
      </div>
      <iframe
        src="/help.html"
        class="flex-1 w-full border-0"
        title="Help"
      ></iframe>
    </div>
  </div>
{/if}
